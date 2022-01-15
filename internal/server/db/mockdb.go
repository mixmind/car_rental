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

type DBStruct struct {
	internalDB *sql.DB
	carsArray  []domain.Car
}

/*
Generates and fills Inmemory db with cars
*/
func NewDBStruct() (*DBStruct, error) {
	inMemoryDB, err := sql.Open("sqlite3", "file:rental.db?cache=shared&mode=memory&_fk=true")
	if err != nil {
		return nil, err
	}
	db := DBStruct{internalDB: inMemoryDB}
	log.Info("Prefilling DB")
	if err := db.createInMemoryTables(); err != nil {
		return nil, errors.Wrap(err, "Failed to create tables")
	}
	db.generateCarsData()
	if err := db.insertCarsIntoDB(); err != nil {
		return nil, errors.Wrap(err, "Failed to insert cars into DB")
	}
	log.Info("DB filled sussesfully")
	return &db, nil
}

func NewDBStructWithDBProvided(inMemoryDB *sql.DB) *DBStruct {
	return &DBStruct{internalDB: inMemoryDB}
}

/*
Creates in memory tables such as cars and rents
*/
func (db *DBStruct) createInMemoryTables() error {
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

/*
Creates mock car data
*/
func (db *DBStruct) generateCarsData() {
	log.Info("Generating Cars data")

	numberOfCars := rand.Intn(50)
	for i := 0; i < numberOfCars; i++ {
		mockCar := cars.GenerateNewCar(time.Now().UTC().UnixNano())
		db.carsArray = append(db.carsArray, mockCar)
	}

	log.Info("Cars data generated sussesfully")
}

/*
Inserts mock car data into DB
*/
func (db *DBStruct) insertCarsIntoDB() error {
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

func (db *DBStruct) BeginTransaction() (*sql.Tx, error) {
	return db.internalDB.Begin()
}

func (db *DBStruct) Query(query string) (*sql.Rows, error) {
	return db.internalDB.Query(query)
}

func (db *DBStruct) Prepare(query string) (*sql.Stmt, error) {
	return db.internalDB.Prepare(query)
}
