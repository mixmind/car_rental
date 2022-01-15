package cmds

import (
	"car-rental/internal/server/db"
	"car-rental/internal/server/domain"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type RentProcessor struct {
	dbStruct *db.DBStruct
}

func NewRentProcessor(dbStruct *db.DBStruct) *RentProcessor {
	return &RentProcessor{dbStruct: dbStruct}
}

/*
Insert rent into DB
*/
func (rentPr *RentProcessor) InsertRentInDB(rent domain.RentInfo, car domain.Car) (int64, error) {
	tx, err := rentPr.dbStruct.BeginTransaction()
	if err != nil {
		return 0, errors.Wrap(err, "Failed to start a transaction")
	}
	statement, err := tx.Prepare(db.InsertIntoRentTable)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to prepare stmt")
	}
	if len(rent.FromDate) == 0 || rent.CarID == 0 || len(rent.ToDate) == 0 {
		return 0, fmt.Errorf("Rent dates and car ID should be provided")
	}
	conditionerText := "With Air Conditioner"
	if !car.AirConditioner {
		conditionerText = "Without Air Conditioner"
	}
	if !checkCarProps(rent, car) {
		return 0, fmt.Errorf("Some of new rent props are incorrect. Please check them again!")
	}
	if isExists, err := rentPr.checkCarAvailability(rent, car); err != nil || isExists {
		if err != nil {
			log.Error(err)
		}
		return 0, fmt.Errorf("Car is not available in such dates")
	}
	res, err := statement.Exec(rent.CarID,
		rent.FromDate,
		rent.ToDate,
		rent.Location,
		strings.Join(rent.AvailableExtras, ","),
		strings.Join(rent.Discounts, ","),
		fmt.Sprintf(`%s %s.Part of %d group. With %d doors, %d adult places, %d big luggage and %d small luggage places.%s. For drivers with minimal age %d`,
			car.CarCompanyName,
			car.Description,
			car.CarGroup,
			car.Doors,
			car.AdultPlaces,
			car.BigLuggage,
			car.SmallLuggage,
			conditionerText,
			car.MinimumAge),
	)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to execute a prepared statement")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "Failed to execute a extract last id")
	}
	if err := tx.Commit(); err != nil {
		return 0, errors.Wrap(err, "Failed to commit a transaction")
	}

	return id, nil
}

/*
Get rents from DB
*/
func (rentPr *RentProcessor) GetRentsFromDB() ([]domain.RentInfo, error) {
	rows, err := rentPr.dbStruct.Query(db.SelectRents)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute a sql query")
	}

	var result []domain.RentInfo
	for rows.Next() {
		var receivedRow domain.RentInfo
		var extras string
		var discounts string

		err = rows.Scan(
			&receivedRow.RentID,
			&receivedRow.CarID,
			&receivedRow.FromDate,
			&receivedRow.ToDate,
			&receivedRow.Location,
			&extras,
			&discounts,
			&receivedRow.CarDetails,
		)
		if err != nil {
			log.Error(err)
			continue
		}
		receivedRow.AvailableExtras = strings.Split(extras, ",")
		receivedRow.Discounts = strings.Split(discounts, ",")
		result = append(result, receivedRow)
	}
	return result, nil
}

/*
Get rent from DB
*/
func (rentPr *RentProcessor) GetRentFromDB(rentID int) (*domain.RentInfo, error) {
	stmt, err := rentPr.dbStruct.Prepare(fmt.Sprintf("%s WHERE rent_id=?", db.SelectRents))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to prepare an sql query")
	}
	row := stmt.QueryRow(rentID)
	var receivedRow domain.RentInfo
	var extras string
	var discounts string

	err = row.Scan(
		&receivedRow.RentID,
		&receivedRow.CarID,
		&receivedRow.FromDate,
		&receivedRow.ToDate,
		&receivedRow.Location,
		&extras,
		&discounts,
		&receivedRow.CarDetails,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to query rent")
	}
	receivedRow.AvailableExtras = strings.Split(extras, ",")
	receivedRow.Discounts = strings.Split(discounts, ",")

	return &receivedRow, nil
}

