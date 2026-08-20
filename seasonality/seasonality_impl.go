package seasonality

import (
	"goapp/db"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type seasonalityRepository struct {
	db *db.UserDB
}

func NewSeasonalityRepository(db *db.UserDB) GetSeasonalityInterface {
	return &seasonalityRepository{
		db: db,
	}
}

type seasonalityValue struct {
	SeasonalityID   string `db:"id"`
	SeasonalityName string `db:"seasonality_name"`
}

type seasonalityPrieValue struct {
	SeasonalityNumberID       string          `db:"id"`
	SeasonalityNumberOfPeople int             `db:"number_of_people"`
	SeasonalityPrice          decimal.Decimal `db:"price"`
}

func (r *seasonalityRepository) GetSeasonalityByFacility(facilityID string) ([]*Seasonality, error) {
	parsedFacilityID, err := uuid.Parse(facilityID)
	if err != nil {
		return nil, err
	}
	sql := `
		SELECT id,seasonality_name FROM seasonality WHERE facility_id = ?
	`

	var seasonalityValues []seasonalityValue

	err = r.db.UserDB.Select(&seasonalityValues, sql, parsedFacilityID[:])
	if err != nil {
		return nil, err
	}
	var seaonalities []*Seasonality
	for _, row := range seasonalityValues {
		parsedSeasonalityID, err := uuid.FromBytes([]byte(row.SeasonalityID))
		if err != nil {
			return nil, err
		}
		seasonality := &Seasonality{
			ID:              parsedSeasonalityID.String(),
			SeasonalityName: row.SeasonalityName,
		}
		seaonalities = append(seaonalities, seasonality)
	}
	return seaonalities, nil
}

func (r *seasonalityRepository) GetSeasonalityNumberOfPeopleBySeasonality(seasonalityID string) ([]*SeasonalityPrice, error) {
	parsedSeasonalityID, err := uuid.Parse(seasonalityID)
	if err != nil {
		return nil, err
	}
	sql :=
		`
		SELECT id,number_of_people,price FROM seasonality_number_of_people WHERE seasonality_id = ?
	`

	var seasonalityPricesValues []seasonalityPrieValue

	err = r.db.UserDB.Select(&seasonalityPricesValues, sql, parsedSeasonalityID[:])
	if err != nil {
		return nil, err
	}
	var seasonalityPrices []*SeasonalityPrice
	for _, row := range seasonalityPricesValues {
		seasonalityPriceParsedID, err := uuid.FromBytes([]byte(row.SeasonalityNumberID))
		if err != nil {
			return nil, err
		}
		seasonalityPrice := &SeasonalityPrice{
			ID:             seasonalityPriceParsedID.String(),
			NumberOfPeople: row.SeasonalityNumberOfPeople,
			Price:          row.SeasonalityPrice,
		}
		seasonalityPrices = append(seasonalityPrices, seasonalityPrice)
	}
	return seasonalityPrices, nil
}
