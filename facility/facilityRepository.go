package facility

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (f FacilityHandler) Get(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		slog.Error("userID is not exist")
		c.JSON(http.StatusForbidden, gin.H{"message": "userID is invalid"})
		return
	}
	facility, err := f.Facility.GetFacilityInfo(userID)
	if err != nil {
		slog.Error("can't get facility Info", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "500 Internal Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "getFacility success", "value": facility})
}
