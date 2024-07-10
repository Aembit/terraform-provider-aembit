package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testCredentialProvidersDataSource string = "data.aembit_credential_providers.test"

func TestAccCredentialProvidersDataSource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/data/TestAccCredentialProvidersDataSource.tf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of Credential Providers returned
					resource.TestCheckResourceAttrSet(testCredentialProvidersDataSource, "credential_providers.#"),
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr(testCredentialProvidersDataSource, "credential_providers.0.name", "TF Acceptance OAuth"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProvidersDataSource, "credential_providers.0.id"),
				),
			},
		},
	})
}