package pricing

import (
	"testing"

	"github.com/vmware/purser/test/utils"

	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph"
	"github.com/vmware/purser/pkg/controller/dgraph/models"
	"github.com/vmware/purser/pkg/pricing/gcp"
)

func TestGCPPricingFlow(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	dgraph.Start("localhost", "9080")
	rateCard := gcp.GetRateCardForGCP("us-east1")
	models.StoreRateCard(rateCard)
	defer dgraph.Close()
	utils.Assert(t, rateCard != nil, "rate card is nil")
	utils.Assert(t, len(rateCard.NodePrices) != 0, "no node prices found")
	utils.Assert(t, len(rateCard.StoragePrices) != 0, "no storage prices found")
}
