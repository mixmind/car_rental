package server

import (
	"bytes"
	"car-rental/internal/server/cmds"
	"car-rental/internal/server/db"
	"car-rental/internal/server/domain"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	testCar = domain.Car{
		CarCompanyName:     "My company",
		Doors:              66,
		BigLuggage:         52,
		SmallLuggage:       33,
		AdultPlaces:        1,
		AirConditioner:     false,
		MinimumAge:         100,
		Price:              3,
		AvailableLocations: []string{"New York"},
		CarGroup:           14,
		Description:        "Some car description. Nothing special",
	}
	updateCar = domain.Car{
		CarCompanyName:     "My company1",
		Doors:              3,
		BigLuggage:         552,
		SmallLuggage:       323,
		AdultPlaces:        15,
		AirConditioner:     true,
		MinimumAge:         140,
		Price:              43,
		AvailableLocations: []string{"New Jersey"},
		CarGroup:           145,
		Description:        "Some car description. Nothing special!",
	}
	testRent = domain.RentInfo{
		FromDate:        "2022-01-15T15:13:30Z",
		ToDate:          "2022-01-16T15:13:30Z",
		Location:        "New York",
		Discounts:       []string{"5%"},
		AvailableExtras: []string{"Free day"},
		CarDetails: fmt.Sprintf(`%s %s.Part of %d group. With %d doors, %d adult places, %d big luggage and %d small luggage places.%s. For drivers with minimal age %d`,
			testCar.CarCompanyName,
			testCar.Description,
			testCar.CarGroup,
			testCar.Doors,
			testCar.AdultPlaces,
			testCar.BigLuggage,
			testCar.SmallLuggage,
			"Without Air Conditioner",
			testCar.MinimumAge),
		AgeGroup: "130",
		CarGroup: 14,
	}
	testRentError = domain.RentInfo{
		FromDate:        "2022-01-15T15:13:30Z",
		ToDate:          "2022-01-16T15:13:30Z",
		Location:        "New York",
		Discounts:       []string{"5%"},
		AvailableExtras: []string{"Free day"},
		CarDetails: fmt.Sprintf(`%s %s.Part of %d group. With %d doors, %d adult places, %d big luggage and %d small luggage places.%s. For drivers with minimal age %d`,
			testCar.CarCompanyName,
			testCar.Description,
			testCar.CarGroup,
			testCar.Doors,
			testCar.AdultPlaces,
			testCar.BigLuggage,
			testCar.SmallLuggage,
			"Without Air Conditioner",
			testCar.MinimumAge),
		AgeGroup: "130",
		CarGroup: 14,
	}
	carID         = ""
	inMemoryDB    *sql.DB
	err           error
	carProcessor  *cmds.CarProcessor
	rentProcessor *cmds.RentProcessor
)

func init() {
	go Launch()

	time.Sleep(time.Second * 1)
	inMemoryDB, err = sql.Open("sqlite3", "file:rental.db?cache=shared&mode=memory")
	if err != nil {
		log.Fatal("failed to create to inmemoryDB")
	}
	carProcessor = cmds.NewCarProcessor(db.NewDBStructWithDBProvided(inMemoryDB))
	rentProcessor = cmds.NewRentProcessor(db.NewDBStructWithDBProvided(inMemoryDB))
}

func TestAPICars(test *testing.T) {
	for i := 0; i < 1000; i++ {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/api/cars", restPort))
		if err != nil {
			test.Error(errors.Wrap(err, "Faled to request cars"))
			test.FailNow()
		}
		if resp.StatusCode != http.StatusOK {
			test.Error(fmt.Errorf("Status is incorrect. Received %d, want %d", resp.StatusCode, http.StatusOK))
			test.FailNow()
		}
		//We Read the response body on the line below.
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			test.Error(errors.Wrap(err, "Faled to read response body of cars"))
			test.FailNow()
		}
		var responseMessage domain.RestResponse
		var restCarsResult []domain.Car
		err = json.Unmarshal(body, &responseMessage)
		if err != nil {
			test.Error(errors.Wrap(err, "Faled to unpack response"))
			test.FailNow()
		}
		bytes, err := json.Marshal(responseMessage.ResponseMessage)
		if err != nil {
			test.Error(errors.Wrap(err, "Faled to marshal cars from response"))
			test.FailNow()
		}
		err = json.Unmarshal(bytes, &restCarsResult)
		if err != nil {
			test.Error(errors.Wrap(err, "Faled to unmarshal cars from response"))
			test.FailNow()
		}
		carsFromDB, err := carProcessor.GetCarsFromDB()
		if err != nil {
			test.Error(errors.Wrap(err, "Faled to extract cars from DB"))
			test.FailNow()
		}
		if !reflect.DeepEqual(carsFromDB, restCarsResult) || len(carsFromDB) != len(restCarsResult) {
			test.Error("Rest response is different from database data.")
			test.FailNow()
		}
	}
}

