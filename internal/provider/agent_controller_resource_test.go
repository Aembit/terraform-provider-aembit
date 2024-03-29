package provider

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAgentControllerResources(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/agent_controllers/TestAccAgentControllerResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/agent_controllers/TestAccAgentControllerResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr("aembit_agent_controller.azure_tp", "name", "TF Acceptance Azure Trust Provider"),
					resource.TestCheckResourceAttr("aembit_agent_controller.device_code", "name", "TF Acceptance Device Code"),
					// Verify Tags.
					resource.TestCheckResourceAttr("aembit_agent_controller.azure_tp", "tags.%", "2"),
					resource.TestCheckResourceAttr("aembit_agent_controller.azure_tp", "tags.color", "blue"),
					resource.TestCheckResourceAttr("aembit_agent_controller.azure_tp", "tags.day", "Sunday"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_agent_controller.azure_tp", "id"),
					resource.TestCheckResourceAttrSet("aembit_agent_controller.device_code", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_agent_controller.azure_tp", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_agent_controller.azure_tp",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_agent_controller.azure_tp", "name", "TF Acceptance Azure Trust Provider - Modified"),
					resource.TestCheckResourceAttr("aembit_agent_controller.device_code", "name", "TF Acceptance Device Code - Modified"),
					// Verify Tags.
					resource.TestCheckResourceAttr("aembit_agent_controller.azure_tp", "tags.%", "2"),
					resource.TestCheckResourceAttr("aembit_agent_controller.azure_tp", "tags.color", "orange"),
					resource.TestCheckResourceAttr("aembit_agent_controller.azure_tp", "tags.day", "Tuesday"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccAgentControllerResource_Validation(t *testing.T) {
	emptyNameFile, _ := os.ReadFile("../../tests/agent_controllers/TestAccAgentControllerResource.tfempty")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      string(emptyNameFile),
				ExpectError: regexp.MustCompile(`Attribute name string length must be at least 1`), // <-- should match any error at all
			},
		},
	})
}
