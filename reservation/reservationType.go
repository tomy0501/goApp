package reservation

import (
	"goapp/guest"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Reservation struct {
	ReservationID  string        `json:"reservationID" binding:"required,min=32,max=36"`
	RepresentGuest *guest.Guest  `json:"representGuest" binding:"required"`
	Guests         []guest.Guest `json:"guests" binding:"required"`
}

type ReservationHandler struct {
	Reservation ReservationInterface
}

type ReservationInterface interface {
	GetAllReservation(facilityID string) ([]*Reservation, error)
	// GetReservationById(facilityID string, reservationID string) (Reservation, error)
	// PutReservation(facilityID string, reservation Reservation) error
	// PutManyReservation(facilityID string, reservations []Reservation) error
	// PostReservation(facilityID string, reservationID string) error
	// DeleteReservation(facilityID string) error
}

type GetReservationJSON struct {
	FacilityID string `json:"facility_id" binding:"required,min=32,max=36"`
}

func (h *ReservationHandler) GetAll(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		slog.Error("userID is not exist")
		c.JSON(http.StatusForbidden, gin.H{"message": "userID is invalid"})
		return
	}
	var reservationJSON GetReservationJSON
	err := c.ShouldBindJSON(&reservationJSON)
	if err != nil {
		slog.Error("failed in parse json", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error happen"})
		return
	}
	reservation, err := h.Reservation.GetAllReservation(reservationJSON.FacilityID)
	if err != nil {
		slog.Error("can't get reservation Info", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "500 Internal Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success to get reservation", "value": reservation})
}
