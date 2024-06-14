package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAgentControllersDataSource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/agent_controllers/data/TestAccAgentControllersDataSource.tf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of Agent Controllers returned
					resource.TestCheckResourceAttrSet("data.aembit_agent_controllers.test", "agent_controllers.#"),
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr("data.aembit_agent_controllers.test", "agent_controllers.0.name", "TF Acceptance Azure Trust Provider"),
					// Verify Tags.
					resource.TestCheckResourceAttr("data.aembit_agent_controllers.test", "agent_controllers.0.tags.%", "2"),
					resource.TestCheckResourceAttr("data.aembit_agent_controllers.test", "agent_controllers.0.tags.color", "blue"),
					resource.TestCheckResourceAttr("data.aembit_agent_controllers.test", "agent_controllers.0.tags.day", "Sunday"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("data.aembit_agent_controllers.test", "agent_controllers.0.id"),
				),
			},
		},
	})
}
