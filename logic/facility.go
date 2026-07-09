package logic

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FacilityInterface interface {
	GetFacilityInfo(userID string) (*[]FacilityJSON, error)
}

type FacilityJSON struct {
	FacilityName string `json:"hotelName"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
	FacilityID   string `json:"hotelID"`
}

func (d *UserDB) GetFacilityInfo(userID string) (*[]FacilityJSON, error) {
	uuid, err := uuid.Parse(userID)
	if err != nil {
		slog.Error("uuid can't parse", "error", err)
		return nil, err
	}
	row, err := d.userDB.Query(`SELECT 
	facilities.id,facilities.name,facilities.address,facilities.phone 
	FROM facilities INNER JOIN facility_users 
	ON facilities.id = facility_users.facility_id 
	WHERE facility_users.user_id = ?;`, uuid[:])
	if err != nil {
		slog.Error("failed in get facility", "error", err)
		return nil, err
	}
	var facilities []FacilityJSON
	for row.Next() {
		var facility FacilityJSON
		if err := row.Scan(&facility.FacilityID, &facility.FacilityName, &facility.Address, &facility.Phone); err != nil {
			slog.Error("can't read facility row", "error", err)
			return nil, err
		}
		facilities = append(facilities, facility)
	}
	return &facilities, nil
}

func (h *AuthHandlers) FacilityGetLogic(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		slog.Error("userID is not exist")
		c.JSON(http.StatusForbidden, gin.H{"message": "userID is invalid"})
		return
	}
	facility, err := h.AuthStorage.GetFacilityInfo(userID)
	if err != nil {
		slog.Error("can't get facility Info", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "500 Internal Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "getFacilit success", "value": facility})
}
