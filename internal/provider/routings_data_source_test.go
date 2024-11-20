package provider

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testRoutingsDataSource string = "data.aembit_routings.test"
const testRoutingResource string = "aembit_routing.routing"

func testFindRouting(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var err error
		var ok bool
		var notFound bool
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if _, err, notFound = testClient.GetRouting(rs.Primary.ID, nil); notFound {
			return err
		}
		return nil
	}
}

func TestAccRoutingsDataSource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/routing/data/TestAccRoutingDataSource.tf")
	createFileConfig := string(createFile)

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

		createFileConfig = strings.ReplaceAll(string(createFile), "ffffffff-ffff-ffff-ffff-ffffffffffff", resourceSetId)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify non-zero number of Server Workloads returned
					resource.TestCheckResourceAttrSet(testRoutingsDataSource, "routings.#"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testRoutingsDataSource, "routings.0.id"),
					resource.TestCheckResourceAttrSet(testRoutingsDataSource, "routings.0.proxy_url"),
					resource.TestCheckResourceAttrSet(testRoutingsDataSource, "routings.0.resource_set_name"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testRoutingsDataSource, "routings.0.id"),
					// Find newly created entry
					testFindRouting(testRoutingResource),
				),
			},
		},
	})
}
