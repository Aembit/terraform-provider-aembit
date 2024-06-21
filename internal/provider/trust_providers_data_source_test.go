package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testTrustProvidersDataSource string = "data.aembit_trust_providers.test"

func TestAccTrustProvidersDataSource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/trust/data/TestAccTrustProvidersDataSource.tf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of Trust Providers returned
					resource.TestCheckResourceAttrSet(testTrustProvidersDataSource, "trust_providers.#"),
					// Verify First Trust Provider Name
					resource.TestCheckResourceAttr(testTrustProvidersDataSource, "trust_providers.0.name", "TF Acceptance Kubernetes"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testTrustProvidersDataSource, "trust_providers.0.id"),
					// Verify Tags.
					resource.TestCheckResourceAttr(testTrustProvidersDataSource, "trust_providers.0.tags.%", "2"),
					resource.TestCheckResourceAttr(testTrustProvidersDataSource, "trust_providers.0.tags.color", "blue"),
					resource.TestCheckResourceAttr(testTrustProvidersDataSource, "trust_providers.0.tags.day", "Sunday"),
				),
			},
		},
	})
}
