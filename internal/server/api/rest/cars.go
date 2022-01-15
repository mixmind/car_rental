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

func cars(writer http.ResponseWriter, request *http.Request) {
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
		id, err = cmds.InsertCarInDB(car)
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
			responseMessage, err = cmds.GetCarsFromDB()
		} else {
			responseMessage, err = cmds.GetCarsFromDBWithParams(urlValues)
		}
	default:
		responseCode = http.StatusBadRequest
		responseMessage = "This method is not allowed"
	}

	if _, err := domain.WriteResponse(writer, responseCode, responseMessage, err); err != nil {
		log.Error(errors.Wrap(err, "Error occurred during writing response"))
	}

}

func crudCars(writer http.ResponseWriter, request *http.Request) {
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
			responseMessage, err = cmds.GetCarFromDB(carID)
		case http.MethodPut:
			var car domain.Car
			err = parseBodyToObj(request, &car)
			if err != nil {
				log.Error(err)
				responseCode = http.StatusBadRequest
				break
			}
			_, err = cmds.UpdateCarInDB(car, carID)
			if err != nil {
				log.Error(err)
				responseCode = http.StatusInternalServerError
				responseMessage = "Failed to update car"
			} else {
				responseMessage = "Car sussesfully updated"
			}
		case http.MethodDelete:
			_, err = cmds.RemoveCarFromDB(carID)
			if err != nil {
				log.Error(err)
				responseCode = http.StatusInternalServerError
				responseMessage = "Failed to remove car"
			} else {
				responseMessage = "Car sussesfully removed"
			}

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
