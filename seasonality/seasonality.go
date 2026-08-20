package seasonality

import (
	"goapp/role"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type Seasonality struct {
	ID              string `json:"id"`
	SeasonalityName string `json:"seasonalityName"`
}

type SeasonalityPrice struct {
	ID             string          `json:"id"`
	Price          decimal.Decimal `json:"price"`
	NumberOfPeople int             `json:"numberOfPeople"`
}

type GetSeasonalityInterface interface {
	GetSeasonalityByFacility(facilityID string) ([]*Seasonality, error)
	GetSeasonalityNumberOfPeopleBySeasonality(seasonalityID string) ([]*SeasonalityPrice, error)
}

type getSeasonalityHandler struct {
	SeasonalityHandler GetSeasonalityInterface
	RoleHandler        role.RoleInterface
}

func NewSeasonalityHandler(seasonalityIF GetSeasonalityInterface, roleIF role.RoleInterface) getSeasonalityHandler {
	return getSeasonalityHandler{
		SeasonalityHandler: seasonalityIF,
		RoleHandler:        roleIF,
	}
}

type SeasonalityUri struct {
	FacilityID string `uri:"facilityID" binding:"required,min=32,max=36"`
}

type SeasonalityPriceUri struct {
	FacilityID    string `uri:"facilityID" binding:"required,min=32,max=36"`
	SeasonalityID string `uri:"seasonalityID" binding:"required,min=32,max=36"`
}

func (h *getSeasonalityHandler) GetSeasonality(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		slog.Error("not found User ID")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid Data"})
		return
	}

	var seasonalityUri SeasonalityUri

	if err := c.ShouldBindUri(&seasonalityUri); err != nil {
		slog.Error("url is invalid", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal Sever error"})
		return
	}
	role, err := h.RoleHandler.GetUserRoleByFacility(userID, seasonalityUri.FacilityID)
	if err != nil {
		slog.Error("can't get role", "role error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	if !role.Role.CanAccessReadRole() {
		slog.Error("can't access role", "role", role.Role, "errorID", userID)
		c.JSON(http.StatusForbidden, gin.H{"message": "invalid status"})
		return
	}
	seasonality, err := h.SeasonalityHandler.GetSeasonalityByFacility(seasonalityUri.FacilityID)
	if err != nil {
		slog.Error("can't Get seasonality", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK", "value": seasonality})
}

func (h *getSeasonalityHandler) GetSeasonalityPrice(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		slog.Error("not found User ID")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid Data"})
		return
	}
	var seasonalityPriceUri SeasonalityPriceUri
	if err := c.ShouldBindUri(&seasonalityPriceUri); err != nil {
		slog.Error("url is invalid", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal Sever error"})
		return
	}
	role, err := h.RoleHandler.GetUserRoleByFacility(userID, seasonalityPriceUri.FacilityID)
	if err != nil {
		slog.Error("cannot access role", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "authorize failed"})
		return
	}
	if !role.Role.CanAccessReadRole() {
		slog.Error("can't access role", "role", role.Role, "errorID", userID)
		c.JSON(http.StatusForbidden, gin.H{"message": "invalid status"})
		return
	}
	seasonalityPrices, err := h.SeasonalityHandler.GetSeasonalityNumberOfPeopleBySeasonality(seasonalityPriceUri.SeasonalityID)
	if err != nil {
		slog.Error("cannot access seasonalityPrice", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK", "value": seasonalityPrices})
}
