package provider

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testResourceSetDefault    string = "data.aembit_resource_set.default"
	testResourceSetsAll       string = "data.aembit_resource_sets.all"
	testResourceSetDataSource string = "data.aembit_resource_sets.aembit_resource_set_datasource"
	testResourceSetResource   string = "aembit_resource_set.crs"
)

func testFindResourceSet(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var err error
		var ok bool
		var notFound bool
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if _, err, notFound = testClient.GetResourceSet(rs.Primary.ID, nil); notFound {
			return err
		}
		return nil
	}
}

func TestAccDefaultResourceSetDataSource(t *testing.T) {
	readFile, _ := os.ReadFile("../../tests/resource_set/data/TestAccDefaultResourceSet.tf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: string(readFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceSetDefault, "name", "Default"),
					resource.TestCheckResourceAttr(
						testResourceSetDefault,
						"id",
						"ffffffff-ffff-ffff-ffff-ffffffffffff",
					),
					resource.TestCheckResourceAttr(
						testResourceSetDefault,
						"effective_roles.#",
						"2",
					),

					resource.TestMatchResourceAttr(
						testResourceSetsAll,
						"resource_sets.#",
						regexp.MustCompile(`[2-9]`),
					),
					resource.TestCheckResourceAttr(
						testResourceSetsAll,
						"resource_sets.0.effective_roles.#",
						"2",
					),
				),
			},
		},
	})
}

func TestAccResourceSetDataSource(t *testing.T) {
	readFile, _ := os.ReadFile("../../tests/resource_set/data/TestAccResourceSetsDataSource.tf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: string(readFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						testResourceSetDataSource,
						"resource_sets.#",
						regexp.MustCompile(`[3-9]`),
					),
					resource.TestCheckResourceAttr(
						testResourceSetDataSource,
						"resource_sets.0.effective_roles.#",
						"2",
					),
					// Find newly created entry
					testFindResourceSet(testResourceSetResource),
				),
			},
		},
	})
}
