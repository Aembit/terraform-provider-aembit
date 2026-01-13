package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testServerWorkloadsDataSource string = "data.aembit_server_workloads.test"

func testFindServerWorkload(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var err error
		var ok bool
		var notFound bool
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if _, err, notFound = testClient.GetServerWorkload(rs.Primary.ID, nil); notFound {
			return err
		}
		return nil
	}
}

func TestAccServerWorkloadsDataSource(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile("../../tests/server/data/TestAccServerWorkloadsDataSource.tf")
	createFileConfig, _, _ := randomizeFileConfigs(string(createFile), "", "Unit Test 1")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify non-zero number of Server Workloads returned
					resource.TestCheckResourceAttrSet(
						testServerWorkloadsDataSource,
						"server_workloads.#",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(
						testServerWorkloadsDataSource,
						"server_workloads.0.id",
					),
					resource.TestCheckResourceAttrSet(
						testServerWorkloadsDataSource,
						"server_workloads.0.service_endpoint.external_id",
					),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(
						testServerWorkloadsDataSource,
						"server_workloads.0.id",
					),
					// Find newly created entry
					testFindServerWorkload(testServerWorkloadResource),
				),
			},
		},
	})
}
