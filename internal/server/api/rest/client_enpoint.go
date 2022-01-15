package rest

import (
	"car-rental/internal/server/domain"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

func NewServer() (*mux.Router, error) {
	log.Info("Launching REST API's")
	rtr := mux.NewRouter()

	rtr.HandleFunc("/api/cars", cars).Methods(http.MethodGet, http.MethodPost)
	rtr.HandleFunc(fmt.Sprintf("/api/cars/{%s}", domain.CarIDPathParam), crudCars).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)
	rtr.HandleFunc("/api/rents", rents).Methods(http.MethodGet, http.MethodPost)
	rtr.HandleFunc(fmt.Sprintf("/api/rents/{%s}", domain.RentIDPathParam), rentDetails).Methods(http.MethodGet, http.MethodDelete)

	return rtr, nil
}
