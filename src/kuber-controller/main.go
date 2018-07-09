package main

import (
	//"fmt"
	"kuber-controller/controller"
	"fmt"
	"log"
	//"kuber-controller/config"
	"kuber-controller/config"
)

func main() {
	fmt.Println("Hello world")
	log.Println("Hello world")
	conf := config.Config{Resource: config.Resource{Pod: true}}
	controller.Start(&conf)
	//controller.CreateCRDDefinition()
	//controller.TestCrdFlow()
}
