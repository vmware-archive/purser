package pricing

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph"
	"github.com/vmware/purser/pkg/controller/dgraph/models"
	"github.com/vmware/purser/pkg/pricing/aws"
)

// TestAWSPricingFlow it should populate your dgraph running at localhost 9080 port with aws compute and storage prices
func TestAWSPricingFlow(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	dgraph.Start("localhost", "9080")
	var rateCard *models.RateCard
	rateCard = aws.GetRateCardForAWS("us-east-1")
	models.StoreRateCard(rateCard)
	defer dgraph.Close()
}