func TestAPIAddCarAndGetCar(test *testing.T) {
	jsonStr, err := json.Marshal(testCar)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to marshal car"))
		test.FailNow()
	}
	resp, err := http.Post(fmt.Sprintf("http://localhost:%d/api/cars", restPort), "application/json; charset=utf-8", bytes.NewBuffer(jsonStr))
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to create cars"))
		test.FailNow()
	}
	if resp.StatusCode != http.StatusCreated {
		test.Error(fmt.Errorf("Status is incorrect. Received %d, want %d", resp.StatusCode, http.StatusCreated))
		test.FailNow()
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to read response body of cars"))
		test.FailNow()
	}
	var responseMessage domain.RestResponse
	err = json.Unmarshal(body, &responseMessage)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to unpack response"))
		test.FailNow()
	}
	respMessage := ""
	switch v := responseMessage.ResponseMessage.(type) {
	case string:
		respMessage = v
	default:
		log.Infof("%t", v)
	}
	re := regexp.MustCompile("[0-9]+")
	carID = re.FindString(respMessage)
	resp, err = http.Get(fmt.Sprintf("http://localhost:%d/api/cars/%s", restPort, carID))
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to request car"))
		test.FailNow()
	}
	if resp.StatusCode != http.StatusOK {
		test.Error(fmt.Errorf("Status is incorrect. Received %d, want %d", resp.StatusCode, http.StatusOK))
		test.FailNow()
	}
	//We Read the response body on the line below.
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to read response body of cars"))
		test.FailNow()
	}
	err = json.Unmarshal(body, &responseMessage)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to unpack response"))
		test.FailNow()
	}
	var car domain.Car
	bytes, err := json.Marshal(responseMessage.ResponseMessage)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to marshal cars from response"))
		test.FailNow()
	}
	err = json.Unmarshal(bytes, &car)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to unmarshal cars from response"))
		test.FailNow()
	}
	carFromDB, err := carProcessor.GetCarFromDB(car.CarID)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to extract cars from DB"))
		test.FailNow()
	}
	testCar.CarID = car.CarID
	if !reflect.DeepEqual(carFromDB, &testCar) {
		test.Errorf("Rest response is different from database data. Post car\n[%v]\n Created car\n[%v]", &testCar, carFromDB)
		test.FailNow()
	} else {
		test.Log("Car created sussesfully")
	}
}

func TestAPIPutCar(test *testing.T) {
	var responseMessage domain.RestResponse
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/api/cars/1", restPort))
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to request car"))
		test.FailNow()
	}
	if resp.StatusCode != http.StatusOK {
		test.Error(fmt.Errorf("Status is incorrect. Received %d, want %d", resp.StatusCode, http.StatusOK))
		test.FailNow()
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to read response body of cars"))
		test.FailNow()
	}
	err = json.Unmarshal(body, &responseMessage)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to unpack response"))
		test.FailNow()
	}

	jsonStr, err := json.Marshal(updateCar)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to marshal car"))
		test.FailNow()
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:%d/api/cars/1", restPort), bytes.NewBuffer(jsonStr))
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to update car"))
		test.FailNow()
	}
	client := &http.Client{}
	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err = client.Do(req)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to update car"))
		test.FailNow()
	}

	if resp.StatusCode != http.StatusOK {
		test.Error(fmt.Errorf("Status is incorrect. Received %d, want %d", resp.StatusCode, http.StatusOK))
		test.FailNow()
	}
	carFromDB, err := carProcessor.GetCarFromDB(1)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to extract cars from DB"))
		test.FailNow()
	}
	updateCar.CarID = 1
	if !reflect.DeepEqual(carFromDB, &updateCar) {
		test.Errorf("Rest response is different from database data.Put car\n[%+v]\n Updated car\n[%+v]", &updateCar, carFromDB)
		test.FailNow()
	} else {
		test.Log("Car created sussesfully")
	}
}

func TestAPIDeleteCar(test *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:%d/api/cars/1", restPort), nil)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to delete car"))
		test.FailNow()
	}
	client := &http.Client{}
	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to delete car"))
		test.FailNow()
	}

	if resp.StatusCode != http.StatusOK {
		test.Error(fmt.Errorf("Status is incorrect. Received %d, want %d", resp.StatusCode, http.StatusOK))
		test.FailNow()
	}

	test.Log("Car deleted sussesfully")
	_, err = carProcessor.GetCarFromDB(1)
	if err == nil {
		test.Error("Error should be produces")
		test.FailNow()
	} else {
		test.Log("Car not found in DB")

	}
}

