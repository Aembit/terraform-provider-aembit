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
resource "aembit_credential_provider" "api_key" {
	name = "TF Acceptance API Key"
	api_key = {
		api_key = "test"
	}
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr("aembit_credential_provider.api_key", "name", "TF Acceptance API Key"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.api_key", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_credential_provider.api_key", "id"),
				),
			},
			// ImportState testing
			//{
			//	ResourceName:      "aembit_credential_provider.api_key",
			//	ImportState:       false,
			//	ImportStateVerify: false,
			//},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "aembit_credential_provider" "api_key" {
	name = "TF Acceptance API Key - Modified"
	api_key = {
		api_key = "test"
	}
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_credential_provider.api_key", "name", "TF Acceptance API Key - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderResource_OAuthClientCredentials(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "aembit_credential_provider" "oauth" {
	name = "TF Acceptance OAuth"
	oauth_client_credentials = {
		token_url = "https://aembit.io/token"
		client_id = "test"
		client_secret = "test"
		scopes = "test_scope"
	}
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr("aembit_credential_provider.oauth", "name", "TF Acceptance OAuth"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.oauth", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_credential_provider.oauth", "id"),
				),
			},
			// ImportState testing
			//{
			//	ResourceName:      "aembit_credential_provider.oauth",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "aembit_credential_provider" "oauth" {
	name = "TF Acceptance OAuth - Modified"
	oauth_client_credentials = {
		token_url = "https://aembit.io/token"
		client_id = "test"
		client_secret = "test"
		scopes = ""
	}
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_credential_provider.oauth", "name", "TF Acceptance OAuth - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
