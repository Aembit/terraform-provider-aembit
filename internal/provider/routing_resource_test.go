package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const routingResourceID string = "aembit_routing.default"

func testDeleteRouting(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var ok bool
		var err error
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if ok, err = testClient.DeleteRouting(context.Background(), rs.Primary.ID, nil); !ok {
			return err
		}
		return nil
	}
}

func TestAccRoutingResource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/routing/TestAccRoutingResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/routing/TestAccRoutingResource.tfmod")
	createFileConfig, modifyFileConfig, newName := randomizeFileConfigs(string(createFile), string(modifyFile), "TF Acceptance Routing")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(routingResourceID, "name", newName),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(routingResourceID, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(routingResourceID, "id"),
					resource.TestCheckResourceAttrSet(routingResourceID, "resource_set_id"),
					resource.TestCheckResourceAttrSet(routingResourceID, "proxy_url"),
				),
			},
			// Test Aembit API Removal causes re-create with non-empty plan
			{Config: createFileConfig, Check: testDeleteRouting(routingResourceID), ExpectNonEmptyPlan: true},
			// Recreate the resource from the first test step
			{Config: createFileConfig},
			// ImportState testing
			{ResourceName: routingResourceID, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(routingResourceID, "name", newName),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
