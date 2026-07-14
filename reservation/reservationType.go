package reservation

import (
	"goapp/guest"
	"goapp/role"
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
	Roll        role.RoleInterface
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

type GetReservationUri struct {
	FacilityID string `uri:"facilityID" binding:"required,min=32,max=36"`
}

func (h *ReservationHandler) GetAll(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		slog.Error("userID is not exist")
		c.JSON(http.StatusForbidden, gin.H{"message": "userID is invalid"})
		return
	}
	var facilityRequest GetReservationUri
	if err := c.ShouldBindUri(&facilityRequest); err != nil {
		slog.Error("invalid facility Uri")
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request, stop process"})
		return
	}
	role, err := h.Roll.GetUserRoleByFacility(userID, facilityRequest.FacilityID)
	if err != nil {
		slog.Error("can't get role", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"message": "request is invalid"})
		return
	}
	if role.Role != "staff" {
		c.JSON(http.StatusForbidden, gin.H{"message": "sorry you have not permission to access reservation"})
		return
	}
	reservation, err := h.Reservation.GetAllReservation(facilityRequest.FacilityID)
	if err != nil {
		slog.Error("can't get reservation Info", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "500 Internal Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success to get reservation", "value": reservation})
}
