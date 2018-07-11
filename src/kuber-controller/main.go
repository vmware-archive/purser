package main

import (
	"kuber-controller/controller"
	"kuber-controller/config"
)

func main() {
	conf := config.Config{Resource: config.Resource{Pod: true}}
	controller.Start(&conf)
	//controller.CreateCRDDefinition()
	//controller.TestCrdFlow()
}
