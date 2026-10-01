package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Keys under which RequireAuth stores the caller's identity in the gin context.
const (
	ContextUserID   = "user_id"
	ContextUsername = "username"
)

// RequireAuth rejects requests that don't carry a valid, unexpired token in the
// "Authorization: Bearer <jwt>" cookie set by signup and login. On success it
// stores the user's id and name in the context (ContextUserID as int32,
// ContextUsername as string) for the handlers that follow.
func (h *AuthHandler) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := h.parseAuthCookie(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or missing token"})
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 32)
		if err != nil || userID <= 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or missing token"})
			return
		}

		c.Set(ContextUserID, int32(userID))
		c.Set(ContextUsername, claims.Username)
		c.Next()
	}
}

func (h *AuthHandler) parseAuthCookie(c *gin.Context) (*authClaims, error) {
	cookie, err := c.Cookie(AuthCookieName)
	if err != nil {
		return nil, err
	}

	raw, ok := strings.CutPrefix(cookie, "Bearer ")
	if !ok || raw == "" {
		return nil, errors.New("missing bearer prefix")
	}

	claims := &authClaims{}
	token, err := jwt.ParseWithClaims(raw, claims, func(*jwt.Token) (any, error) {
		return h.secret, nil
	},
		// Pin the algorithm so a token can't pick its own ("none", RS256, ...).
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
