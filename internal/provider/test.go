package provider

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var maxRand = big.NewInt(10000000)

const tagsCount = "tags.%"
const tagsColor = "tags.color"
const tagsDay = "tags.day"

func randomizeFileConfigs(newConfig, modifyConfig, startValue string) (string, string, string) {
	randID, _ := rand.Int(rand.Reader, maxRand)

	endValue := fmt.Sprintf("%s%d", startValue, randID)
	return strings.ReplaceAll(newConfig, startValue, endValue), strings.ReplaceAll(modifyConfig, startValue, endValue), endValue
}

func checkValidClientID(resourceName, attributeName, arnType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		var clientID string = rs.Primary.Attributes[attributeName]
		if len(clientID) <= 0 {
			return fmt.Errorf("empty client id: %s %s", resourceName, attributeName)
		}

		if len(strings.Split(clientID, ":")) != 6 {
			return fmt.Errorf("clientID does not have the correct number of blocks: %s", clientID)
		}
		if !strings.HasPrefix(clientID, "aembit:") {
			return fmt.Errorf("clientID does not have the correct prefix: %s", clientID)
		}
		if !strings.Contains(clientID, arnType) {
			return fmt.Errorf("clientID does not have the expected ARN type: %s", clientID)
		}
		if len(strings.Split(strings.Split(clientID, ":")[5], "-")) != 5 {
			return fmt.Errorf("clientID identifier is not a valid GUID: %s", clientID)
		}
		return nil
	}
}
