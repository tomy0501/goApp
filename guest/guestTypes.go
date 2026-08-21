package guest

import (
	"goapp/role"
	"goapp/room"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Guest struct {
	GuestID   string    `json:"GuestID" binding:"required"`
	GuestName string    `json:"GuestName" binding:"required"`
	GuestAge  *int      `json:"GuestAge" binding:"gte=0,lte=130"`
	GuestSex  *string   `json:"GuestSex" binding:"oneof=male female other unknown"`
	GuestFrom *string   `json:"GuestFrom"`
	Room      room.Room `json:"Room" binding:"required"`
	Checkin   *string   `json:"Checkin" binding:"datetime=2006-01-02"`
	Checkout  *string   `json:"Checkout" binding:"datetime=2006-01-02"`
}

type GuestOnly struct {
	GuestID   string  `json:"GuestID" binding:"required"`
	GuestName string  `json:"GuestName" binding:"required"`
	GuestAge  *int    `json:"GuestAge" binding:"gte=0,lte=130"`
	GuestSex  *string `json:"GuestSex" binding:"oneof=male female other unknown"`
	GuestFrom *string `json:"GuestFrom"`
	Checkin   *string `json:"Checkin"`
}

type GuestInterface interface {
	GetGuestAll(facilityID string) ([]*GuestOnly, error)
	GetGuestById(facilityID string, guestID string) (*GuestOnly, error)
}

type guestHandler struct {
	GuestIF GuestInterface
	RoleIF  role.RoleInterface
}

type GetGuestUri struct {
	FacilityID string `uri:"facilityID" binding:"required,min=32,max=36"`
}

type GetGuestByIDUri struct {
	FacilityID string `uri:"facilityID" binding:"required,min=32,max=36"`
	GuestID    string `uri:"guestID" binding:"required,min=32,max=36"`
}

func NewGuestHandler(guestIF GuestInterface, roleIF role.RoleInterface) *guestHandler {
	return &guestHandler{
		GuestIF: guestIF,
		RoleIF:  roleIF,
	}
}

func (h *guestHandler) GetAll(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		slog.Error("can't Get UserID")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "please login again"})
		return
	}
	var facilityUri GetGuestUri

	if err := c.ShouldBindUri(&facilityUri); err != nil {
		slog.Error("uri is invalid")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	role, err := h.RoleIF.GetUserRoleByFacility(userID, facilityUri.FacilityID)
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
	guests, err := h.GuestIF.GetGuestAll(facilityUri.FacilityID)
	if err != nil {
		slog.Error("can't get guests", "guest error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, guests)
}

func (h *guestHandler) Get(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		slog.Error("can't Get UserID")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "please login again"})
		return
	}
	var guestUri GetGuestByIDUri
	if err := c.ShouldBindUri(&guestUri); err != nil {
		slog.Error("uri is invalid")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	role, err := h.RoleIF.GetUserRoleByFacility(userID, guestUri.FacilityID)
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
	guest, err := h.GuestIF.GetGuestById(guestUri.FacilityID, guestUri.GuestID)
	if err != nil {
		slog.Error("can't get guest", "guest error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK", "value": guest})
}