func TestAPIAddRentAndGetRent(test *testing.T) {
	carIDNumber, err := strconv.Atoi(carID)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to extract carID"))
		test.FailNow()
	}
	testRent.CarID = carIDNumber
	jsonStr, err := json.Marshal(testRent)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to marshal rent"))
		test.FailNow()
	}
	resp, err := http.Post(fmt.Sprintf("http://localhost:%d/api/rents", restPort), "application/json; charset=utf-8", bytes.NewBuffer(jsonStr))
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to create rents"))
		test.FailNow()
	}
	if resp.StatusCode != http.StatusCreated {
		test.Error(fmt.Errorf("Status is incorrect. Received %d, want %d", resp.StatusCode, http.StatusCreated))
		test.FailNow()
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to read response body of rents"))
		test.FailNow()
	}
	var responseMessage domain.RestResponse
	err = json.Unmarshal(body, &responseMessage)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to unpack response"))
		test.FailNow()
	}
	respMessage := ""
	switch v := responseMessage.ResponseMessage.(type) {
	case string:
		respMessage = v
	default:
		log.Infof("%t", v)
	}
	re := regexp.MustCompile("[0-9]+")
	rentID := re.FindString(respMessage)
	resp, err = http.Get(fmt.Sprintf("http://localhost:%d/api/rents/%s", restPort, rentID))
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to request rent"))
		test.FailNow()
	}
	if resp.StatusCode != http.StatusOK {
		test.Error(fmt.Errorf("Status is incorrect. Received %d, want %d", resp.StatusCode, http.StatusOK))
		test.FailNow()
	}
	//We Read the response body on the line below.
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to read response body of rents"))
		test.FailNow()
	}
	err = json.Unmarshal(body, &responseMessage)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to unpack response"))
		test.FailNow()
	}
	var rent domain.RentInfo
	bytes, err := json.Marshal(responseMessage.ResponseMessage)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to marshal rents from response"))
		test.FailNow()
	}
	err = json.Unmarshal(bytes, &rent)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to unmarshal rents from response"))
		test.FailNow()
	}
	rentFromDB, err := rentProcessor.GetRentFromDB(rent.RentID)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to extract rents from DB"))
		test.FailNow()
	}
	testRent.RentID = rent.RentID
	testRent.AgeGroup = rent.AgeGroup
	testRent.CarGroup = rent.CarGroup

	if !reflect.DeepEqual(rentFromDB, &testRent) {
		test.Errorf("Rest response is different from database data. Post rent\n[%+v]\n Created rent\n[%+v]", &testRent, rentFromDB)
		test.FailNow()
	} else {
		test.Log("Rent created sussesfully")
	}
}

