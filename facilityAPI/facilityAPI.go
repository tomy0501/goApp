package facilityapi

import (
	"crypto/rsa"
	"goapp/role"
	"os"
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
	userID := g.Request.Header.Get("userID")
	facilityID := g.Request.Header.Get("facilityID")
	if userID == "" || facilityID == "" {
		g.JSON(400, gin.H{"error": "userID and facilityID are required"})
		return
	}
	APIKey := g.Request.Header.Get("APIKey")
	if APIKey == "" {
		g.JSON(400, gin.H{"error": "APIKey is required"})
		return
	}
	APIKeyFromEnv := os.Getenv("AUTHORIZATION_API_KEY")
	if APIKey != APIKeyFromEnv {
		g.JSON(401, gin.H{"error": "Invalid APIKey"})
		return
	}
	role, err := h.roleIF.GetUserRoleByFacility(userID, facilityID)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to get role"})
		return
	}
	facilityRoleJWT, err := jwt.NewWithClaims(jwt.SigningMethodRS256, FacilityAPIJWT{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "authorizationService",
		},
		Role: role.Role.GetRoleName(),
	}).SignedString(&h.stockPemKey)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to generate JWT"})
		return
	}
	g.JSON(200, gin.H{"facilityRoleJWT": facilityRoleJWT})
}
