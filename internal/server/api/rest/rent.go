package rest

import (
	"car-rental/internal/server/cmds"
	"car-rental/internal/server/domain"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

/*
Method responsible for rents listing and rent info creation
*/
func (restPr *RestProcessor) rents(writer http.ResponseWriter, request *http.Request) {
	carProcessor := cmds.NewCarProcessor(restPr.dbStruct)
	rentProcessor := cmds.NewRentProcessor(restPr.dbStruct)

	responseCode := http.StatusOK
	var responseMessage interface{}
	var err error

	switch request.Method {
	case http.MethodGet:
		responseMessage, err = rentProcessor.GetRentsFromDB()
	case http.MethodPost:
		var rent domain.RentInfo
		err = parseBodyToObj(request, &rent)
		if err != nil {
			log.Error(err)
			responseCode = http.StatusBadRequest
			break
		}
		insertRentProcessing := func() {
			restPr.carMutex.Lock()
			defer restPr.carMutex.Unlock()
			var car *domain.Car
			car, err = carProcessor.GetCarFromDB(rent.CarID)
			if err != nil {
				log.Error(err)
				responseCode = http.StatusBadRequest
				return
			}
			var id int64
			id, err = rentProcessor.InsertRentInDB(rent, *car)
			if err != nil {
				log.Error(err)
				responseCode = http.StatusConflict
				responseMessage = "Failed to insert rent info"
			} else {
				responseCode = http.StatusCreated
				responseMessage = fmt.Sprintf("Rent info sussesfully inserted. Rent ID number = %d", id)
			}
		}
		insertRentProcessing()
	}

	if _, err := domain.WriteResponse(writer, responseCode, responseMessage, err); err != nil {
		log.Error(errors.Wrap(err, "Error occurred during writing response"))
	}
}

/*
Method responsible for rent listing and rent deletion
*/
func (restPr *RestProcessor) rentDetails(writer http.ResponseWriter, request *http.Request) {
	rentProcessor := cmds.NewRentProcessor(restPr.dbStruct)

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
			responseMessage, err = rentProcessor.GetRentFromDB(rentID)

		case http.MethodDelete:
			_, err = rentProcessor.RemoveRentFromDB(rentID)
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
