package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testIntegrationsDataSource string = "data.aembit_integrations.test"
	testIntegrationResource    string = "aembit_integration.wiz"
)

func testFindIntegration(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var err error
		var ok bool
		var notFound bool
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if _, err, notFound = testClient.GetIntegration(rs.Primary.ID, nil); notFound {
			return err
		}
		return nil
	}
}

func TestAccIntegrationsDataSource(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile("../../tests/integration/data/TestAccIntegrationsDataSource.tf")
	createFileConfig, _, _ := randomizeFileConfigs(string(createFile), "", "TF Acceptance Wiz")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify non-zero number of Integrations returned
					resource.TestCheckResourceAttrSet(testIntegrationsDataSource, "integrations.#"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(
						testIntegrationsDataSource,
						"integrations.0.id",
					),
					// Find newly created entry
					testFindIntegration(testIntegrationResource),
				),
			},
		},
	})
}
