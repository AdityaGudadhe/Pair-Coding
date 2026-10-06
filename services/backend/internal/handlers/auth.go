package handlers

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

const (
	// AuthCookieName matches the cookie the admin service uses, so both read
	// the token the same way: "Bearer <jwt>".
	AuthCookieName = "Authorization"
	tokenDuration  = time.Hour

	bcryptCost   = 12
	maxBodyBytes = 4 << 10
	dbTimeout    = 5 * time.Second

	maxEmailLen = 254

	pgUniqueViolation = "23505"
)

var (
	usernamePattern = regexp.MustCompile(`^[A-Za-z0-9_.-]{1,20}$`)
	emailPattern    = regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9-]+(\.[A-Za-z0-9-]+)*\.[A-Za-z]{2,}$`)
	// Printable ASCII without space: letters, digits and the usual symbols.
	passwordPattern = regexp.MustCompile(`^[!-~]{8,30}$`)

	errInvalidCredentials = errors.New("invalid email or password")
)

// dummyHash is compared against when the email is unknown so a failed login
// takes the same time whether or not the account exists.
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("dummy-password-for-timing"), bcryptCost)

type AuthHandler struct {
	db     *pgxpool.Pool
	secret []byte
}

func NewAuthHandler(db *pgxpool.Pool, secret []byte) *AuthHandler {
	return &AuthHandler{db: db, secret: secret}
}

type signupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (h *AuthHandler) SignupHandler(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBodyBytes)

	var req signupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "invalid request body")
		return
	}

	username := strings.TrimSpace(req.Username)
	email := strings.ToLower(strings.TrimSpace(req.Email))
	password := req.Password

	if username == "" || email == "" || password == "" {
		c.String(http.StatusBadRequest, "username, email and password are required")
		return
	}
	if err := validateUsername(username); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if err := validateEmail(email); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if err := validatePassword(password); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), dbTimeout)
	defer cancel()

	var usernameTaken, emailTaken bool
	err := h.db.QueryRow(ctx,
		`SELECT
			EXISTS (SELECT 1 FROM users WHERE LOWER(user_name) = LOWER($1)),
			EXISTS (SELECT 1 FROM users WHERE email = $2)`,
		username, email,
	).Scan(&usernameTaken, &emailTaken)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to sign up")
		return
	}
	if usernameTaken {
		c.String(http.StatusConflict, "username is already taken")
		return
	}
	if emailTaken {
		c.String(http.StatusConflict, "email is already registered")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to sign up")
		return
	}

	var userID int32
	err = h.db.QueryRow(ctx,
		`INSERT INTO users (user_name, email, password) VALUES ($1, $2, $3) RETURNING user_id`,
		username, email, string(hash),
	).Scan(&userID)
	if err != nil {
		// Another request may have registered the same name or email between
		// the existence check and this insert; the unique indexes catch it.
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			c.String(http.StatusConflict, "username or email is already registered")
			return
		}
		c.String(http.StatusInternalServerError, "failed to sign up")
		return
	}

	if err := h.setAuthCookie(c, userID, username); err != nil {
		c.String(http.StatusInternalServerError, "failed to create session")
		return
	}

	c.String(http.StatusOK, "signup successful")
}

func (h *AuthHandler) LoginHandler(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBodyBytes)

	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "invalid request body")
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	password := req.Password

	if email == "" || password == "" {
		c.String(http.StatusBadRequest, "email and password are required")
		return
	}
	if validateEmail(email) != nil || validatePassword(password) != nil {
		c.String(http.StatusUnauthorized, errInvalidCredentials.Error())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), dbTimeout)
	defer cancel()

	var (
		userID   int32
		username string
		hash     string
	)
	err := h.db.QueryRow(ctx,
		`SELECT user_id, user_name, password FROM users WHERE email = $1`,
		email,
	).Scan(&userID, &username, &hash)
	if errors.Is(err, pgx.ErrNoRows) {
		_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(password))
		c.String(http.StatusUnauthorized, errInvalidCredentials.Error())
		return
	}
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to log in")
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		c.String(http.StatusUnauthorized, errInvalidCredentials.Error())
		return
	}

	if err := h.setAuthCookie(c, userID, username); err != nil {
		c.String(http.StatusInternalServerError, "failed to create session")
		return
	}

	c.String(http.StatusOK, "login successful")
}

// setAuthCookie signs a one-hour HS256 token for the user and stores it in an
// HttpOnly cookie. The token carries the user's id and name, never the
// password: a JWT payload is only base64-encoded, so anyone holding it can read it.
func (h *AuthHandler) setAuthCookie(c *gin.Context, userID int32, username string) error {
	now := time.Now()
	claims := authClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(int64(userID), 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenDuration)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(h.secret)
	if err != nil {
		return err
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		AuthCookieName,
		"Bearer "+token,
		int(tokenDuration.Seconds()),
		"/",
		"",
		gin.Mode() == gin.ReleaseMode,
		true,
	)
	return nil
}

func validateUsername(username string) error {
	if !usernamePattern.MatchString(username) {
		return errors.New("username must be 1-20 characters: letters, digits, '_', '.' or '-'")
	}
	return nil
}

func validateEmail(email string) error {
	if len(email) > maxEmailLen || !emailPattern.MatchString(email) {
		return errors.New("email must look like name@example.com")
	}
	return nil
}

func validatePassword(password string) error {
	if !passwordPattern.MatchString(password) {
		return errors.New("password must be 8-30 characters: letters, digits or symbols, no spaces")
	}
	return nil
}
