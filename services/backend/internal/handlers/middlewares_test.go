package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var testSecret = []byte("0123456789abcdef0123456789abcdef")

func sign(t *testing.T, method jwt.SigningMethod, key any, sub string, exp time.Time) string {
	t.Helper()
	claims := authClaims{
		Username: "alice",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	s, err := jwt.NewWithClaims(method, claims).SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestRequireAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewAuthHandler(nil, testSecret)

	r := gin.New()
	r.GET("/me", h.RequireAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"id": c.GetInt32(ContextUserID), "name": c.GetString(ContextUsername)})
	})

	valid := sign(t, jwt.SigningMethodHS256, testSecret, "7", time.Now().Add(time.Hour))
	expired := sign(t, jwt.SigningMethodHS256, testSecret, "7", time.Now().Add(-time.Minute))
	wrongKey := sign(t, jwt.SigningMethodHS256, []byte("another-secret-another-secret-xx"), "7", time.Now().Add(time.Hour))
	noneAlg := sign(t, jwt.SigningMethodNone, jwt.UnsafeAllowNoneSignatureType, "7", time.Now().Add(time.Hour))
	badSub := sign(t, jwt.SigningMethodHS256, testSecret, "abc", time.Now().Add(time.Hour))

	tests := []struct {
		name   string
		cookie string
		want   int
	}{
		{"valid token", "Bearer " + valid, http.StatusOK},
		{"no cookie", "", http.StatusUnauthorized},
		{"missing bearer prefix", valid, http.StatusUnauthorized},
		{"expired", "Bearer " + expired, http.StatusUnauthorized},
		{"wrong signing key", "Bearer " + wrongKey, http.StatusUnauthorized},
		{"alg none", "Bearer " + noneAlg, http.StatusUnauthorized},
		{"non-numeric subject", "Bearer " + badSub, http.StatusUnauthorized},
		{"garbage", "Bearer not.a.jwt", http.StatusUnauthorized},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/me", nil)
			if tt.cookie != "" {
				req.AddCookie(&http.Cookie{Name: AuthCookieName, Value: tt.cookie})
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tt.want {
				t.Fatalf("status = %d, want %d (body %s)", w.Code, tt.want, w.Body.String())
			}
		})
	}
}
