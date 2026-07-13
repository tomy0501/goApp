package logic

import (
	"goapp/myTypes"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandlers) GetAllReservation(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		slog.Error("userID is not exist")
		c.JSON(http.StatusForbidden, gin.H{"message": "userID is invalid"})
		return
	}
	var reservationJSON myTypes.BaseGetReservationInfo
	err := c.ShouldBindJSON(&reservationJSON)
	if err != nil {
		slog.Error("failed in parse json", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error happen"})
		return
	}
	reservation, err := h.AuthStorage.GetAllReservationInfo(reservationJSON.FacilityID)
	if err != nil {
		slog.Error("can't get reservation Info", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "500 Internal Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success to get reservation", "value": reservation})
}