/*
Remove rent from DB
*/
func (rentPr *RentProcessor) RemoveRentFromDB(rentID int) (int64, error) {
	stmt, err := rentPr.dbStruct.Prepare(db.RemoveRent)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to prepare an sql query")
	}

	res, err := stmt.Exec(rentID)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to execute rent delete")
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "Failed to extract rows deleted number")
	}

	return affect, nil
}

/*
Check if car information is correct in provided rent data
*/
func checkCarProps(rent domain.RentInfo, car domain.Car) bool {
	locationExists := false
	groupIsOK := false
	ageIsOK := false
	for _, loc := range car.AvailableLocations {
		if loc == rent.Location {
			locationExists = true
			break
		}
	}
	if len(rent.AgeGroup) > 0 {
		failedConversation := false
		splittedAge := strings.Split(rent.AgeGroup, "-")
		parsedMin, err := strconv.Atoi(splittedAge[0])
		if err != nil {
			log.Error(err)
			failedConversation = true
		}
		if !failedConversation {
			if len(splittedAge) > 1 {
				parsedMax, err := strconv.Atoi(splittedAge[1])
				if err != nil {
					log.Error(err)
					failedConversation = true
				}
				if parsedMin < parsedMax && car.MinimumAge > parsedMin && car.MinimumAge < parsedMax {
					ageIsOK = true
				} else {
					log.Error("Please provide correct age interval")
				}
			} else {
				if parsedMin > car.MinimumAge {
					ageIsOK = true
				}
			}
		}

	}
	if rent.CarGroup == car.CarGroup {
		groupIsOK = true
	} else {
		log.Error("Please provide correct car group")
	}

	if !locationExists {
		log.Error("Please provide correct location for this car")
	}
	if !ageIsOK {
		log.Error("Please provide correct age interval")
	}
	return locationExists && groupIsOK && ageIsOK
}

/*
Check if car still available
*/
func (rentPr *RentProcessor) checkCarAvailability(rent domain.RentInfo, car domain.Car) (bool, error) {
	searchParams := buildFromToFilter(rent.FromDate, rent.ToDate, true)
	query := db.SelectCarsRents + " WHERE " + searchParams.DateFilter
	log.Debugln(query)
	rows, err := rentPr.dbStruct.Query(query)
	if err != nil {
		return false, errors.Wrap(err, "Failed to execute a sql query")
	}

	var result []domain.CombinedRentInfo
	for rows.Next() {
		var receivedRow domain.CombinedRentInfo
		var locations string
		var rentID sql.NullInt32
		var fromDate sql.NullString
		var toDate sql.NullString
		var location sql.NullString
		var extras sql.NullString
		var discounts sql.NullString
		var carDetails sql.NullString
		err = rows.Scan(&receivedRow.CarID,
			&receivedRow.CarCompanyName,
			&receivedRow.Doors,
			&receivedRow.BigLuggage,
			&receivedRow.SmallLuggage,
			&receivedRow.AdultPlaces,
			&receivedRow.AirConditioner,
			&receivedRow.MinimumAge,
			&locations,
			&receivedRow.CarGroup,
			&receivedRow.Description,
			&receivedRow.Price,
			&rentID,
			&fromDate,
			&toDate,
			&location,
			&extras,
			&discounts,
			&carDetails,
		)
		if err != nil {
			log.Error(err)
			continue
		}

		receivedRow.RentID = rentID.Int32
		receivedRow.FromDate = fromDate.String
		receivedRow.ToDate = toDate.String
		receivedRow.Location = location.String
		receivedRow.AvailableExtras = strings.Split(extras.String, ",")
		receivedRow.Discounts = strings.Split(discounts.String, ",")
		receivedRow.AvailableLocations = strings.Split(locations, ",")
		receivedRow.CarDetails = carDetails.String
		result = append(result, receivedRow)
	}
	return len(result) != 0, nil
}
