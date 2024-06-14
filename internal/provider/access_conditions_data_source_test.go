package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccessConditionsDataSource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/condition/data/TestAccAccessConditionsDataSource.tf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of Access Conditions returned
					resource.TestCheckResourceAttrSet("data.aembit_access_conditions.test", "access_conditions.#"),
					// Verify AccessCondition Name
					resource.TestCheckResourceAttr("data.aembit_access_conditions.test", "access_conditions.0.name", "TF Acceptance Crowdstrike"),
					// Verify Tags.
					resource.TestCheckResourceAttr("data.aembit_access_conditions.test", "access_conditions.0.tags.%", "2"),
					resource.TestCheckResourceAttr("data.aembit_access_conditions.test", "access_conditions.0.tags.color", "blue"),
					resource.TestCheckResourceAttr("data.aembit_access_conditions.test", "access_conditions.0.tags.day", "Sunday"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("data.aembit_access_conditions.test", "access_conditions.0.id"),
				),
			},
		},
	})
}
