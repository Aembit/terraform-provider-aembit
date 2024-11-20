package provider

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const routingTestResource string = "aembit_routing.routing"

func TestAccRoutingResource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/routing/TestAccRoutingResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/routing/TestAccRoutingResource.tfmod")
	createFileConfig, modifyFileConfig, newName := randomizeFileConfigs(string(createFile), string(modifyFile), "TF Acceptance Routing")

	if os.Getenv("CI") != "" { // indicates it is running in CI
		terraformVersion := getTerraformVersion()
		resourceSetId := "ffffffff-ffff-ffff-ffff-ffffffffffff"

		if strings.Contains(terraformVersion, "v1.6") {
			resourceSetId = "9538a706-936f-4fb9-8710-cc8f3096fd9b"
		} else if strings.Contains(terraformVersion, "v1.8") {
			resourceSetId = "a2302939-1232-4d40-8d0a-5af08115aa06"
		} else if strings.Contains(terraformVersion, "v1.9") {
			resourceSetId = "e3c81619-f708-47d4-a72f-0b6a296c5833"
		}

		createFileConfig = strings.ReplaceAll(createFileConfig, "ffffffff-ffff-ffff-ffff-ffffffffffff", resourceSetId)
		modifyFileConfig = strings.ReplaceAll(modifyFileConfig, "ffffffff-ffff-ffff-ffff-ffffffffffff", resourceSetId)
	}

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
