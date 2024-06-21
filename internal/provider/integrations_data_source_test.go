package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testIntegrationsDataSource string = "data.aembit_integrations.test"

func TestAccIntegrationsDataSource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/integration/data/TestAccIntegrationsDataSource.tf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of Integrations returned
					resource.TestCheckResourceAttrSet(testIntegrationsDataSource, "integrations.#"),
					// Verify Integration Name
					resource.TestCheckResourceAttr(testIntegrationsDataSource, "integrations.0.name", "TF Acceptance Wiz"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testIntegrationsDataSource, "integrations.0.id"),
					// Verify Tags.
					resource.TestCheckResourceAttr(testIntegrationsDataSource, "integrations.0.tags.%", "2"),
					resource.TestCheckResourceAttr(testIntegrationsDataSource, "integrations.0.tags.color", "blue"),
					resource.TestCheckResourceAttr(testIntegrationsDataSource, "integrations.0.tags.day", "Sunday"),
				),
			},
		},
	})
}
