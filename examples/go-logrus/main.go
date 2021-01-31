package main

import (
	"os"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	log.Out = os.Stdout
	log.SetLevel(logrus.TraceLevel)
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	log.WithFields(logrus.Fields{
		"animal": "walrus",
	}).Trace("A walrus appears")

	log.WithFields(logrus.Fields{
		"animal": "walrus",
	}).Debug("A walrus appears")

	log.WithFields(logrus.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")

	log.WithFields(logrus.Fields{
		"animal": "walrus",
	}).Warn("A walrus appears")

	log.WithFields(logrus.Fields{
		"animal": "walrus",
	}).Error("A walrus appears")

	log.WithFields(logrus.Fields{
		"animal": "walrus",
	}).Fatal("A walrus appears")

	log.WithFields(logrus.Fields{
		"animal": "walrus",
	}).Panic("A walrus appears")
}
