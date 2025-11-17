package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	trustProviderSecretPath = "aembit_trust_provider_secret.secret"
)

func TestAccTrustProviderSecretResource(t *testing.T) {
	createFile, _ := os.ReadFile(
		"../../tests/trust_provider_secret/TestAccTrustProviderSecretResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/trust_provider_secret/TestAccTrustProviderSecretResource.tfmod",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderSecretPath, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderSecretPath, "id"),
					// Verify Subject
					resource.TestCheckResourceAttr(
						trustProviderSecretPath,
						"subject",
						"CN=Test Certificate",
					),
				),
			},
			// ImportState testing
			{ResourceName: trustProviderSecretPath, ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Subject updated
					resource.TestCheckResourceAttr(
						trustProviderSecretPath,
						"subject",
						"CN=Different Test Certificate",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
