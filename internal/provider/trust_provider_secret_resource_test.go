package provider

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	trustProviderSecretCertificatePath   = "aembit_trust_provider_secret.secret_certificate"
	trustProviderSecretSymmetricKeyPath1 = "aembit_trust_provider_secret.symmetric_key_secret1"
	trustProviderSecretSymmetricKeyPath2 = "aembit_trust_provider_secret.symmetric_key_secret2"
)

func TestAccTrustProviderSecretResource_Certificate(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile(
		"../../tests/trust_provider_secret/certificate/TestAccTrustProviderSecretResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/trust_provider_secret/certificate/TestAccTrustProviderSecretResource.tfmod",
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
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderSecretCertificatePath, "id"),
					// Verify Subject
					resource.TestCheckResourceAttr(
						trustProviderSecretCertificatePath,
						"subject",
						"CN=Aembit Unit Test",
					),
					// Verify type
					resource.TestCheckResourceAttr(
						trustProviderSecretCertificatePath,
						"type",
						"Certificate",
					),
				),
			},
			// ImportState testing
			{
				ResourceName:      trustProviderSecretCertificatePath,
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: modifyFileString,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Subject updated
					resource.TestCheckResourceAttr(
						trustProviderSecretCertificatePath,
						"subject",
						"CN=Updated Aembit Unit Test",
					),
					// Verify type
					resource.TestCheckResourceAttr(
						trustProviderSecretCertificatePath,
						"type",
						"Certificate",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTrustProviderSecretResource_SymmetricKey(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile(
		"../../tests/trust_provider_secret/symmetric_key/TestAccTrustProviderSecretResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/trust_provider_secret/symmetric_key/TestAccTrustProviderSecretResource.tfmod",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderSecretSymmetricKeyPath1, "id"),
					// Verify type
					resource.TestCheckResourceAttr(
						trustProviderSecretSymmetricKeyPath1,
						"type",
						"SymmetricKey",
					),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderSecretSymmetricKeyPath2, "id"),
					// Verify type
					resource.TestCheckResourceAttr(
						trustProviderSecretSymmetricKeyPath2,
						"type",
						"SymmetricKey",
					),
				),
			},
			// ImportState testing
			{
				ResourceName:      trustProviderSecretSymmetricKeyPath1,
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify type
					resource.TestCheckResourceAttr(
						trustProviderSecretSymmetricKeyPath1,
						"type",
						"SymmetricKey",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
