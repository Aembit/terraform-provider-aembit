package provider

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testAgentControllerAzure      string = "aembit_agent_controller.azure_tp"
	testAgentControllerDeviceCode string = "aembit_agent_controller.device_code"
)

func testDeleteAgentController(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var ok bool
		var err error
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if ok, err = testClient.DeleteAgentController(context.Background(), rs.Primary.ID, nil); !ok {
			return err
		}
		return nil
	}
}

func TestAccAgentControllerResources(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/agent_controllers/TestAccAgentControllerResource.tf")
	modifyFile, _ := os.ReadFile(
		"../../tests/agent_controllers/TestAccAgentControllerResource.tfmod",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(
						testAgentControllerAzure,
						"name",
						"TF Acceptance Azure Trust Provider",
					),
					resource.TestCheckResourceAttr(
						testAgentControllerDeviceCode,
						"name",
						"TF Acceptance Device Code",
					),
					// Verify Tags.
					resource.TestCheckResourceAttr(testAgentControllerAzure, tagsCount, "2"),
					resource.TestCheckResourceAttr(testAgentControllerAzure, tagsColor, "blue"),
					resource.TestCheckResourceAttr(testAgentControllerAzure, tagsDay, "Sunday"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testAgentControllerAzure, "id"),
					resource.TestCheckResourceAttrSet(testAgentControllerDeviceCode, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testAgentControllerAzure, "id"),
				),
			},
			// Test Aembit API Removal causes re-create with non-empty plan
			{
				Config:             string(createFile),
				Check:              testDeleteAgentController(testAgentControllerAzure),
				ExpectNonEmptyPlan: true,
			},
			// Recreate the resource from the first test step
			{Config: string(createFile)},
			// ImportState testing
			{ResourceName: testAgentControllerAzure, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						testAgentControllerAzure,
						"name",
						"TF Acceptance Azure Trust Provider - Modified",
					),
					resource.TestCheckResourceAttr(
						testAgentControllerDeviceCode,
						"name",
						"TF Acceptance Device Code - Modified",
					),
					// Verify Tags.
					resource.TestCheckResourceAttr(testAgentControllerAzure, tagsCount, "2"),
					resource.TestCheckResourceAttr(testAgentControllerAzure, tagsColor, "orange"),
					resource.TestCheckResourceAttr(testAgentControllerAzure, tagsDay, "Tuesday"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccAgentControllerResource_Validation(t *testing.T) {
	emptyNameFile, _ := os.ReadFile(
		"../../tests/agent_controllers/TestAccAgentControllerResource.tfempty",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(emptyNameFile),
				ExpectError: regexp.MustCompile(
					`Attribute name string length must be between 1 and 128`,
				),
			},
		},
	})
}
