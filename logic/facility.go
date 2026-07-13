package logic

import (
	"goapp/myTypes"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (d *UserDB) GetFacilityInfo(userID string) (*[]myTypes.FacilityJSON, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		slog.Error("uuid can't parse", "error", err)
		return nil, err
	}
	row, err := d.userDB.Query(`SELECT 
	facilities.id,facilities.name,facilities.address,facilities.phone 
	FROM facilities INNER JOIN facility_users 
	ON facilities.id = facility_users.facility_id 
	WHERE facility_users.user_id = ?;`, parsedUserID[:])
	if err != nil {
		slog.Error("failed in get facility", "error", err)
		return nil, err
	}
	err = row.Err()
	if err != nil {
		slog.Error("failed when fetch facility row", "error", err)
		return nil, err
	}
	defer row.Close()
	var facilities []myTypes.FacilityJSON
	for row.Next() {
		var facility myTypes.FacilityJSON
		var rawFacilityID []byte
		if err := row.Scan(&rawFacilityID, &facility.FacilityName, &facility.Address, &facility.Phone); err != nil {
			slog.Error("can't read facility row", "error", err)
			return nil, err
		}
		parsedFacilityID, err := uuid.FromBytes(rawFacilityID)
		if err != nil {
			slog.Error("can't read faclity ID", "error", err)
			return nil, err
		}
		facility.FacilityID = parsedFacilityID.String()
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
