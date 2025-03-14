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

const testClientWorkloadsDataSource string = "data.aembit_client_workloads.test"
const testClientWorkloadResource string = "aembit_client_workload.test"

func testFindClientWorkload(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var err error
		var ok bool
		var notFound bool
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if _, err, notFound = testClient.GetClientWorkload(context.Background(), rs.Primary.ID, nil); notFound {
			return err
		}
		return nil
	}
}

func TestAccClientWorkloadsDataSource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/client/data/TestAccClientWorkloadsDataSource.tf")

	randID := rand.Intn(10000000)
	createFileConfig := strings.ReplaceAll(string(createFile), "unittest1namespace", fmt.Sprintf("unittest1namespace%d", randID))
	createFileConfig, _, _ = randomizeFileConfigs(createFileConfig, "", "Acceptance Test client workload")
	createFileConfig, _, _ = randomizeFileConfigs(createFileConfig, "", "unittestname")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify non-zero number of Client Workloads returned
					resource.TestCheckResourceAttrSet(testClientWorkloadsDataSource, "client_workloads.#"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testClientWorkloadsDataSource, "client_workloads.0.id"),
					// Find newly created entry
					testFindClientWorkload(testClientWorkloadResource),
				),
			},
		},
	})
}
