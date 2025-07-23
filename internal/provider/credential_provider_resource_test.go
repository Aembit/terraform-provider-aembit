package provider

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testCredentialProviderAembit    = "aembit_credential_provider.aembit"
	testCredentialProviderApiKey    = "aembit_credential_provider.api_key"
	testCredentialProviderAWS       = "aembit_credential_provider.aws"
	testCredentialProviderGCP       = "aembit_credential_provider.gcp"
	testCredentialProviderSnowflake = "aembit_credential_provider.snowflake"
	testCredentialProviderUserPass  = "aembit_credential_provider.userpass"
	testCredentialProviderVault     = "aembit_credential_provider.vault"
)

func testDeleteCredentialProvider(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var ok bool
		var err error
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if ok, err = testClient.DeleteCredentialProviderV2(context.Background(), rs.Primary.ID, nil); !ok {
			return err
		}
		return nil
	}
}

func TestAccCredentialProviderResource_AembitToken(t *testing.T) {
	createFile, _ := os.ReadFile(
		"../../tests/credential/aembit/TestAccCredentialProviderResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/credential/aembit/TestAccCredentialProviderResource.tfmod",
	)
	createFileConfig, modifyFileConfig, _ := randomizeFileConfigs(
		string(createFile),
		string(modifyFile),
		"TF Acceptance Role for Token",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider set values
					resource.TestCheckResourceAttr(
						testCredentialProviderAembit,
						"name",
						"TF Acceptance Aembit Token",
					),
					resource.TestCheckResourceAttr(
						testCredentialProviderAembit,
						"aembit_access_token.lifetime",
						"1800",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderAembit, "id"),
					resource.TestCheckResourceAttrSet(
						testCredentialProviderAembit,
						"aembit_access_token.audience",
					),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testCredentialProviderAembit, "id"),
				),
			},
			// Test Aembit API Removal causes re-create with non-empty plan
			{
				Config:             createFileConfig,
				Check:              testDeleteCredentialProvider(testCredentialProviderAembit),
				ExpectNonEmptyPlan: true,
			},
			// Recreate the resource from the first test step
			{Config: createFileConfig},
			// ImportState testing
			{
				ResourceName:      testCredentialProviderAembit,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						testCredentialProviderAembit,
						"name",
						"TF Acceptance Aembit Token - Modified",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(
						testCredentialProviderAembit,
						"aembit_access_token.audience",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderResource_ApiKey(t *testing.T) {
	createFile, _ := os.ReadFile(
		"../../tests/credential/apikey/TestAccCredentialProviderResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/credential/apikey/TestAccCredentialProviderResource.tfmod",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr(
						testCredentialProviderApiKey,
						"name",
						"TF Acceptance API Key",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderApiKey, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testCredentialProviderApiKey, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      testCredentialProviderApiKey,
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						testCredentialProviderApiKey,
						"name",
						"TF Acceptance API Key - Modified",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderResource_AwsSTS(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/aws/TestAccCredentialProviderResource.tf")
	modifyFile, _ := os.ReadFile(
		"../../tests/credential/aws/TestAccCredentialProviderResource.tfmod",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential set values
					resource.TestCheckResourceAttr(
						testCredentialProviderAWS,
						"name",
						"TF Acceptance AWS STS",
					),
					resource.TestCheckResourceAttr(
						testCredentialProviderAWS,
						"aws_sts.lifetime",
						"1800",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderAWS, "id"),
					resource.TestCheckResourceAttrSet(
						testCredentialProviderAWS,
						"aws_sts.oidc_issuer",
					),
					resource.TestCheckResourceAttrSet(
						testCredentialProviderAWS,
						"aws_sts.token_audience",
					),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testCredentialProviderAWS, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      testCredentialProviderAWS,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						testCredentialProviderAWS,
						"name",
						"TF Acceptance AWS STS - Modified",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderAWS, "id"),
					resource.TestCheckResourceAttrSet(
						testCredentialProviderAWS,
						"aws_sts.oidc_issuer",
					),
					resource.TestCheckResourceAttrSet(
						testCredentialProviderAWS,
						"aws_sts.token_audience",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderResource_GoogleWorkload(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/gcp/TestAccCredentialProviderResource.tf")
	modifyFile, _ := os.ReadFile(
		"../../tests/credential/gcp/TestAccCredentialProviderResource.tfmod",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider set values
					resource.TestCheckResourceAttr(
						testCredentialProviderGCP,
						"name",
						"TF Acceptance GCP Workload",
					),
					resource.TestCheckResourceAttr(
						testCredentialProviderGCP,
						"google_workload_identity.lifetime",
						"1800",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderGCP, "id"),
					resource.TestCheckResourceAttrSet(
						testCredentialProviderGCP,
						"google_workload_identity.oidc_issuer",
					),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testCredentialProviderGCP, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      testCredentialProviderGCP,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						testCredentialProviderGCP,
						"name",
						"TF Acceptance GCP Workload - Modified",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderGCP, "id"),
					resource.TestCheckResourceAttrSet(
						testCredentialProviderGCP,
						"google_workload_identity.oidc_issuer",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderResource_AzureEntraToken(t *testing.T) {
	const credentialProviderName string = "aembit_credential_provider.ae"
	createFile, _ := os.ReadFile(
		"../../tests/credential/azure-entra/TestAccCredentialProviderResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/credential/azure-entra/TestAccCredentialProviderResource.tfmod",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider set values
					resource.TestCheckResourceAttr(
						credentialProviderName,
						"name",
						"TF Acceptance Azure Entra Workload",
					),
					resource.TestCheckResourceAttr(
						credentialProviderName,
						"azure_entra_workload_identity.audience",
						"audience",
					),
					resource.TestCheckResourceAttr(
						credentialProviderName,
						"azure_entra_workload_identity.subject",
						"subject",
					),
					resource.TestCheckResourceAttr(
						credentialProviderName,
						"azure_entra_workload_identity.scope",
						"scope",
					),
					resource.TestCheckResourceAttr(
						credentialProviderName,
						"azure_entra_workload_identity.azure_tenant",
						"00000000-0000-0000-0000-000000000000",
					),
					resource.TestCheckResourceAttr(
						credentialProviderName,
						"azure_entra_workload_identity.client_id",
						"00000000-0000-0000-0000-000000000000",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(
						credentialProviderName,
						"azure_entra_workload_identity.oidc_issuer",
					),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(credentialProviderName, "id"),
				),
			},
			// ImportState testing
			{ResourceName: credentialProviderName, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						credentialProviderName,
						"name",
						"TF Acceptance Azure Entra Workload - Modified",
					),
					resource.TestCheckResourceAttr(
						credentialProviderName,
						"azure_entra_workload_identity.audience",
						"new audience",
					),
					resource.TestCheckResourceAttr(
						credentialProviderName,
						"azure_entra_workload_identity.subject",
						"new subject",
					),
					resource.TestCheckResourceAttr(
						credentialProviderName,
						"azure_entra_workload_identity.scope",
						"new scope",
					),
					resource.TestCheckResourceAttr(
						credentialProviderName,
						"azure_entra_workload_identity.azure_tenant",
						"11111111-1111-1111-1111-111111111111",
					),
					resource.TestCheckResourceAttr(
						credentialProviderName,
						"azure_entra_workload_identity.client_id",
						"11111111-1111-1111-1111-111111111111",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(
						credentialProviderName,
						"azure_entra_workload_identity.oidc_issuer",
					),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(credentialProviderName, "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderResource_SnowflakeToken(t *testing.T) {
	createFile, _ := os.ReadFile(
		"../../tests/credential/snowflake/TestAccCredentialProviderResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/credential/snowflake/TestAccCredentialProviderResource.tfmod",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr(
						testCredentialProviderSnowflake,
						"name",
						"TF Acceptance Snowflake Token",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderSnowflake, "id"),
					resource.TestCheckResourceAttrSet(
						testCredentialProviderSnowflake,
						"snowflake_jwt.alter_user_command",
					),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testCredentialProviderSnowflake, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      testCredentialProviderSnowflake,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						testCredentialProviderSnowflake,
						"name",
						"TF Acceptance Snowflake Token - Modified",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderSnowflake, "id"),
					resource.TestCheckResourceAttrSet(
						testCredentialProviderSnowflake,
						"snowflake_jwt.alter_user_command",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

const (
	testOAuthClientCredentialsAuthHeader string = "aembit_credential_provider.oauth_authHeader"
	testOAuthClientCredentialsPostBody   string = "aembit_credential_provider.oauth_postBody"
)

func TestAccCredentialProviderResource_OAuthClientCredentialsAuthHeader(t *testing.T) {
	createFile, _ := os.ReadFile(
		"../../tests/credential/oauth-client-credentials/TestAccCredentialProviderResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/credential/oauth-client-credentials/TestAccCredentialProviderResource.tfmod",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr(
						testOAuthClientCredentialsAuthHeader,
						"name",
						"TF Acceptance OAuth",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testOAuthClientCredentialsAuthHeader, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testOAuthClientCredentialsAuthHeader, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      testOAuthClientCredentialsAuthHeader,
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						testOAuthClientCredentialsAuthHeader,
						"name",
						"TF Acceptance OAuth - Modified",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderResource_OAuthClientCredentialsPostBody(t *testing.T) {
	createFile, _ := os.ReadFile(
		"../../tests/credential/oauth-client-credentials/TestAccCredentialProviderResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/credential/oauth-client-credentials/TestAccCredentialProviderResource.tfmod",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr(
						testOAuthClientCredentialsPostBody,
						"name",
						"TF Acceptance OAuth",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testOAuthClientCredentialsPostBody, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testOAuthClientCredentialsPostBody, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      testOAuthClientCredentialsPostBody,
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						testOAuthClientCredentialsPostBody,
						"name",
						"TF Acceptance OAuth - Modified",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

const (
	testOAuthAuthCodeResource              string = "aembit_credential_provider.oauth_authorization_code"
	testOAuthAuthCodeEmptyCustomParameters string = "aembit_credential_provider.oauth_authorization_code_empty_custom_parameters"
	testOAuthAuthCodeUserAuthUrl                  = "oauth_authorization_code.user_authorization_url"
)

func TestAccCredentialProviderResource_OAuthAuthorizationCode(t *testing.T) {
	createFile, _ := os.ReadFile(
		"../../tests/credential/oauth-authorization-code/TestAccCredentialProviderResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/credential/oauth-authorization-code/TestAccCredentialProviderResource.tfmod",
	)

	firstID := uuid.New().String()
	secondID := uuid.New().String()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: strings.ReplaceAll(
					strings.ReplaceAll(string(createFile), "replace-with-uuid-first", firstID),
					"replace-with-uuid-second",
					secondID,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr(
						testOAuthAuthCodeResource,
						"name",
						"TF Acceptance OAuth Authorization Code",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testOAuthAuthCodeResource, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testOAuthAuthCodeResource, "id"),
					// Verify we get back a user_authorization_url
					resource.TestCheckResourceAttrSet(
						testOAuthAuthCodeResource,
						testOAuthAuthCodeUserAuthUrl,
					),
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr(
						testOAuthAuthCodeEmptyCustomParameters,
						"name",
						"TF Acceptance OAuth Authorization Code",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testOAuthAuthCodeEmptyCustomParameters, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testOAuthAuthCodeEmptyCustomParameters, "id"),
					// Verify we get back a user_authorization_url
					resource.TestCheckResourceAttrSet(
						testOAuthAuthCodeEmptyCustomParameters,
						testOAuthAuthCodeUserAuthUrl,
					),
				),
			},
			// ImportState testing
			{ResourceName: testOAuthAuthCodeResource, ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						testOAuthAuthCodeResource,
						"name",
						"TF Acceptance OAuth Authorization Code - Modified",
					),
					// Verify we get back a user_authorization_url
					resource.TestCheckResourceAttrSet(
						testOAuthAuthCodeResource,
						testOAuthAuthCodeUserAuthUrl,
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderResource_UsernamePassword(t *testing.T) {
	createFile, _ := os.ReadFile(
		"../../tests/credential/userpass/TestAccCredentialProviderResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/credential/userpass/TestAccCredentialProviderResource.tfmod",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr(
						testCredentialProviderUserPass,
						"name",
						"TF Acceptance Username Password",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderUserPass, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testCredentialProviderUserPass, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      testCredentialProviderUserPass,
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						testCredentialProviderUserPass,
						"name",
						"TF Acceptance Username Password - Modified",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderResource_VaultClientToken(t *testing.T) {
	createFile, _ := os.ReadFile(
		"../../tests/credential/vault/TestAccCredentialProviderResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/credential/vault/TestAccCredentialProviderResource.tfmod",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr(
						testCredentialProviderVault,
						"name",
						"TF Acceptance Vault",
					),
					// Verify Tags.
					resource.TestCheckResourceAttr(
						testCredentialProviderVault,
						tagsCount,
						"2",
					),
					resource.TestCheckResourceAttr(
						testCredentialProviderVault,
						tagsColor,
						"blue",
					),
					resource.TestCheckResourceAttr(
						testCredentialProviderVault,
						tagsDay,
						"Sunday",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderVault, "id"),
					resource.TestCheckResourceAttr(
						testCredentialProviderVault,
						"vault_client_token.vault_forwarding",
						"",
					),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testCredentialProviderVault, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      testCredentialProviderVault,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						testCredentialProviderVault,
						"name",
						"TF Acceptance Vault - Modified",
					),
					// Verify Tags.
					resource.TestCheckResourceAttr(
						testCredentialProviderVault,
						tagsCount,
						"2",
					),
					resource.TestCheckResourceAttr(
						testCredentialProviderVault,
						tagsColor,
						"orange",
					),
					resource.TestCheckResourceAttr(
						testCredentialProviderVault,
						tagsDay,
						"Tuesday",
					),
					// Verify Vault_Forwarding update
					resource.TestCheckResourceAttr(
						testCredentialProviderVault,
						"vault_client_token.vault_forwarding",
						"conditional",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

var gitlabManagedAccountResourcePath = "aembit_credential_provider.gitlab_managed_account"

func TestAccCredentialProviderResource_ManagedGitlabAccount(t *testing.T) {
	t.Skip("skipping test until we figure out a way to handle the GitLab tokens appropriately")

	createFile, _ := os.ReadFile(
		"../../tests/credential/gitlab-managed-account/TestAccCredentialProviderResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/credential/gitlab-managed-account/TestAccCredentialProviderResource.tfmod",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr(
						gitlabManagedAccountResourcePath,
						"name",
						"TF Acceptance Managed Gitlab Account",
					),
					resource.TestCheckResourceAttr(
						gitlabManagedAccountResourcePath,
						"managed_gitlab_account.service_account_username",
						"test_service_account",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(gitlabManagedAccountResourcePath, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(gitlabManagedAccountResourcePath, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      gitlabManagedAccountResourcePath,
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						gitlabManagedAccountResourcePath,
						"name",
						"TF Acceptance Managed Gitlab Account - Updated",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

var oidcIdTokenResourcePath = "aembit_credential_provider.oidc_id_token"

func TestAccCredentialProviderResource_OidcIdToken(t *testing.T) {
	createFile, _ := os.ReadFile(
		"../../tests/credential/oidc-id-token/TestAccCredentialProviderResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/credential/oidc-id-token/TestAccCredentialProviderResource.tfmod",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr(
						oidcIdTokenResourcePath,
						"name",
						"TF Acceptance OIDC ID Token",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(oidcIdTokenResourcePath, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(oidcIdTokenResourcePath, "id"),
				),
			},
			// ImportState testing
			{ResourceName: oidcIdTokenResourcePath, ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						oidcIdTokenResourcePath,
						"name",
						"TF Acceptance OIDC ID Token - Modified",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccAwsSecretsManagerValueCP(t *testing.T) {
	cpResourcePath := "aembit_credential_provider.aws_sm_value"
	createFile, _ := os.ReadFile(
		"../../tests/credential/aws-secrets-manager/TestAwsSecretsManagerValueResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/credential/aws-secrets-manager/TestAwsSecretsManagerValueResource.tfmod",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr(
						cpResourcePath,
						"name",
						"TF Acceptance AWS Secrets Manager Value CP",
					),
					resource.TestCheckResourceAttr(
						cpResourcePath,
						"description",
						"TF Acceptance AWS Secrets Manager Value CP Description",
					),
					resource.TestCheckResourceAttr(
						cpResourcePath,
						"aws_secrets_manager_value.secret_arn",
						"arn:aws:secretsmanager:us-east-2:123456789012:secret:secretname-ABCDEF",
					),
					resource.TestCheckResourceAttr(
						cpResourcePath,
						"aws_secrets_manager_value.secret_key_1",
						"key1",
					),
					resource.TestCheckResourceAttr(
						cpResourcePath,
						"aws_secrets_manager_value.secret_key_2",
						"key2",
					),
					resource.TestCheckResourceAttr(
						cpResourcePath,
						"aws_secrets_manager_value.private_network_access",
						"false",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(cpResourcePath, "id"),
				),
			},
			// ImportState testing
			{ResourceName: cpResourcePath, ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						cpResourcePath,
						"name",
						"TF Acceptance AWS Secrets Manager Value CP - Updated",
					),
					resource.TestCheckResourceAttr(
						cpResourcePath,
						"description",
						"TF Acceptance AWS Secrets Manager Value CP Description - Updated",
					),
					resource.TestCheckResourceAttr(
						cpResourcePath,
						"aws_secrets_manager_value.secret_arn",
						"arn:aws:secretsmanager:us-east-2:123456789012:secret:anothersecretname-ABCDEF",
					),
					resource.TestCheckResourceAttr(
						cpResourcePath,
						"aws_secrets_manager_value.secret_key_1",
						"key1-updated",
					),
					resource.TestCheckResourceAttr(
						cpResourcePath,
						"aws_secrets_manager_value.secret_key_2",
						"key2-updated",
					),
					resource.TestCheckResourceAttr(
						cpResourcePath,
						"aws_secrets_manager_value.private_network_access",
						"true",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
