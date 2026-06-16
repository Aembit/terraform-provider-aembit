package provider

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testClientWorkloadsDataSource string = "data.aembit_client_workloads.test"
	testClientWorkloadResource    string = "aembit_client_workload.test"
)

func testFindClientWorkload(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var err error
		var ok bool
		var notFound bool
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceSetID := rs.Primary.Attributes["resource_set_id"]

		if _, err, notFound = testClient.GetClientWorkload(context.Background(), rs.Primary.ID, nil, &resourceSetID); notFound {
			return err
		}
		return nil
	}
}

func TestAccClientWorkloadsDataSource(t *testing.T) {
	createFile1, _ := os.ReadFile("../../tests/client/data/TestAccClientWorkloadsDataSource_ResourceSet.tf")
	createFile2, _ := os.ReadFile("../../tests/client/data/TestAccClientWorkloadsDataSource_ProviderResourceSet.tf")
	createFile3, _ := os.ReadFile("../../tests/client/data/TestAccClientWorkloadsDataSource.tf")

	files := [3]string{string(createFile1), string(createFile2), string(createFile3)}

	for _, createFile := range files {
		randID := rand.Intn(10000000)
		createFileConfig := strings.ReplaceAll(
			createFile,
			"unittest1namespace",
			fmt.Sprintf("unittest1namespace%d", randID),
		)
		createFileConfig, _, _ = randomizeFileConfigs(
			createFileConfig,
			"",
			"Acceptance Test client workload",
		)
		createFileConfig, _, _ = randomizeFileConfigs(createFileConfig, "", "unittestname")

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: createFileConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						// Verify non-zero number of Client Workloads returned
						resource.TestCheckResourceAttrSet(
							testClientWorkloadsDataSource,
							"client_workloads.#",
						),
						// Verify dynamic values have any value set in the state.
						resource.TestCheckResourceAttrSet(
							testClientWorkloadsDataSource,
							"client_workloads.0.id",
						),
						// Find newly created entry
						testFindClientWorkload(testClientWorkloadResource),
					),
				},
			},
		})
	}

}
