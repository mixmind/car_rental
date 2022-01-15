package rest

import (
	"car-rental/internal/server/db"
	"car-rental/internal/server/domain"
	"fmt"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type RestProcessor struct {
	dbStruct *db.DBStruct
	Router   *mux.Router
	carMutex *sync.RWMutex
}

/*
Creates router and defines REST API's
*/
func NewServer(dbStruct *db.DBStruct) (*mux.Router, error) {
	log.Info("Launching REST API's")
	rtr := mux.NewRouter()
	restProcessor := RestProcessor{dbStruct: dbStruct, carMutex: &sync.RWMutex{}}
	rtr.Handle("/api/cars", domain.WrapREST(restProcessor.cars)).Methods(http.MethodGet, http.MethodPost)
	rtr.Handle(fmt.Sprintf("/api/cars/{%s}", domain.CarIDPathParam), domain.WrapREST(restProcessor.crudCars)).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)
	rtr.Handle("/api/rents", domain.WrapREST(restProcessor.rents)).Methods(http.MethodGet, http.MethodPost)
	rtr.Handle(fmt.Sprintf("/api/rents/{%s}", domain.RentIDPathParam), domain.WrapREST(restProcessor.rentDetails)).Methods(http.MethodGet, http.MethodDelete)
	restProcessor.Router = rtr
	return rtr, nil
}
