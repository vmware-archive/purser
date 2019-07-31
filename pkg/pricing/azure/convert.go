package azure

import (
	"github.com/vmware/purser/pkg/controller/dgraph"
	"github.com/vmware/purser/pkg/controller/dgraph/models"
)

//constant values
const (
	coreToCPUConvertor = 2.0
	linux              = "Linux"
	windows            = "Windows"
	deliminator        = "-"
)

// PricingData structure
type PricingData struct {
	Unit         string
	PricePerUnit map[string]string
}

// Product structure
type Product struct {
	ProductFamily string
	Attributes    ProductAttributes
}

// ProductAttributes structure
type ProductAttributes struct {
	InstanceType    string
	InstanceFamily  string
	OperatingSystem string
	PreInstalledSW  string
	VolumeType      string
	UsageType       string
	Vcpu            string
	Memory          string
}

// GetRateCardForAzure takes region as input and returns RateCard and error if any.
func GetRateCardForAzure(region string) *models.RateCard {
	azurePricingArray, err := getAzureRateCard(region)
	if err == nil {
		return getPurserRateCard(region, azurePricingArray)
	}

	return nil

}

//getPurserRateCard take region and pricingArray as input and returns RateCard for Azure of that region.
func getPurserRateCard(region string, pricingArray []*Pricing) *models.RateCard {
	nodePrices, storagePrices := getResourceRateCard(pricingArray)
	return &models.RateCard{
		ID:            dgraph.ID{Xid: models.RateCardXID},
		IsRateCard:    true,
		CloudProvider: models.AZURE,
		Region:        region,
		NodePrices:    nodePrices,
		StoragePrices: storagePrices,
	}
}

//take azurePricing as input and reutrn the nodePrice and storagePrice Array
func getResourceRateCard(azurePricing []*Pricing) ([]*models.NodePrice, []*models.StoragePrice) {
	var nodePrices []*models.NodePrice
	var storagePrices []*models.StoragePrice

	for _, azurePrice := range azurePricing {

		nodePrices = updateComputePrices(azurePrice, nodePrices, "linux")
		nodePrices = updateComputePrices(azurePrice, nodePrices, "windows")
		storagePrices = updateStoragePrices(azurePrice, storagePrices)
	}
	return nodePrices, storagePrices
}

func updateComputePrices(azurePrice *Pricing, nodePrices []*models.NodePrice, operatingSystem string) []*models.NodePrice {
	productXID := azurePrice.CanonicalName + deliminator + operatingSystem
	pricePerCPU, pricePerGB := getPriceForUnitResources(azurePrice, operatingSystem)
	nodePrice := &models.NodePrice{
		ID:              dgraph.ID{Xid: productXID},
		IsNodePrice:     true,
		InstanceType:    azurePrice.Name,
		InstanceFamily:  azurePrice.CanonicalName,
		OperatingSystem: operatingSystem,
		Price:           azurePrice.LinuxPrice,
		PricePerCPU:     pricePerCPU,
		PricePerMemory:  pricePerGB,
		CPU:             azurePrice.NumberOfCores,
		Memory:          azurePrice.MemoryInMB / 1024,
	}
	uid := models.StoreNodePrice(nodePrice, productXID)
	if uid != "" {
		nodePrice.ID = dgraph.ID{UID: uid, Xid: productXID}
		nodePrices = append(nodePrices, nodePrice)
	}
	return nodePrices

}

func updateStoragePrices(azurePrice *Pricing, storagePrices []*models.StoragePrice) []*models.StoragePrice {
	productXID := azurePrice.CanonicalName
	storagePrice := &models.StoragePrice{
		ID:             dgraph.ID{Xid: productXID},
		IsStoragePrice: true,
		//Todo joshipr:Update pricing for storage
		Price: float64(azurePrice.ResourceDiskSizeInMB) * azurePrice.PricePerMemoryLinux,
	}
	uid := models.StoreStoragePrice(storagePrice, productXID)
	if uid != "" {
		storagePrice.ID = dgraph.ID{UID: uid, Xid: productXID}
		storagePrices = append(storagePrices, storagePrice)
	}
	return storagePrices
}

func getPriceForUnitResources(azurePrice *Pricing, operatingSystem string) (float64, float64) {
	pricePerCPU := models.DefaultCPUCostInFloat64
	pricePerMemory := models.DefaultMemCostInFloat64
	switch operatingSystem {
	case linux:
		{
			pricePerCPU = float64(azurePrice.NumberOfCores) * coreToCPUConvertor * azurePrice.PricePerCoreLinux
			//TODO joshipr: update price calculator for memory
			pricePerMemory = float64(azurePrice.MaxDataDiskCount) * azurePrice.PricePerMemoryLinux

		}
	case windows:
		{
			pricePerCPU = float64(azurePrice.NumberOfCores) * coreToCPUConvertor * azurePrice.PricePerCoreWindows
			//TODO joshipr: update price calculator for memory
			pricePerMemory = float64(azurePrice.MaxDataDiskCount) * azurePrice.PricePerMemoryWindows
		}
	}
	return pricePerCPU, pricePerMemory
}
