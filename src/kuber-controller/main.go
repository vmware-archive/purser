package main

import (
	"kuber-controller/controller"
	//"kuber-controller/config"
	log "github.com/Sirupsen/logrus"
	"os"
)

func init() {
	setlogFile()
}

func setlogFile() {
	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	//conf := config.Config{Resource: config.Resource{Pod: true}}
	//controller.Start(&conf)
	//controller.CreateCRDDefinition()
	controller.TestCrdFlow()
}
