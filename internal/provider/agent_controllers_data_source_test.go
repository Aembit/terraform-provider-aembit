package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testAgentControllersDataSource string = "data.aembit_agent_controllers.test"
	testAgentControllerResource    string = "aembit_agent_controller.azure_tp"
)

func testFindAgentController(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var err error
		var ok bool
		var notFound bool
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if _, err, notFound = testClient.GetAgentController(rs.Primary.ID, nil); notFound {
			return err
		}
		return nil
	}
}

func TestAccAgentControllersDataSource(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile(
		"../../tests/agent_controllers/data/TestAccAgentControllersDataSource.tf",
	)
	createFileConfig, _, _ := randomizeFileConfigs(
		string(createFile),
		"",
		"TF Acceptance Azure Trust Provider",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify non-zero number of Agent Controllers returned
					resource.TestCheckResourceAttrSet(
						testAgentControllersDataSource,
						"agent_controllers.#",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(
						testAgentControllersDataSource,
						"agent_controllers.0.id",
					),
					// Find newly created entry
					testFindAgentController(testAgentControllerResource),
				),
			},
		},
	})
}
