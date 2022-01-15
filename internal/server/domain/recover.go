package domain

import (
	"fmt"
	"net/http"
	"runtime/debug"

	log "github.com/sirupsen/logrus"
)

// RecoverHTTPPanic - recovers possible panic during HTTP processing and writes it into output HTTP message
func RecoverHTTPPanic(writer http.ResponseWriter) {
	if r := recover(); r != nil {
		err := fmt.Errorf("Detected panic: %s", r)
		log.Errorf("Error on processing the http request: %s\n%s", err, string(debug.Stack())) // print panic message and stack trace
		WriteResponse(writer, http.StatusBadRequest, "", err)
	}
}
