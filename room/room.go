package room

import (
	"goapp/role"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Room struct {
	RoomID          string  `json:"RoomID" binding:"required"`
	RoomName        string  `json:"RoomName" binding:"required"`
	RoomTypeID      string  `json:"RoomTypeID" binding:"required,min=32,max=36"`
	RoomType        string  `json:"RoomType"`
	SeasonalityID   *string `json:"SeasonalityID"`
	SeasonalityName *string `json:"SeasonalityName"`
}

type GetRoomInterface interface {
	GetRoomAll(facilityID string) ([]*Room, error)
	GetRoomById(facilityID string, roomID string) (*Room, error)
}

type getRoomHandler struct {
	RoomHandler GetRoomInterface
	RoleHandler role.RoleInterface
}

type GetRoomUri struct {
	FacilityID string `uri:"facilityID" binding:"required,min=32,max=36"`
}

type GetRoomByIDUri struct {
	FacilityID string `uri:"facilityID" binding:"required,min=32,max=36"`
	RoomID     string `uri:"roomID" binding:"required,min=32,max=36"`
}

func NewGetRoomHandler(roomHandler GetRoomInterface, roleHandler role.RoleInterface) *getRoomHandler {
	return &getRoomHandler{
		RoomHandler: roomHandler,
		RoleHandler: roleHandler,
	}
}

func (h *getRoomHandler) GetAll(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		slog.Error("can't Get UserID")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "please login again"})
		return
	}
	var facilityUri GetRoomUri

	if err := c.ShouldBindUri(&facilityUri); err != nil {
		slog.Error("uri is invalid")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	role, err := h.RoleHandler.GetUserRoleByFacility(userID, facilityUri.FacilityID)
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
	rooms, err := h.RoomHandler.GetRoomAll(facilityUri.FacilityID)
	if err != nil {
		slog.Error("can't access room", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "can't get room"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "get room!", "value": rooms})
}

func (h *getRoomHandler) Get(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		slog.Error("can't Get UserID")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "please login again"})
		return
	}
	var facilityUri GetRoomByIDUri

	if err := c.ShouldBindUri(&facilityUri); err != nil {
		slog.Error("uri is invalid")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	role, err := h.RoleHandler.GetUserRoleByFacility(userID, facilityUri.FacilityID)
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
	room, err := h.RoomHandler.GetRoomById(facilityUri.FacilityID, facilityUri.RoomID)
	if err != nil {
		slog.Error("can't access room", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "can't get room"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK", "value": room})
}
