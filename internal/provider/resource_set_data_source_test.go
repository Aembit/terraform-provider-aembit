package provider

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	testResourceSetDefault string = "data.aembit_resource_set.default"
	testResourceSetsAll    string = "data.aembit_resource_sets.all"
)

func TestAccDefaultResourceSet(t *testing.T) {
	t.Parallel()
	readFile, _ := os.ReadFile("../../tests/resource_set/TestAccDefaultResourceSet.tf")

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

					resource.TestMatchResourceAttr(
						testResourceSetsAll,
						"resource_sets.#",
						regexp.MustCompile(`[2-9]`),
					),
				),
			},
		},
	})
}
