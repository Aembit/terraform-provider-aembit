package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testDiscoveryIntegrationResourceWiz string = "aembit_discovery_integration.wiz"

func TestAccDiscoveryIntegrationResource_Wiz(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile(
		"../../tests/discovery_integration/wiz/TestAccDiscoveryIntegrationResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/discovery_integration/wiz/TestAccDiscoveryIntegrationResource.tfmod",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify AccessCondition Name
					resource.TestCheckResourceAttr(
						testDiscoveryIntegrationResourceWiz,
						"name",
						"TF Acceptance Wiz",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testDiscoveryIntegrationResourceWiz, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testDiscoveryIntegrationResourceWiz, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      testDiscoveryIntegrationResourceWiz,
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						testDiscoveryIntegrationResourceWiz,
						"name",
						"TF Acceptance Wiz - Modified",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
