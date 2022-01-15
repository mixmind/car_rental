package rest

import (
	"car-rental/internal/server/cmds"
	"car-rental/internal/server/domain"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

/*
Method responsible for cars listing and new car creating
*/
func (restPr *RestProcessor) cars(writer http.ResponseWriter, request *http.Request) {
	carProcessor := cmds.NewCarProcessor(restPr.dbStruct)
	responseCode := http.StatusOK
	var responseMessage interface{}
	var err error
	switch request.Method {
	case http.MethodPost:
		var car domain.Car
		err = parseBodyToObj(request, &car)
		if err != nil {
			log.Error(err)
			responseCode = http.StatusBadRequest
			break
		}
		var id int64
		id, err = carProcessor.InsertCarInDB(car)
		if err != nil {
			log.Error(err)
			responseCode = http.StatusInternalServerError
			responseMessage = "Failed to insert car"
		} else {
			responseCode = http.StatusCreated
			responseMessage = fmt.Sprintf("New Car Sussesfully Added. Car ID number = %d", id)
		}
	case http.MethodGet:
		urlValues := request.URL.Query()
		if len(urlValues) == 0 {
			responseMessage, err = carProcessor.GetCarsFromDB()
		} else {
			responseMessage, err = carProcessor.GetCarsFromDBWithParams(urlValues)
		}
	default:
		responseCode = http.StatusBadRequest
		responseMessage = "This method is not allowed"
	}

	if _, err := domain.WriteResponse(writer, responseCode, responseMessage, err); err != nil {
		log.Error(errors.Wrap(err, "Error occurred during writing response"))
	}

}

/*
Method responsible for car listing, car update and car deletion
*/
func (restPr *RestProcessor) crudCars(writer http.ResponseWriter, request *http.Request) {
	carProcessor := cmds.NewCarProcessor(restPr.dbStruct)
	responseCode := http.StatusOK
	var responseMessage interface{}
	var err error
	var carID int
	skipProcessing := false
	carID, err = extractPathID(request, domain.CarIDPathParam)
	if err != nil {
		log.Error(err)
		responseCode = http.StatusNotFound
		skipProcessing = true
	}
	if !skipProcessing {
		switch request.Method {
		case http.MethodGet:
			responseMessage, err = carProcessor.GetCarFromDB(carID)
		case http.MethodPut:
			updateCarProcessing := func() {
				restPr.carMutex.Lock()
				defer restPr.carMutex.Unlock()
				var car domain.Car
				err = parseBodyToObj(request, &car)
				if err != nil {
					log.Error(err)
					responseCode = http.StatusBadRequest
					return
				}
				_, err = carProcessor.UpdateCarInDB(car, carID)
				if err != nil {
					log.Error(err)
					responseCode = http.StatusInternalServerError
					responseMessage = "Failed to update car"
				} else {
					responseMessage = "Car sussesfully updated"
				}
			}
			updateCarProcessing()
		case http.MethodDelete:
			removeCarProcessing := func() {
				restPr.carMutex.Lock()
				defer restPr.carMutex.Unlock()
				_, err = carProcessor.RemoveCarFromDB(carID)
				if err != nil {
					log.Error(err)
					responseCode = http.StatusInternalServerError
					responseMessage = "Failed to remove car"
				} else {
					responseMessage = "Car sussesfully removed"
				}
			}
			removeCarProcessing()
		default:
			responseCode = http.StatusBadRequest
			responseMessage = "This method is not allowed"
		}
	}

	if _, err := domain.WriteResponse(writer, responseCode, responseMessage, err); err != nil {
		log.Error(errors.Wrap(err, "Error occurred during writing response"))
	}
}

func parseBodyToObj(request *http.Request, obj interface{}) error {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return errors.Wrap(err, "Incorrect body in request")
	}
	log.Debugf("Request body [%s]", string(requestBody))
	err = json.Unmarshal(requestBody, obj)
	if err != nil {
		return errors.Wrap(err, "Incorrect format of body")
	}
	return nil
}

func extractPathID(request *http.Request, stringParam string) (int, error) {
	requestPathParams := mux.Vars(request)

	pathParamValueString, ok := requestPathParams[stringParam]
	if !ok {
		return 0, fmt.Errorf("Provided path parameter [%s] not found", stringParam)
	}

	intParamValue, err := strconv.ParseInt(pathParamValueString, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(intParamValue), nil
}
