package provider

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
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

		resourceSetID := rs.Primary.Attributes["resource_set_id"]

		if _, err, notFound = testClient.GetServerWorkload(rs.Primary.ID, nil, &resourceSetID); notFound {
			return err
		}
		return nil
	}
}

func TestAccServerWorkloadsDataSource(t *testing.T) {
	//createFile1, _ := os.ReadFile("../../tests/server/data/TestAccServerWorkloadsDataSource_ResourceSet.tf")
	createFile2, _ := os.ReadFile("../../tests/server/data/TestAccServerWorkloadsDataSource_ProviderResourceSet.tf")
	//createFile3, _ := os.ReadFile("../../tests/server/data/TestAccServerWorkloadsDataSource.tf")

	files := [1]string{string(createFile2)}

	for _, createFile := range files {
		randID := rand.Intn(10000000)
		createFileConfig := strings.ReplaceAll(
			createFile,
			"unittest.testhost.com",
			fmt.Sprintf("unittest%d.testhost.com", randID),
		)
		createFileConfig, _, _ = randomizeFileConfigs(createFileConfig, "", "Unit Test 1")

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
}
