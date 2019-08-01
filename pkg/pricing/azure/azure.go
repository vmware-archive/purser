package azure

import (
	json "encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	s "strings"
)

//Pricing stores pricing
type Pricing struct {
	Name                  string  `json:"name"`
	CanonicalName         string  `json:"canonicalname"`
	NumberOfCores         float64 `json:"numberOfCores"`
	OsDiskSizeInMB        int     `json:"osDiskSizeInMB"`
	ResourceDiskSizeInMB  int     `json:"resourceDiskSizeInMB"`
	MemoryInMB            float64 `json:"memoryInMB"`
	MaxDataDiskCount      int     `json:"maxDataDiskCount"`
	RegionName            string  `json:"regionName"`
	LinuxPrice            float64 `json:"linuxPrice"`
	WidowsPrice           float64 `json:"windowsPrice"`
	PricePerCoreLinux     float64 `json:"pricePerCoreLinux"`
	PricePerMemoryLinux   float64 `json:"pricePerMemoryLinux"`
	PricePerCoreWindows   float64 `json:"pricePerCoreWindows"`
	PricePerMemoryWindows float64 `json:"pricePerMemoryWindows"`
}

func getAzureRateCard(region string) ([]*Pricing, error) {
	// Make HTTP GET request
	resp, err := http.Get(getAzureURLForRegion(region))
	if err != nil {
		log.Fatal(err)
	}
	defer Check(resp.Body.Close)

	var Pricing []*Pricing
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		bodyString := string(bodyBytes)
		jsonStringStart := s.Split(bodyString, "json = ")[1]
		jsonStringEnd := s.Replace(s.Split(jsonStringStart, "</script>")[0], ";", "", 1)
		error := json.Unmarshal([]byte(jsonStringEnd), &Pricing)
		if error == nil {
			return Pricing, nil
		}
	}
	return nil, nil
}

//GetAzurePricingUrl function details
//input:region
//return azurePrice url for the region
func getAzureURLForRegion(region string) string {
	return "https://www.azureprice.net/?region=" + region + "&timeoption=hour"
}

//Check takes a fuction as input and checks if it returns an error
func Check(f func() error) {
	if err := f(); err != nil {
		fmt.Println("Received error:", err)
	}
}
