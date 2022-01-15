package db

import (
	"car-rental/internal/server/cars"
	"car-rental/internal/server/domain"
	"database/sql"
	"math/rand"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

var (
	inMemoryDB *sql.DB
	db         dbStruct
)

type dbStruct struct {
	internalDB *sql.DB
	carsArray  []domain.Car
}

func init() {
	var err error
	inMemoryDB, err = sql.Open("sqlite3", "file:rental.db?cache=shared&mode=memory")
	if err != nil {
		log.Fatal("Failed to create inmemoryDB")
	}
}

func NewDBStruct() error {
	db = dbStruct{internalDB: inMemoryDB}
	if err := db.createInMemoryTables(); err != nil {
		return errors.Wrap(err, "Failed to create tables")
	}
	db.generateCarsData()
	if err := db.insertCarsIntoDB(); err != nil {
		return errors.Wrap(err, "Failed to insert cars into DB")
	}
	return nil
}

func (db *dbStruct) createInMemoryTables() error {
	log.Info("Creating cars table")
	if _, err := db.internalDB.Exec(createCarTable); err != nil {
		return errors.Wrapf(err, "Failed to create cars table")
	}
	log.Info("Cars table created sussesfully")
	log.Info("Creating rents table")
	if _, err := db.internalDB.Exec(createRentTable); err != nil {
		return errors.Wrapf(err, "Failed to create rent table")
	}
	log.Info("Rents table created sussesfully")
	return nil
}

func (db *dbStruct) generateCarsData() {
	log.Info("Generating Cars data")

	numberOfCars := rand.Intn(50)
	for i := 0; i < numberOfCars; i++ {
		mockCar := cars.GenerateNewCar(time.Now().UTC().UnixNano())
		db.carsArray = append(db.carsArray, mockCar)
	}

	log.Info("Cars data generated sussesfully")
}

func (db *dbStruct) insertCarsIntoDB() error {
	log.Println("Inserting cars records")
	tx, err := db.internalDB.Begin()
	if err != nil {
		return errors.Wrap(err, "Failed to start a transaction")
	}
	statement, err := tx.Prepare(InsertIntoCarTable)
	if err != nil {
		return errors.Wrap(err, "Failed to prepare stmt")
	}
	for _, car := range db.carsArray {
		_, err := statement.Exec(car.CarCompanyName,
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
			return errors.Wrap(err, "Failed to execute a prepared statement")
		}
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "Failed to commit a transaction")
	}
	log.Info("Cars records inserted sussesfully")

	return nil
}
