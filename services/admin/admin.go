package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	adminJWTCookieName = "Authorization"
	adminJWTDuration   = 10 * 24 * time.Hour
)

type User struct {
	userName string
	email    string
	password string
}

type Problem struct {
	ProblemID   int32  `json:"problemId"`
	AuthorID    int32  `json:"authorId"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	Constraints string `json:"constraints"`
	TimeLimit   int32  `json:"timeLimit"`
	MemoryLimit int32  `json:"memoryLimit"`
}

type adminJWTClaims struct {
	UserID    int32 `json:"userId"`
	IssuedAt  int64 `json:"iat"`
	ExpiresAt int64 `json:"exp"`
}

func adminProtectedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := verifyAdminJWTFromRequest(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or missing admin token",
			})
			return
		}

		c.Set("admin_user_id", claims.UserID)
		c.Next()
	}
}

func populateUser(c *gin.Context) (User, error) {
	var body struct {
		UserName string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		return User{}, err
	}

	u := User{
		userName: strings.TrimSpace(body.UserName),
		email:    strings.TrimSpace(body.Email),
		password: strings.TrimSpace(body.Password),
	}
	if u.userName == "" || u.email == "" || u.password == "" {
		return User{}, errors.New("username, email, and password are required")
	}

	return u, nil
}

func populateProblem(c *gin.Context) (Problem, error) {
	var body struct {
		AuthorID    int32  `json:"authorid"`
		Name        string `json:"name"`
		Body        string `json:"body"`
		Constraints string `json:"constraints"`
		TimeLimit   int32  `json:"time_limit"`
		MemoryLimit int32  `json:"memory_limit"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		return Problem{}, err
	}

	problem := Problem{
		AuthorID:    body.AuthorID,
		Name:        strings.TrimSpace(body.Name),
		Body:        strings.TrimSpace(body.Body),
		Constraints: strings.TrimSpace(body.Constraints),
		TimeLimit:   body.TimeLimit,
		MemoryLimit: body.MemoryLimit,
	}
	if problem.AuthorID == 0 || problem.Name == "" || problem.Body == "" ||
		problem.Constraints == "" || problem.TimeLimit == 0 || problem.MemoryLimit == 0 {
		return Problem{}, errors.New("authorid, name, body, constraints, time_limit, and memory_limit are required")
	}

	return problem, nil
}

func handleAdminSignIn(c *gin.Context) {
	user, err := populateUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "database is not configured",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID, err := database.AdminUserIDByCredentials(ctx, user)
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid admin credentials",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to sign in",
		})
		return
	}

	token, err := signAdminJWT(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to sign jwt",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		adminJWTCookieName,
		"Bearer "+token,
		int(adminJWTDuration.Seconds()),
		"/v1/admin",
		"",
		gin.Mode() == gin.ReleaseMode,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"userId": userID,
	})
}

func handleAdminCreateProblem(c *gin.Context) {
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "database is not configured",
		})
		return
	}

	problem, err := populateProblem(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	problemID, err := database.createNewProblem(ctx, problem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create problem",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"problemId": problemID,
	})
}

func handleAdminGetProblem(c *gin.Context) {
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "database is not configured",
		})
		return
	}

	problemID, err := int32Param(c, "problemId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid problemId",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	problem, err := database.getProblem(ctx, problemID)
	if err != nil {
		if errors.Is(err, errProblemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "problem not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get problem",
		})
		return
	}

	c.JSON(http.StatusOK, problem)
}

func handleAdminUploadInput(c *gin.Context) {
	handleAdminProblemFileUpload(c, uploadInput, database.uploadInput)
}

func handleAdminUploadOutput(c *gin.Context) {
	handleAdminProblemFileUpload(c, uploadOutput, database.uploadOutput)
}

func handleAdminUploadChecker(c *gin.Context) {
	handleAdminProblemFileUpload(c, uploadChecker, database.uploadChecker)
}

