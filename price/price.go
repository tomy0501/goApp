package price

import (
	"goapp/role"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type SeasonalityPrice struct {
	Price          decimal.Decimal `json:"price" binding:"required gte=0"`
	NumberOfPeople int             `json:"number_of_people" binding:"required gte=0"`
}

type SeasonalityDate struct {
	Date        string             `json:"date" binding:"datetime=2006-01-02"`
	Seasonality []SeasonalityPrice `json:"seasonality"`
}

type GetSeasonalityInterface interface {
	GetSeasonalityByMonth(month string, facilityID string, roomID string) ([]*SeasonalityDate, error)
	GetSeasonalityByDate(date string, facilityID string, roomID string) (*SeasonalityDate, error)
}

type priceHandler struct {
	SeasonalityIF GetSeasonalityInterface
	RoleIF        role.RoleInterface
}

func NewPriceHandler(seasonalityIF GetSeasonalityInterface, roleIF role.RoleInterface) *priceHandler {
	return &priceHandler{
		SeasonalityIF: seasonalityIF,
		RoleIF:        roleIF,
	}
}

type GetSeasonalityByMonthUri struct {
	Month         string `uri:"month" binding:"datetime=2006-01"`
	FacilityID    string `uri:"facilityID" binding:"required,min=32,max=36"`
	SeasonalityID string `uri:"seasonalityID" binding:"required,min=32,max=36"`
}

type GetSeasonalityByDateUri struct {
	Date          string `uri:"month" binding:"datetime=2006-01"`
	FacilityID    string `uri:"facilityID" binding:"required,min=32,max=36"`
	SeasonalityID string `uri:"seasonalityID" binding:"required,min=32,max=36"`
}

func (h *priceHandler) GetSeasonalityByMonth(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		slog.Error("can't Get UserID")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "please login again"})
		return
	}
	var priceByMonthUri GetSeasonalityByMonthUri
	if err := c.ShouldBindUri(&priceByMonthUri); err != nil {
		slog.Error("can't get uri", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "invalid uri status"})
		return
	}
	roles, err := h.RoleIF.GetUserRoleByFacility(userID, priceByMonthUri.FacilityID)
	if err != nil {
		slog.Error("can't get role", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "invalid uri status"})
		return
	}

	if !roles.Role.CanAccessReadRole() {
		slog.Error("access deny beacause of can't permitted")
		c.JSON(http.StatusForbidden, gin.H{"message": "you cannot access"})
		return
	}
	prices, err := h.SeasonalityIF.GetSeasonalityByMonth(priceByMonthUri.Month, priceByMonthUri.FacilityID, priceByMonthUri.SeasonalityID)
	if err != nil {
		slog.Error("can't get seasonality", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "getSeasonality", "value": prices})
}
