package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const resourceSetID string = "data.aembit_resource_set.default"

func TestAccDefaultResourceSet(t *testing.T) {
	readFile, _ := os.ReadFile("../../tests/resource_set/TestAccDefaultResourceSet.tf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: string(readFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceSetID, "name", "Default"),
					resource.TestCheckResourceAttr(resourceSetID, "id", "ffffffff-ffff-ffff-ffff-ffffffffffff"),
				),
			},
		},
	})
}
