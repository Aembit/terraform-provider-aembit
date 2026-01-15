package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testCredentialProvidersDataSource  string = "data.aembit_credential_providers.test"
	testCredentialProviderResource     string = "aembit_credential_provider.oauth"
	testAzureKeyVaultValueCpDataSource string = "data.aembit_credential_providers.azure_key_vault_value_test"
	testAzureKeyVaultValueCpResource   string = "aembit_credential_provider.azure_key_vault_value_cp"
)

func testFindCredentialProvider(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var err error
		var ok bool
		var notFound bool
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if _, err, notFound = testClient.GetCredentialProviderV2(rs.Primary.ID, nil); notFound {
			return err
		}
		return nil
	}
}

func TestAccCredentialProvidersDataSource(t *testing.T) {
	createFile, _ := os.ReadFile(
		"../../tests/credential/data/TestAccCredentialProvidersDataSource.tf",
	)
	createFileConfig, _, _ := randomizeFileConfigs(string(createFile), "", "TF Acceptance OAuth")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify non-zero number of Credential Providers returned
					resource.TestCheckResourceAttrSet(
						testCredentialProvidersDataSource,
						"credential_providers.#",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(
						testCredentialProvidersDataSource,
						"credential_providers.0.id",
					),
					// Find newly created entry
					testFindCredentialProvider(testCredentialProviderResource),
				),
			},
		},
	})
}

func TestAccAzureKeyVaultValueCredentialProvidersDataSource(t *testing.T) {
	createFile, _ := os.ReadFile(
		"../../tests/credential/data/TestAccAzureKeyVaultCpDataSource.tf",
	)
	createFileConfig, _, _ := randomizeFileConfigs(
		string(createFile),
		"",
		"TF Acceptance Azure Key Vault Value",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify non-zero number of Credential Providers returned
					resource.TestCheckResourceAttrSet(
						testAzureKeyVaultValueCpDataSource,
						"credential_providers.#",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(
						testAzureKeyVaultValueCpDataSource,
						"credential_providers.0.id",
					),
					// Find newly created entry
					testFindCredentialProvider(testAzureKeyVaultValueCpResource),
				),
			},
		},
	})
}
