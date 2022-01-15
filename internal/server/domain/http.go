package domain

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func WriteResponse(writer http.ResponseWriter, responseCode int, message interface{}, err error) (int, error) {
	writer.WriteHeader(responseCode)
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
		log.Error(errMsg)
	}
	responseMessage := RestResponse{ResponseMessage: message, ResponseError: errMsg}
	marshaledResponse, err := json.Marshal(responseMessage)
	if err != nil {
		return -1, err
	}
	return writer.Write(marshaledResponse)
}
