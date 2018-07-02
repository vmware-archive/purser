package main

import (
	"fmt"
	"kuber-controller/controller"
	"kuber-controller/config"
)

func main() {
	fmt.Println("Hello World")

	conf := config.Config{Resource: config.Resource{Pod: true}}
	controller.Start(&conf)
}
