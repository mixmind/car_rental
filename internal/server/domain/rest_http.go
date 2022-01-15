package domain

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// RequestHTTPFilterFunc - prefilters incoming requests
type RequestHTTPFilterFunc = func(request *http.Request) (int, error)

func WrapREST(handleFunc func(writer http.ResponseWriter, request *http.Request), filterFunc ...RequestHTTPFilterFunc) http.Handler {

	processingFunc := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer RecoverHTTPPanic(writer)

		log.Debugln("Start handle request", request)
		handleStartTime := time.Now()

		handleFunc(writer, request) // call original

		handleEndTime := time.Now()
		log.Debugln("Request handled in", handleEndTime.Sub(handleStartTime))
	})

	return processingFunc
}
