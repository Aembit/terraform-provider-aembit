package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const routingTestResource string = "aembit_routing.routing"

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
					// Verify Routing Name
					resource.TestCheckResourceAttr(routingTestResource, "name", newName),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(routingTestResource, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(routingTestResource, "id"),
				),
			},
			// Recreate the resource from the first test step
			{Config: createFileConfig},
			// ImportState testing
			{ResourceName: routingTestResource, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(routingTestResource, "name", newName),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