func handleAdminProblemFileUpload(
	c *gin.Context,
	uploadFile func(*multipart.FileHeader) (int32, error),
	saveFile func(context.Context, int32, int32, int32) error,
) {
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "database is not configured",
		})
		return
	}

	problemID, err := int32Param(c, "problemId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid problemId",
		})
		return
	}

	userID, err := int32Param(c, "userId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid userId",
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file is required",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	problem, err := database.getProblem(ctx, problemID)
	if err != nil {
		if errors.Is(err, errProblemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "problem not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get problem",
		})
		return
	}
	if problem.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "user is not the problem author",
		})
		return
	}

	uploadedID, err := uploadFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to upload file",
		})
		return
	}

	if err := saveFile(ctx, userID, problemID, uploadedID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save uploaded file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": uploadedID,
	})
}

func uploadInput(file *multipart.FileHeader) (int32, error) {
	return 0, nil
}

func uploadOutput(file *multipart.FileHeader) (int32, error) {
	return 0, nil
}

func uploadChecker(file *multipart.FileHeader) (int32, error) {
	return 0, nil
}

func signAdminJWT(userID int32) (string, error) {
	secret := jwtSecret()
	if secret == "" {
		return "", errors.New("secret_key is required")
	}

	now := time.Now()
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	claims := map[string]any{
		"userId": userID,
		"iat":    now.Unix(),
		"exp":    now.Add(adminJWTDuration).Unix(),
	}

	unsignedToken, err := encodedJWTPart(header, claims)
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(unsignedToken))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return unsignedToken + "." + signature, nil
}

func verifyAdminJWTFromRequest(c *gin.Context) (*adminJWTClaims, error) {
	cookieValue, err := c.Cookie(adminJWTCookieName)
	if err != nil {
		return nil, err
	}

	token, ok := strings.CutPrefix(cookieValue, "Bearer ")
	if !ok || token == "" {
		return nil, errors.New("admin token is missing bearer prefix")
	}

	return verifyAdminJWT(token)
}

func verifyAdminJWT(token string) (*adminJWTClaims, error) {
	secret := jwtSecret()
	if secret == "" {
		return nil, errors.New("secret_key is required")
	}

	header, payload, signature, ok := splitJWT(token)
	if !ok {
		return nil, errors.New("admin token is malformed")
	}

	unsignedToken := header + "." + payload
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(unsignedToken))
	expectedSignature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return nil, errors.New("admin token signature is invalid")
	}

	headerJSON, err := base64.RawURLEncoding.DecodeString(header)
	if err != nil {
		return nil, err
	}

	var tokenHeader struct {
		Algorithm string `json:"alg"`
		Type      string `json:"typ"`
	}
	if err := json.Unmarshal(headerJSON, &tokenHeader); err != nil {
		return nil, err
	}
	if tokenHeader.Algorithm != "HS256" || tokenHeader.Type != "JWT" {
		return nil, errors.New("admin token header is invalid")
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}

	var claims adminJWTClaims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return nil, err
	}
	if claims.UserID == 0 || claims.ExpiresAt == 0 {
		return nil, errors.New("admin token claims are invalid")
	}
	if time.Now().Unix() > claims.ExpiresAt {
		return nil, errors.New("admin token is expired")
	}

	return &claims, nil
}

func splitJWT(token string) (string, string, string, bool) {
	firstDot := strings.IndexByte(token, '.')
	if firstDot == -1 {
		return "", "", "", false
	}

	secondDot := strings.IndexByte(token[firstDot+1:], '.')
	if secondDot == -1 {
		return "", "", "", false
	}
	secondDot += firstDot + 1

	header := token[:firstDot]
	payload := token[firstDot+1 : secondDot]
	signature := token[secondDot+1:]
	if header == "" || payload == "" || signature == "" {
		return "", "", "", false
	}

	return header, payload, signature, true
}

func encodedJWTPart(header map[string]string, claims map[string]any) (string, error) {
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(headerJSON) + "." +
		base64.RawURLEncoding.EncodeToString(claimsJSON), nil
}

func jwtSecret() string {
	if secret := os.Getenv("secret_key"); secret != "" {
		return secret
	}
	if secret := os.Getenv("SECRET_KEY"); secret != "" {
		return secret
	}
	return ""
}

func int32Param(c *gin.Context, name string) (int32, error) {
	value, err := strconv.ParseInt(c.Param(name), 10, 32)
	if err != nil || value <= 0 {
		return 0, errors.New("invalid int32 param")
	}

	return int32(value), nil
}
