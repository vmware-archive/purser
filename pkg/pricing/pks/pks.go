package pks

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

//Pricing srtucutre
type Pricing struct {
	Name               string  `json:"Name"`
	Region             string  `json:"Region"`
	VCPU               float64 `json:"VCPU"`
	GBMemory           float64 `json:"GBMemory"`
	GBStorage          float64 `json:"GBStorage"`
	LoadBalancer       float64 `json:"LoadBalancer"`
	DevelopmentCluster float64 `json:"DevelopmentCluster"`
	ProductionCluster  float64 `json:"ProductionCluster "`
	PKStoAWS           float64 `json:"PKStoAWS"`
	PKStoInternet      float64 `json:"PKStoInternet "`
	DataIO             float64 `json:"DataIO"`
}

func getPricingAttributes() []*Pricing {
	var pricing []*Pricing

	gp := os.Getenv("GOPATH")
	ap := filepath.Join(gp, "src/github.com/vmware/purser")
	fp := filepath.Join(ap, "pkg/utils/Cost.json")
	raw, err := ioutil.ReadFile(filepath.Clean(fp))
	if err != nil {
		return nil
	}
	err = json.Unmarshal(raw, &pricing)
	if err != nil {
		return nil
	}

	return pricing
}
