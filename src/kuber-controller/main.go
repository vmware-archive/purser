package main

import (
	"fmt"
	//"kuber-controller/config"
	"kuber-controller/controller"
)

func main() {
	fmt.Println("Hello World")

	//conf := config.Config{Resource: config.Resource{Node: true}}
	//controller.Start(&conf)
	//controller.CreateCRDDefinition()
	controller.TestCrdFlow()
}
