package pricing

import (
	"testing"

	"github.com/vmware/purser/test/utils"

	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph"
	"github.com/vmware/purser/pkg/controller/dgraph/models"
	"github.com/vmware/purser/pkg/pricing/aws"
)

// TestAWSPricingFlow it should populate your dgraph running at localhost 9080 port with aws compute and storage prices
// The following dgraph query will give the rate card data
// {
//		rateCard(func: has(isRateCard)) {
//			cloudProvider
//			region
//			nodePrices {
//				instanceType
//				operatingSystem
//				price
//				instanceFamily
//			}
//			storagePrices {
//				volumeType
//				usageType
//				price
//			}
//		}
// }
func TestAWSPricingFlow(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	dgraph.Start("localhost", "9080")
	rateCard := aws.GetRateCardForAWS("us-east-1")
	models.StoreRateCard(rateCard)
	defer dgraph.Close()
	utils.Assert(t, rateCard != nil, "rate card is nil")
}
