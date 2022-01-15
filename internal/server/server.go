package server

import (
	"car-rental/internal/server/api/rest"
	"car-rental/internal/server/db"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

const restPort = 1020

var (
	server *http.Server
	ctx    context.Context
)

func Launch() error {
	err := db.NewDBStruct()
	if err != nil {
		return err
	}
	rtr, err := rest.NewServer()
	if err != nil {
		return err
	}
	go func() {
		ctx = context.Background()
		server = &http.Server{Addr: fmt.Sprintf(":%d", restPort), Handler: rtr}
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Errorln("HTTP endpoint returned error: ", err)
			os.Exit(1)
		} else {
			log.Info("Server closed gracefully")
		}

	}()
	select {}
}

func ShutDown() {
	fmt.Println("Shutting down the HTTP server...")
	server.Shutdown(ctx)
}
