package main

import (
	"car-rental/internal/server"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := server.Launch(); err != nil {
		log.Error(errors.Wrap(err, "Failed to launch server"))
		os.Exit(1)
	}
}
