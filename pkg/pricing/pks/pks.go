package pks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//Pricing srtucutre
type Pricing struct {
	Name               string
	Region             string
	VCPU               float64
	GBMemory           float64
	GBStorage          float64
	LoadBalancer       float64
	DevelopmentCluster float64
	ProductionCluster  float64
	PKStoAWS           float64
	PKStoInternet      float64
	DataIO             float64
}

func getPricingAttributes() []*Pricing {
	var pricing []*Pricing
	raw, err := ioutil.ReadFile("pkg/utils/Cost.json") //How to make particular file read !!
	if err != nil {
		fmt.Println(err.Error())
		// os.Exit(1)
		return nil
	}
	json.Unmarshal(raw, &pricing)
	fmt.Println("getPricingAttributes", pricing)
	return pricing
}
