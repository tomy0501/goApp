package authorize

import (
	"errors"
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

var skippaths = map[string]bool{
	"/facilityAPI": true,
}

func CheckAccessToken(h AuthHandlerInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Info("check path", "path", c.Request)
		if skippaths[c.Request.URL.Path] {
			slog.Info("skip access token because path is skiptarget", "path", c.Request.URL)
			c.Next()
			return
		}
		accessToken, err := c.Cookie("accessToken")
		if err != nil {
			if errors.Is(http.ErrNoCookie, err) {
				accessToken = c.Request.Header["Authorization"][0]
				if accessToken == "" {
					slog.Error("cannnot found access token")
					c.JSON(http.StatusUnauthorized, gin.H{"message": "accessToken is invalid", "code": "TOKEN_INVALID"})
					c.Abort()
					return
				}
			} else {
				slog.Error("unknown error happen", "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error", "code": "INTERNAL_SERVER_ERROR"})
				c.Abort()
				return
			}
		}
		var claims CustomClaims
		token, err := jwt.ParseWithClaims(accessToken, &claims, func(tokenString *jwt.Token) (interface{}, error) {
			if _, ok := tokenString.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, jwt.ErrSignatureInvalid
			}

			return h.GetPublicKey(), nil
		})

		if err != nil || !token.Valid {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "token expired", "code": "TOKEN_EXPIRED"})
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
