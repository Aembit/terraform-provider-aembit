package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testDiscoveryIntegrationssDataSource string = "data.aembit_discovery_integrations.test"
const testDiscoveryIntegrationResource string = "aembit_discovery_integration.wiz"

func testFindDiscoveryIntegration(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var err error
		var ok bool
		var notFound bool
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if _, err, notFound = testClient.GetDiscoveryIntegration(rs.Primary.ID, nil); notFound {
			return err
		}
		return nil
	}
}

func TestAccDiscoveryIntegrationsDataSource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/discovery_integration/data/TestAccDiscoveryIntegrationsDataSource.tf")
	createFileConfig, _, _ := randomizeFileConfigs(string(createFile), "", "TF Acceptance Wiz")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify non-zero number of Discovery Integrations returned
					resource.TestCheckResourceAttrSet(testDiscoveryIntegrationssDataSource, "discovery_integrations.#"),
					// Find newly created entry
					testFindDiscoveryIntegration(testDiscoveryIntegrationResource),
				),
			},
		},
	})
}
