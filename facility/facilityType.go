package facility

type FacilityHandler struct {
	Facility FacilityInterface
}

func NewFacilityHandler(f FacilityInterface) FacilityHandler {
	return FacilityHandler{
		Facility: f,
	}
}

type FacilityInterface interface {
	GetFacilityInfo(userID string) (*[]FacilityJSON, error)
}

type FacilityJSON struct {
	FacilityName string `json:"hotelName"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
	FacilityID   string `json:"hotelID"`
}
