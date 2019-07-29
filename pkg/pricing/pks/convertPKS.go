package pks

import (
	"github.com/vmware/purser/pkg/controller/dgraph"
	"github.com/vmware/purser/pkg/controller/dgraph/models"
)

//constant values
const (
	linux = "linux"
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

func getPKSRateCard(pararegion string) *Pricing {
	pksPricing := getPricingAttributes()
	for _, specificRegion := range pksPricing {
		if specificRegion.Region == pararegion {
			return specificRegion
		}
	}
	return nil
}

// GetRateCardForPKS takes region as input and returns RateCard and error if any.
func GetRateCardForPKS(region string) *models.RateCard {
	specificRegionPreicing := getPKSRateCard(region)
	return getPurserRateCard(region, specificRegionPreicing)
}

//getPurserRateCard take region and pricingArray as input and returns RateCard for PKS of that region.
func getPurserRateCard(region string, pricing *Pricing) *models.RateCard {
	nodePrices, storagePrices := getResourceRateCard(pricing)

	return &models.RateCard{
		ID:            dgraph.ID{Xid: models.RateCardXID},
		IsRateCard:    true,
		CloudProvider: models.PKS,
		Region:        region,
		NodePrices:    nodePrices,
		StoragePrices: storagePrices,
	}
}

//take pksPricing as input and reutrn the nodePrice and storagePrice Array
func getResourceRateCard(pksPrice *Pricing) ([]*models.NodePrice, []*models.StoragePrice) {
	var nodePrices []*models.NodePrice
	var storagePrices []*models.StoragePrice

	nodePrices = updateComputePrices(pksPrice, nodePrices)
	storagePrices = updateStoragePrices(pksPrice, storagePrices)
	return nodePrices, storagePrices
}

func updateComputePrices(pksPrice *Pricing, nodePrices []*models.NodePrice) []*models.NodePrice {
	productXID := pksPrice.Region
	nodePrice := &models.NodePrice{
		ID:              dgraph.ID{Xid: productXID},
		IsNodePrice:     true,
		InstanceType:    pksPrice.Name,
		OperatingSystem: linux,
		PricePerCPU:     pksPrice.VCPU,
		PricePerMemory:  pksPrice.GBMemory,
	}
	uid := models.StoreNodePrice(nodePrice, productXID)
	if uid != "" {
		nodePrice.ID = dgraph.ID{UID: uid, Xid: productXID}
		nodePrices = append(nodePrices, nodePrice)
	}
	return nodePrices

}

func updateStoragePrices(pksPrice *Pricing, storagePrices []*models.StoragePrice) []*models.StoragePrice {
	productXID := pksPrice.Name
	storagePrice := &models.StoragePrice{
		ID:             dgraph.ID{Xid: productXID},
		IsStoragePrice: true,
		Price:          pksPrice.GBStorage,
	}
	uid := models.StoreStoragePrice(storagePrice, productXID)
	if uid != "" {
		storagePrice.ID = dgraph.ID{UID: uid, Xid: productXID}
		storagePrices = append(storagePrices, storagePrice)
	}
	return storagePrices
}
