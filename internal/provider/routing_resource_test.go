package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const routingTestResource string = "aembit_routing.routing"

func TestAccRoutingResource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/routing/TestAccRoutingResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/routing/TestAccRoutingResource.tfmod")
	createFileConfig, modifyFileConfig, newName := randomizeFileConfigs(string(createFile), string(modifyFile), "TF Acceptance Routing")

	fmt.Println(os.Environ())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Routing Name testing
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
