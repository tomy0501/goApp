package facilityapi

import (
	"crypto/rsa"
	"encoding/json"
	"goapp/role"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type FacilityAPIJWT struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

type facilityAPIHandler struct {
	stockPemKey rsa.PrivateKey
	roleIF      role.RoleInterface
}

func NewFacilityAPIHandler(roleIF role.RoleInterface, stockPemKey rsa.PrivateKey) *facilityAPIHandler {
	return &facilityAPIHandler{
		stockPemKey: stockPemKey,
		roleIF:      roleIF,
	}
}

func (h *facilityAPIHandler) Get(g *gin.Context) {
	// Implementation for getting facility authorization
	userID := g.GetString("userID")
	if userID == "" {
		slog.Error("can't Get UserID")
		g.JSON(http.StatusUnauthorized, gin.H{"message": "please login again"})
		return
	}

	roles, err := h.roleIF.GetUserRoles(userID)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to get role"})
		return
	}
	jsonRoles, err := json.Marshal(roles)
	if err != nil {
		slog.Error("cannnot marshaling json", "error", err)
		g.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_SERVER_ERROR"})
		return
	}
	facilityRoleJWT, err := jwt.NewWithClaims(jwt.SigningMethodRS256, FacilityAPIJWT{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "authorizationService",
		},
		Role: string(jsonRoles),
	}).SignedString(&h.stockPemKey)
	if err != nil {
		slog.Error("failed to change to claim", "error", err)
		g.JSON(500, gin.H{"error": "Failed to generate JWT"})
		return
	}
	g.Header("X-AuthorizationToken", facilityRoleJWT)
	g.JSON(200, gin.H{"message": "ok"})
}
