package rest

import (
	"car-rental/internal/server/cmds"
	"car-rental/internal/server/domain"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func rents(writer http.ResponseWriter, request *http.Request) {
	responseCode := http.StatusOK
	var responseMessage interface{}
	var err error

	switch request.Method {
	case http.MethodGet:
		responseMessage, err = cmds.GetRentsFromDB()
	case http.MethodPost:
		var rent domain.RentInfo
		err = parseBodyToObj(request, &rent)
		if err != nil {
			log.Error(err)
			responseCode = http.StatusBadRequest
			break
		}
		var car *domain.Car
		car, err = cmds.GetCarFromDB(rent.CarID)
		if err != nil {
			log.Error(err)
			responseCode = http.StatusBadRequest
			break
		}
		var id int64
		id, err = cmds.InsertRentInDB(rent, *car)
		if err != nil {
			log.Error(err)
			responseCode = http.StatusConflict
			responseMessage = "Failed to insert rent info"
		} else {
			responseCode = http.StatusCreated
			responseMessage = fmt.Sprintf("Rent info sussesfully inserted. Rent ID number = %d", id)
		}
	}

	if _, err := domain.WriteResponse(writer, responseCode, responseMessage, err); err != nil {
		log.Error(errors.Wrap(err, "Error occurred during writing response"))
	}
}

func rentDetails(writer http.ResponseWriter, request *http.Request) {
	responseCode := http.StatusOK
	var responseMessage interface{}
	var err error
	var rentID int
	skipProcessing := false
	rentID, err = extractPathID(request, domain.RentIDPathParam)
	if err != nil {
		log.Error(err)
		responseCode = http.StatusNotFound
		skipProcessing = true
	}
	if !skipProcessing {
		switch request.Method {
		case http.MethodGet:
			responseMessage, err = cmds.GetRentFromDB(rentID)

		case http.MethodDelete:
			_, err = cmds.RemoveRentFromDB(rentID)
			if err != nil {
				log.Error(err)
				responseCode = http.StatusInternalServerError
				responseMessage = "Failed to remove rent"
			} else {
				responseMessage = "Rent sussesfully removed"
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
