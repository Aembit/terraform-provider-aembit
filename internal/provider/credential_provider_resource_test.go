package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCredentialProviderResource_ApiKey(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "aembit_credential_provider" "test" {
	name = "Unit Test 1"
	type = "apikey"
	api_key {
		value = "test"
	}
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr("aembit_credential_provider.test", "name", "Unit Test 1"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.test", "id"),
					resource.TestCheckResourceAttrSet("aembit_credential_provider.test", "type"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_credential_provider.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_credential_provider.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "aembit_credential_provider" "test" {
	name = "Unit Test 1 - Modified"
	type = "apikey"
	api_key {
		value = "test"
	}
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_credential_provider.test", "name", "Unit Test 1 - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
