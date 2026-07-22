package hotelindex

import (
	"goapp/role"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HotelIndex struct {
	HotelName string `json:"hotelName" binding:"require"`
	Authority string `json:"authority" binding:"require"`
}

type HotelIndexInterface interface {
	GetHotelIndexInfo(facilityID string, userID string) (*HotelIndex, error)
}
type hotelIndexHandler struct {
	HotelInterface HotelIndexInterface
	RoleInterfaces role.RoleInterface
}

func NewHotelIndexHandler(hotelIndexHandle HotelIndexInterface, roleHandler role.RoleInterface) *hotelIndexHandler {
	return &hotelIndexHandler{
		HotelInterface: hotelIndexHandle,
		RoleInterfaces: roleHandler,
	}
}

type HotelGetUri struct {
	FacilityID string `uri:"facilityID" binding:"required,min=32,max=36"`
}

func (h *hotelIndexHandler) Get(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		slog.Error("can't Get userID")
		c.JSON(http.StatusForbidden, gin.H{"message": "invalid data"})
		return
	}
	var hotelUri HotelGetUri
	if err := c.ShouldBindUri(&hotelUri); err != nil {
		slog.Error("invalid facility Uri")
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request, stop process"})
		return
	}
	role, err := h.RoleInterfaces.GetUserRoleByFacility(userID, hotelUri.FacilityID)
	if err != nil {
		slog.Error("can't get role", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"message": "request is invalid"})
		return
	}
	if role.Role != "staff" {
		c.JSON(http.StatusForbidden, gin.H{"message": "sorry you have not permission to access reservation"})
		return
	}
	hotelIndex, err := h.HotelInterface.GetHotelIndexInfo(hotelUri.FacilityID, userID)
	if err != nil {
		slog.Error("can't get hotelIndex", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK", "value": hotelIndex})
}
