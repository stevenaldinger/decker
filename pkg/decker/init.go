package decker

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	environment := os.Getenv("ENVIRONMENT")

	if environment == "production" {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.TraceLevel)
	}
}