func TestAPIAddRentError(test *testing.T) {
	carIDNumber, err := strconv.Atoi(carID)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to extract carID"))
		test.FailNow()
	}
	testRentError.CarID = carIDNumber
	jsonStr, err := json.Marshal(testRentError)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to marshal rent"))
		test.FailNow()
	}
	resp, err := http.Post(fmt.Sprintf("http://localhost:%d/api/rents", restPort), "application/json; charset=utf-8", bytes.NewBuffer(jsonStr))
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to create rents"))
		test.FailNow()
	}
	if resp.StatusCode != http.StatusConflict {
		test.Error(fmt.Errorf("Status is incorrect. Received %d, want %d", resp.StatusCode, http.StatusConflict))
		test.FailNow()
	} else {
		test.Log("Failed to create rent on same dates")
	}
	testRentError.FromDate = "2022-01-15T15:14:30Z"
	testRentError.ToDate = "2022-01-16T15:12:30Z"
	jsonStr, err = json.Marshal(testRentError)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to marshal rent"))
		test.FailNow()
	}
	resp, err = http.Post(fmt.Sprintf("http://localhost:%d/api/rents", restPort), "application/json; charset=utf-8", bytes.NewBuffer(jsonStr))
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to create rents"))
		test.FailNow()
	}
	if resp.StatusCode != http.StatusConflict {
		test.Error(fmt.Errorf("Status is incorrect. Received %d, want %d", resp.StatusCode, http.StatusConflict))
		test.FailNow()
	} else {
		test.Log("Failed to create rent between already created dates")
	}
	testRentError.FromDate = "2022-01-15T15:12:30Z"
	testRentError.ToDate = "2022-01-16T15:12:30Z"
	jsonStr, err = json.Marshal(testRentError)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to marshal rent"))
		test.FailNow()
	}
	resp, err = http.Post(fmt.Sprintf("http://localhost:%d/api/rents", restPort), "application/json; charset=utf-8", bytes.NewBuffer(jsonStr))
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to create rents"))
		test.FailNow()
	}
	if resp.StatusCode != http.StatusConflict {
		test.Error(fmt.Errorf("Status is incorrect. Received %d, want %d", resp.StatusCode, http.StatusConflict))
		test.FailNow()
	} else {
		test.Log("Failed to create rent when to date is between already created dates")
	}
	testRentError.FromDate = "2022-01-15T15:14:30Z"
	testRentError.ToDate = "2022-01-16T15:14:30Z"
	jsonStr, err = json.Marshal(testRentError)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to marshal rent"))
		test.FailNow()
	}
	resp, err = http.Post(fmt.Sprintf("http://localhost:%d/api/rents", restPort), "application/json; charset=utf-8", bytes.NewBuffer(jsonStr))
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to create rents"))
		test.FailNow()
	}
	if resp.StatusCode != http.StatusConflict {
		test.Error(fmt.Errorf("Status is incorrect. Received %d, want %d", resp.StatusCode, http.StatusConflict))
		test.FailNow()
	} else {
		test.Log("Failed to create rent when from date is between already created dates")
	}
	testRentError.FromDate = "2022-01-15T15:10:30Z"
	testRentError.ToDate = "2022-01-16T15:15:30Z"
	jsonStr, err = json.Marshal(testRentError)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to marshal rent"))
		test.FailNow()
	}
	resp, err = http.Post(fmt.Sprintf("http://localhost:%d/api/rents", restPort), "application/json; charset=utf-8", bytes.NewBuffer(jsonStr))
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to create rents"))
		test.FailNow()
	}
	if resp.StatusCode != http.StatusConflict {
		test.Error(fmt.Errorf("Status is incorrect. Received %d, want %d", resp.StatusCode, http.StatusConflict))
		test.FailNow()
	} else {
		test.Log("Failed to create rent when created date is inside new provided dates")
	}
}

func TestAPIRents(test *testing.T) {
	for i := 0; i < 1000; i++ {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/api/rents", restPort))
		if err != nil {
			test.Error(errors.Wrap(err, "Faled to request rents"))
			test.FailNow()
		}
		if resp.StatusCode != http.StatusOK {
			test.Error(fmt.Errorf("Status is incorrect. Received %d, want %d", resp.StatusCode, http.StatusOK))
			test.FailNow()
		}
		//We Read the response body on the line below.
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			test.Error(errors.Wrap(err, "Faled to read response body of rents"))
			test.FailNow()
		}
		var responseMessage domain.RestResponse
		var restRentsResult []domain.RentInfo
		err = json.Unmarshal(body, &responseMessage)
		if err != nil {
			test.Error(errors.Wrap(err, "Faled to unpack response"))
			test.FailNow()
		}
		bytes, err := json.Marshal(responseMessage.ResponseMessage)
		if err != nil {
			test.Error(errors.Wrap(err, "Faled to marshal rents from response"))
			test.FailNow()
		}
		err = json.Unmarshal(bytes, &restRentsResult)
		if err != nil {
			test.Error(errors.Wrap(err, "Faled to unmarshal rents from response"))
			test.FailNow()
		}
		rentsFromDB, err := rentProcessor.GetRentsFromDB()
		if err != nil {
			test.Error(errors.Wrap(err, "Faled to extract rents from DB"))
			test.FailNow()
		}
		if !reflect.DeepEqual(rentsFromDB, restRentsResult) || len(rentsFromDB) != len(restRentsResult) {
			test.Error("Rest response is different from database data.")
			test.FailNow()
		}
	}
}

func TestAPIDeleteRent(test *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:%d/api/rents/1", restPort), nil)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to delete rent"))
		test.FailNow()
	}
	client := &http.Client{}
	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		test.Error(errors.Wrap(err, "Faled to delete rent"))
		test.FailNow()
	}

	if resp.StatusCode != http.StatusOK {
		test.Error(fmt.Errorf("Status is incorrect. Received %d, want %d", resp.StatusCode, http.StatusOK))
		test.FailNow()
	}

	test.Log("Rent deleted sussesfully")
	_, err = rentProcessor.GetRentFromDB(1)
	if err == nil {
		test.Error("Error should be produces")
		test.FailNow()
	} else {
		test.Log("Rent not found in DB")

	}
	ShutDown()
}
