package cmds

import (
	"car-rental/internal/server/db"
	"car-rental/internal/server/domain"
	"database/sql"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

var inMemoryDB *sql.DB

func init() {
	var err error
	inMemoryDB, err = sql.Open("sqlite3", "file:rental.db?cache=shared&mode=memory")
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed connection to inmemory DB"))
	}
}

func InsertCarInDB(car domain.Car) (int64, error) {
	tx, err := inMemoryDB.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "Failed to start a transaction")
	}
	statement, err := tx.Prepare(db.InsertIntoCarTable)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to prepare stmt")
	}
	res, err := statement.Exec(car.CarCompanyName,
		car.Doors,
		car.BigLuggage,
		car.SmallLuggage,
		car.AdultPlaces,
		car.AirConditioner,
		car.MinimumAge,
		strings.Join(car.AvailableLocations, ","),
		car.CarGroup,
		car.Description,
		car.Price)
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

func GetCarsFromDB() ([]domain.Car, error) {
	rows, err := inMemoryDB.Query(db.SelectCars)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute a sql query")
	}

	var result []domain.Car
	for rows.Next() {
		var receivedRow domain.Car
		var locations string
		err = rows.Scan(&receivedRow.CarID, &receivedRow.CarCompanyName,
			&receivedRow.Doors,
			&receivedRow.BigLuggage,
			&receivedRow.SmallLuggage,
			&receivedRow.AdultPlaces,
			&receivedRow.AirConditioner,
			&receivedRow.MinimumAge,
			&locations,
			&receivedRow.CarGroup,
			&receivedRow.Description,
			&receivedRow.Price)
		if err != nil {
			log.Error(err)
			continue
		}
		receivedRow.AvailableLocations = strings.Split(locations, ",")

		result = append(result, receivedRow)
	}
	return result, nil
}

func GetCarsFromDBWithParams(values map[string][]string) ([]domain.CombinedRentInfo, error) {
	searchParams := extractURLValues(values)
	query := db.SelectCarsRents

	if len(searchParams.DefaultFilter) > 0 {
		query += " WHERE " + searchParams.DefaultFilter
	}
	if len(searchParams.DefaultFilter) == 0 && len(searchParams.DateFilter) > 0 {
		query += " WHERE " + searchParams.DateFilter
	} else if len(searchParams.DefaultFilter) > 0 && len(searchParams.DateFilter) > 0 {
		query += " and " + searchParams.DateFilter
	}
	log.Info(query)
	rows, err := inMemoryDB.Query(query)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute a sql query")
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
	return result, nil
}

func GetCarFromDB(carID int) (*domain.Car, error) {
	stmt, err := inMemoryDB.Prepare(fmt.Sprintf("%s WHERE car_id=?", db.SelectCars))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to prepare an sql query")
	}
	rows := stmt.QueryRow(carID)
	var receivedRow domain.Car
	var locations string
	err = rows.Scan(&receivedRow.CarID, &receivedRow.CarCompanyName,
		&receivedRow.Doors,
		&receivedRow.BigLuggage,
		&receivedRow.SmallLuggage,
		&receivedRow.AdultPlaces,
		&receivedRow.AirConditioner,
		&receivedRow.MinimumAge,
		&locations,
		&receivedRow.CarGroup,
		&receivedRow.Description, &receivedRow.Price)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to query car")
	}
	receivedRow.AvailableLocations = strings.Split(locations, ",")
	return &receivedRow, nil
}

func UpdateCarInDB(car domain.Car, carID int) (int64, error) {
	stmt, err := inMemoryDB.Prepare(db.UpdateCar)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to prepare an sql query")
	}

	res, err := stmt.Exec(car.CarCompanyName,
		car.Doors,
		car.BigLuggage,
		car.SmallLuggage,
		car.AdultPlaces,
		car.AirConditioner,
		car.MinimumAge,
		strings.Join(car.AvailableLocations, ","),
		car.CarGroup,
		car.Description,
		car.Price,
		carID)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to execute car update")
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "Failed to extract rows updated number")
	}

	return affect, nil
}

func RemoveCarFromDB(carID int) (int64, error) {
	stmt, err := inMemoryDB.Prepare(db.RemoveCar)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to prepare an sql query")
	}

	res, err := stmt.Exec(carID)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to execute car delete")
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "Failed to extract rows deleted number")
	}

	return affect, nil
}

func extractURLValues(values map[string][]string) domain.SearchParams {
	fromExists := ""
	toExists := ""
	if from, ok := values[domain.FromDateUrlValue]; ok {
		if len(from) == 1 && len(from[0]) > 0 {
			fromExists = from[0]
		}
	}
	if to, ok := values[domain.ToDateUrlValue]; ok {
		if len(to) == 1 && len(to[0]) > 0 {
			if len(fromExists) > 0 {
				toExists = to[0]
			}
		}
	}
	searchParams := buildFromToFilter(fromExists, toExists, false)
	addedFilter := false
	if location, ok := values[domain.LocationUrlValue]; ok {
		if len(location) == 1 && len(location[0]) > 0 {
			splittedLocations := strings.Split(location[0], ",")
			addedFilter = true
			searchParams.DefaultFilter += "("
			for index, loc := range splittedLocations {
				searchParams.DefaultFilter += "locations like '%" + loc + "%'"
				if index < len(splittedLocations)-1 {
					searchParams.DefaultFilter += " or "
				}
			}
			searchParams.DefaultFilter += ")"
		}
	}
	if age, ok := values[domain.AgeGroupUrlValue]; ok {
		if len(age) == 1 && len(age[0]) > 0 {
			if addedFilter {
				searchParams.DefaultFilter += " and "
			}
			addedFilter = true
			splittedAge := strings.Split(age[0], "-")
			if len(splittedAge) == 1 {
				searchParams.DefaultFilter += "min_age>=" + splittedAge[0]
			} else if len(splittedAge) == 2 {
				searchParams.DefaultFilter += " min_age between " + splittedAge[0] + " and " + splittedAge[1]
			}
		}
	}
	if car, ok := values[domain.CarGroupUrlValue]; ok {
		if len(car) == 1 && len(car[0]) > 0 {
			if addedFilter {
				searchParams.DefaultFilter += " and "
			}
			searchParams.DefaultFilter += " car_group=" + car[0]
		}
	}
	return searchParams
}

func buildFromToFilter(from string, to string, isExists bool) domain.SearchParams {
	var searchParams domain.SearchParams
	fromExists := false
	if len(from) > 0 {
		fromExists = true
	}
	if fromExists && len(to) > 0 {
		brokenTime := false
		covertFrom, err := time.Parse(domain.TimeLayout, from)
		if err != nil {
			brokenTime = true
			log.Error(err)
		}
		covertTo, err := time.Parse(domain.TimeLayout, to)
		if err != nil {
			brokenTime = true
			log.Error(err)
		}
		if !covertFrom.Before(covertTo) {
			log.Error("Please provide correct dates, from must be less than too")
			brokenTime = true
		}
		if !brokenTime && !isExists {
			searchParams.DateFilter += " (to_time IS NULL or to_time<'" + from + "') or"
			searchParams.DateFilter += " (from_time IS NULL or from_time>'" + to + "')"
		} else if !brokenTime && isExists {
			searchParams.DateFilter += " (to_time IS NOT NULL AND from_time IS NOT NULL) and"
			searchParams.DateFilter += " ((from_time between '" + from + "' and '" + to + "') or (to_time between '" + from + "' and '" + to + "')) or "
			searchParams.DateFilter += " (from_time <='" + from + "' and to_time >= '" + to + "')"
		}

	}
	return searchParams
}
