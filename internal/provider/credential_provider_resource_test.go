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

const testCredentialProviderAembit string = "aembit_credential_provider.aembit"

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
						"aembit_credential_provider.api_key",
						"name",
						"TF Acceptance API Key",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.api_key", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_credential_provider.api_key", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_credential_provider.api_key",
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						"aembit_credential_provider.api_key",
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
						"aembit_credential_provider.aws",
						"name",
						"TF Acceptance AWS STS",
					),
					resource.TestCheckResourceAttr(
						"aembit_credential_provider.aws",
						"aws_sts.lifetime",
						"1800",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.aws", "id"),
					resource.TestCheckResourceAttrSet(
						"aembit_credential_provider.aws",
						"aws_sts.oidc_issuer",
					),
					resource.TestCheckResourceAttrSet(
						"aembit_credential_provider.aws",
						"aws_sts.token_audience",
					),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_credential_provider.aws", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_credential_provider.aws",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						"aembit_credential_provider.aws",
						"name",
						"TF Acceptance AWS STS - Modified",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.aws", "id"),
					resource.TestCheckResourceAttrSet(
						"aembit_credential_provider.aws",
						"aws_sts.oidc_issuer",
					),
					resource.TestCheckResourceAttrSet(
						"aembit_credential_provider.aws",
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
						"aembit_credential_provider.gcp",
						"name",
						"TF Acceptance GCP Workload",
					),
					resource.TestCheckResourceAttr(
						"aembit_credential_provider.gcp",
						"google_workload_identity.lifetime",
						"1800",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.gcp", "id"),
					resource.TestCheckResourceAttrSet(
						"aembit_credential_provider.gcp",
						"google_workload_identity.oidc_issuer",
					),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_credential_provider.gcp", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_credential_provider.gcp",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						"aembit_credential_provider.gcp",
						"name",
						"TF Acceptance GCP Workload - Modified",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.gcp", "id"),
					resource.TestCheckResourceAttrSet(
						"aembit_credential_provider.gcp",
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
						"aembit_credential_provider.snowflake",
						"name",
						"TF Acceptance Snowflake Token",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.snowflake", "id"),
					resource.TestCheckResourceAttrSet(
						"aembit_credential_provider.snowflake",
						"snowflake_jwt.alter_user_command",
					),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_credential_provider.snowflake", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_credential_provider.snowflake",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						"aembit_credential_provider.snowflake",
						"name",
						"TF Acceptance Snowflake Token - Modified",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.snowflake", "id"),
					resource.TestCheckResourceAttrSet(
						"aembit_credential_provider.snowflake",
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
						"aembit_credential_provider.userpass",
						"name",
						"TF Acceptance Username Password",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.userpass", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_credential_provider.userpass", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_credential_provider.userpass",
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						"aembit_credential_provider.userpass",
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
						"aembit_credential_provider.vault",
						"name",
						"TF Acceptance Vault",
					),
					// Verify Tags.
					resource.TestCheckResourceAttr(
						"aembit_credential_provider.vault",
						tagsCount,
						"2",
					),
					resource.TestCheckResourceAttr(
						"aembit_credential_provider.vault",
						tagsColor,
						"blue",
					),
					resource.TestCheckResourceAttr(
						"aembit_credential_provider.vault",
						tagsDay,
						"Sunday",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.vault", "id"),
					resource.TestCheckResourceAttr(
						"aembit_credential_provider.vault",
						"vault_client_token.vault_forwarding",
						"",
					),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_credential_provider.vault", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_credential_provider.vault",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						"aembit_credential_provider.vault",
						"name",
						"TF Acceptance Vault - Modified",
					),
					// Verify Tags.
					resource.TestCheckResourceAttr(
						"aembit_credential_provider.vault",
						tagsCount,
						"2",
					),
					resource.TestCheckResourceAttr(
						"aembit_credential_provider.vault",
						tagsColor,
						"orange",
					),
					resource.TestCheckResourceAttr(
						"aembit_credential_provider.vault",
						tagsDay,
						"Tuesday",
					),
					// Verify Vault_Forwarding update
					resource.TestCheckResourceAttr(
						"aembit_credential_provider.vault",
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
