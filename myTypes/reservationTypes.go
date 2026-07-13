package myTypes

type Reservation struct {
	ReservationID  string  `json:"reservationID" binding:"required,min=32,max=36"`
	RepresentGuest Guest   `json:"representGuest" binding:"required"`
	Guests         []Guest `json:"guests" binding:"required"`
}

type BaseGetReservationInfo struct {
	FacilityID string `json:"facilityID" binding:"required,min=32,max=36"`
}

type GetReservationByReservationID struct {
	facilityInfo  BaseGetReservationInfo
	ReservationID string `json:"reservationID" binding:"required,min=32,max=36"`
}

type ReservationInterface interface {
	GetAllReservationInfo(facilityID string) ([]Reservation, error)
	GetReservationById(userID string, facilityID string, reservationID string) (Reservation, error)
	CreateReservation(userID string, facilityID string, reservation Reservation) error
	BulkCreateReservation(userID string, facilityID string, reservations []Reservation) error
	UpdateReservationById(userID string, facilityID string, reservationID string) error
}
