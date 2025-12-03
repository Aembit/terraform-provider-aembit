package provider

import (
	"os"
	"strings"
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
	// Generate random PEM certificates for creation and update
	pemCertificate := generateRandomPEMCertificate(t, "Aembit Unit Test")
	updatedPemCertificate := generateRandomPEMCertificate(t, "Updated Aembit Unit Test")

	// Replace placeholder with generated certificate
	createFileString := strings.ReplaceAll(string(createFile), "PEM_CERTIFICATE", pemCertificate)
	modifyFileString := strings.ReplaceAll(string(modifyFile), "PEM_CERTIFICATE",
		updatedPemCertificate,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileString,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderSecretPath, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderSecretPath, "id"),
					// Verify Subject
					resource.TestCheckResourceAttr(
						trustProviderSecretPath,
						"subject",
						"CN=Aembit Unit Test",
					),
				),
			},
			// ImportState testing
			{ResourceName: trustProviderSecretPath, ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: modifyFileString,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Subject updated
					resource.TestCheckResourceAttr(
						trustProviderSecretPath,
						"subject",
						"CN=Updated Aembit Unit Test",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
