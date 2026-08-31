package price

import (
	"goapp/db"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type priceRepository struct {
	db *db.UserDB
}

func NewPriceRepository(db *db.UserDB) GetSeasonalityInterface {
	return &priceRepository{
		db: db,
	}
}

type CombineSeasonalityData struct {
	RoomName        string          `db:"room_name"`
	RoomTypeName    string          `db:"room_type_name"`
	SeasonalityName string          `db:"seasonality_name"`
	NumberOfPeople  int             `db:"number_of_people"`
	Price           decimal.Decimal `db:"price"`
	Date            string          `db:"date"`
}

func (r *priceRepository) GetSeasonalityByMonth(yearMonth string, facilityID string, roomTypeID string) ([]*SeasonalityDate, error) {
	parsedFacilityID, err := uuid.Parse(facilityID)
	if err != nil {
		return nil, err
	}
	parsedRoomTypeID, err := uuid.Parse(roomTypeID)
	if err != nil {
		return nil, err
	}
	parsedTime, err := time.Parse("2006-01", yearMonth)
	if err != nil {
		return nil, err
	}
	year := parsedTime.Year()
	month := parsedTime.Month()

	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	lastDay := time.Date(year, month+1, 1, 0, 0, 0, 0, time.Local).AddDate(0, 0, -1)

	sql := `
	SELECT rt.room_type_name,s.seasonality_name,snp.number_of_people,snp.price,sod.date FROM facilities f 
	INNER JOIN room_type rt ON rt.facility_id = f.id
	INNER JOIN seasonality s ON f.id = s.facility_id
	INNER JOIN seasonality_number_of_people snp ON s.id = snp.seasonality_id
	INNER JOIN seasonality_of_date sod ON s.id = sod.seasonality_id
	WHERE rt.facility_id = ? AND rt.id = ? AND sod.date BETWEEN ? AND ? ORDER BY sod.date,snp.number_of_people
	`

	var rawSeasonalityDatas []CombineSeasonalityData

	err = r.db.UserDB.Select(&rawSeasonalityDatas, sql, parsedFacilityID[:], parsedRoomTypeID[:], firstDay, lastDay)
	if err != nil {
		return nil, err
	}
	var preDate string
	isFirst := true
	var seasonalityDates []*SeasonalityDate
	var seasonalityPrices []SeasonalityPrice
	for _, row := range rawSeasonalityDatas {
		if isFirst {
			preDate = row.Date
			isFirst = false
		}
		if preDate != row.Date {
			copySeasonalityDate := make([]SeasonalityPrice, len(seasonalityPrices))
			copy(copySeasonalityDate, seasonalityPrices)
			seasonalityDate := &SeasonalityDate{
				Date:        preDate,
				Seasonality: copySeasonalityDate,
			}
			seasonalityDates = append(seasonalityDates, seasonalityDate)
			seasonalityPrices = seasonalityPrices[:0]
			preDate = row.Date
		} else {
			sesonalityPrice := SeasonalityPrice{
				Price:          row.Price,
				NumberOfPeople: row.NumberOfPeople,
			}
			seasonalityPrices = append(seasonalityPrices, sesonalityPrice)
		}
	}
	return seasonalityDates, nil
}

func (r *priceRepository) GetSeasonalityByDate(date string, facilityID string, roomTypeID string) (*SeasonalityDate, error) {
	parsedFacilityID, err := uuid.Parse(facilityID)
	if err != nil {
		return nil, err
	}
	parsedRoomTypeID, err := uuid.Parse(roomTypeID)
	if err != nil {
		return nil, err
	}

	sql := `
	SELECT rt.room_type_name,s.seasonality_name,snp.number_of_people,snp.price,sod.date FROM facilities f 
	INNER JOIN room_type rt ON  rt.facility_id = f.id
	INNER JOIN seasonality s ON f.id = s.facility_id
	INNER JOIN seasonality_number_of_people snp ON s.id = snp.seasonality_id
	INNER JOIN seasonality_of_date sod ON s.id = sod.seasonality_id
	WHERE f.facility_id = ? AND rt.id = ? AND sod.date = ?
	`
	var CombineSeasonality []CombineSeasonalityData
	err = r.db.UserDB.Select(&CombineSeasonality, sql, parsedFacilityID[:], parsedRoomTypeID[:], date)
	if err != nil {
		return nil, err
	}
	var Date string
	var seasonality []SeasonalityPrice
	isFirst := true
	for _, row := range CombineSeasonality {
		if isFirst {
			Date = row.Date
		}
		seasonalityPrice := SeasonalityPrice{
			NumberOfPeople: row.NumberOfPeople,
			Price:          row.Price,
		}
		seasonality = append(seasonality, seasonalityPrice)
	}
	seasonalityDate := &SeasonalityDate{
		Date:        Date,
		Seasonality: seasonality,
	}
	return seasonalityDate, nil
}
