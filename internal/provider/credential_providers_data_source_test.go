package provider

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testCredentialProvidersDataSource          string = "data.aembit_credential_providers.test"
	testCredentialProviderResource             string = "aembit_credential_provider.oauth"
	testAzureKeyVaultValueCpDataSource         string = "data.aembit_credential_providers.azure_key_vault_value_test"
	testAzureKeyVaultValueCpResource           string = "aembit_credential_provider.azure_key_vault_value_cp"
	testMcpUserBasedAccessTokenCpDataSource    string = "data.aembit_credential_providers.mcp_user_based_access_token_test"
	testMcpUserBasedAccessTokenCpDataSourceRes string = "aembit_credential_provider.mcp_user_based_access_token_cp"
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

		resourceSetID := rs.Primary.Attributes["resource_set_id"]

		if _, err, notFound = testClient.GetCredentialProviderV2(rs.Primary.ID, nil, &resourceSetID); notFound {
			return err
		}
		return nil
	}
}

func TestAccCredentialProvidersDataSource(t *testing.T) {
	createFile1, _ := os.ReadFile("../../tests/credential/data/TestAccCredentialProvidersDataSource_ResourceSet.tf")
	createFile2, _ := os.ReadFile("../../tests/credential/data/TestAccCredentialProvidersDataSource_ProviderResourceSet.tf")
	createFile3, _ := os.ReadFile("../../tests/credential/data/TestAccCredentialProvidersDataSource.tf")

	files := [3]string{string(createFile1), string(createFile2), string(createFile3)}

	for _, createFile := range files {
		createFileConfig, _, _ := randomizeFileConfigs(createFile, "", "TF Acceptance OAuth")

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

func TestAccMcpUserBasedAccessTokenCredentialProvidersDataSource(t *testing.T) {
	createFile, _ := os.ReadFile(
		"../../tests/credential/data/TestAccMcpUserBasedAccessTokenCpDataSource.tf",
	)
	createFileConfig, _, endName := randomizeFileConfigs(
		string(createFile),
		"",
		"TF Acceptance McpUserBasedAccessToken DataSource",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						testMcpUserBasedAccessTokenCpDataSource,
						"credential_providers.#",
					),
					testFindCredentialProvider(testMcpUserBasedAccessTokenCpDataSourceRes),
					func(s *terraform.State) error {
						ms, ok := s.RootModule().Resources[testMcpUserBasedAccessTokenCpDataSource]
						if !ok {
							return fmt.Errorf("Not found: %s", testMcpUserBasedAccessTokenCpDataSource)
						}

						countStr := ms.Primary.Attributes["credential_providers.#"]
						count, err := strconv.Atoi(countStr)
						if err != nil {
							return err
						}

						const expectedRevocationUrl = "https://aembit.io/revoke"
						for i := 0; i < count; i++ {
							if ms.Primary.Attributes[fmt.Sprintf("credential_providers.%d.name", i)] == endName {
								actual := ms.Primary.Attributes[fmt.Sprintf(
									"credential_providers.%d.mcp_user_based_access_token.oauth_revocation_url",
									i,
								)]
								if actual != expectedRevocationUrl {
									return fmt.Errorf(
										"expected oauth_revocation_url %q, got %q",
										expectedRevocationUrl,
										actual,
									)
								}
								return nil
							}
						}

						return fmt.Errorf(
							"Credential Provider with name %q not found in %s",
							endName,
							testMcpUserBasedAccessTokenCpDataSource,
						)
					},
				),
			},
		},
	})
}
