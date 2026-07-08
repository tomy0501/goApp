package authorize

import (
	"crypto/rsa"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}

type AuthHandlerInterface interface {
	GetPublicKey() string
}

type AuthHandler struct {
	PublicKey *rsa.PublicKey
}

func (h *AuthHandler) GetPublicKey() *rsa.PublicKey {
	return h.PublicKey
}
func CheckAccessToken(h *AuthHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("accessToken")
		if err != nil {
			slog.Error("not found accessToken")
			c.JSON(http.StatusForbidden, gin.H{"message": "accessToken is invalid"})
			c.Abort()
			return
		}
		var claims CustomClaims
		token, err := jwt.ParseWithClaims(accessToken, &claims, func(tokenString *jwt.Token) (interface{}, error) {
			if _, ok := tokenString.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, jwt.ErrSignatureInvalid
			}

			return h.GetPublicKey(), nil
		})

		if err != nil || !token.Valid {
			if err == jwt.ErrTokenExpired && token.Valid {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "toyken expired"})
				c.Abort()
				return
			} else {
				c.JSON(http.StatusForbidden, gin.H{"message": "token is invalid"})
				c.Abort()
				return
			}
		}
		userID, err := token.Claims.GetSubject()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error happen"})
			c.Abort()
			return
		}
		Role := token.Claims.(*CustomClaims).Role
		c.Set("userID", userID)
		c.Set("Role", Role)
		c.Next()
	}
}
