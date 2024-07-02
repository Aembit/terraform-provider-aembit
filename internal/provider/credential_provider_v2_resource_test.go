package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testCredentialProviderV2Aembit string = "aembit_credential_provider_v2.aembit"
const testCredentialProviderV2ApiKey string = "aembit_credential_provider_v2.api_key"
const testCredentialProviderV2Aws string = "aembit_credential_provider_v2.aws"
const testCredentialProviderV2Gcp string = "aembit_credential_provider_v2.gcp"
const testCredentialProviderV2Snowflake string = "aembit_credential_provider_v2.snowflake"
const testOAuthClientCredentialsV2AuthHeader string = "aembit_credential_provider_v2.oauth_authHeader"
const testOAuthClientCredentialsV2PostBody string = "aembit_credential_provider_v2.oauth_postBody"
const testCredentialProviderV2Username string = "aembit_credential_provider_v2.userpass"
const testCredentialProviderV2Vault string = "aembit_credential_provider_v2.vault"

func testDeleteCredentialProviderV2(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var ok bool
		var err error
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if ok, err = testClient.DeleteCredentialProviderV2(rs.Primary.ID, nil); !ok {
			return err
		}
		return nil
	}
}

func TestAccCredentialProviderV2Resource_AembitToken(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/aembit/TestAccCredentialProviderV2Resource.tf")
	modifyFile, _ := os.ReadFile("../../tests/credential/aembit/TestAccCredentialProviderV2Resource.tfmod")
	createFileConfig, modifyFileConfig, _ := randomizeFileConfigs(string(createFile), string(modifyFile), "IGNORE")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider set values
					resource.TestCheckResourceAttr(testCredentialProviderV2Aembit, "name", "TF Acceptance Aembit Token"),
					resource.TestCheckResourceAttr(testCredentialProviderV2Aembit, "aembit_access_token.lifetime", "1800"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Aembit, "id"),
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Aembit, "aembit_access_token.audience"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Aembit, "id"),
				),
			},
			// Test Aembit API Removal causes re-create with non-empty plan
			{Config: createFileConfig, Check: testDeleteCredentialProviderV2(testCredentialProviderV2Aembit), ExpectNonEmptyPlan: true},
			// Recreate the resource from the first test step
			{Config: createFileConfig},
			// ImportState testing
			{ResourceName: testCredentialProviderV2Aembit, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(testCredentialProviderV2Aembit, "name", "TF Acceptance Aembit Token - Modified"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Aembit, "aembit_access_token.audience"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderV2Resource_ApiKey(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/apikey/TestAccCredentialProviderV2Resource.tf")
	modifyFile, _ := os.ReadFile("../../tests/credential/apikey/TestAccCredentialProviderV2Resource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr(testCredentialProviderV2ApiKey, "name", "TF Acceptance API Key"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderV2ApiKey, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testCredentialProviderV2ApiKey, "id"),
				),
			},
			// ImportState testing
			{ResourceName: testCredentialProviderV2ApiKey, ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(testCredentialProviderV2ApiKey, "name", "TF Acceptance API Key - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderV2Resource_AwsSTS(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/aws/TestAccCredentialProviderV2Resource.tf")
	modifyFile, _ := os.ReadFile("../../tests/credential/aws/TestAccCredentialProviderV2Resource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential set values
					resource.TestCheckResourceAttr(testCredentialProviderV2Aws, "name", "TF Acceptance AWS STS"),
					resource.TestCheckResourceAttr(testCredentialProviderV2Aws, "aws_sts.lifetime", "1800"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Aws, "id"),
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Aws, "aws_sts.oidc_issuer"),
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Aws, "aws_sts.token_audience"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Aws, "id"),
				),
			},
			// ImportState testing
			{ResourceName: testCredentialProviderV2Aws, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(testCredentialProviderV2Aws, "name", "TF Acceptance AWS STS - Modified"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Aws, "id"),
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Aws, "aws_sts.oidc_issuer"),
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Aws, "aws_sts.token_audience"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderV2Resource_GoogleWorkload(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/gcp/TestAccCredentialProviderV2Resource.tf")
	modifyFile, _ := os.ReadFile("../../tests/credential/gcp/TestAccCredentialProviderV2Resource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider set values
					resource.TestCheckResourceAttr(testCredentialProviderV2Gcp, "name", "TF Acceptance GCP Workload"),
					resource.TestCheckResourceAttr(testCredentialProviderV2Gcp, "google_workload_identity.lifetime", "1800"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Gcp, "id"),
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Gcp, "google_workload_identity.oidc_issuer"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Gcp, "id"),
				),
			},
			// ImportState testing
			{ResourceName: testCredentialProviderV2Gcp, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(testCredentialProviderV2Gcp, "name", "TF Acceptance GCP Workload - Modified"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Gcp, "id"),
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Gcp, "google_workload_identity.oidc_issuer"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderV2Resource_SnowflakeToken(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/snowflake/TestAccCredentialProviderV2Resource.tf")
	modifyFile, _ := os.ReadFile("../../tests/credential/snowflake/TestAccCredentialProviderV2Resource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr(testCredentialProviderV2Snowflake, "name", "TF Acceptance Snowflake Token"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Snowflake, "id"),
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Snowflake, "snowflake_jwt.alter_user_command"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Snowflake, "id"),
				),
			},
			// ImportState testing
			{ResourceName: testCredentialProviderV2Snowflake, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(testCredentialProviderV2Snowflake, "name", "TF Acceptance Snowflake Token - Modified"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Snowflake, "id"),
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Snowflake, "snowflake_jwt.alter_user_command"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderV2Resource_OAuthClientCredentialsAuthHeader(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/oauth/TestAccCredentialProviderV2Resource.tf")
	modifyFile, _ := os.ReadFile("../../tests/credential/oauth/TestAccCredentialProviderV2Resource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr(testOAuthClientCredentialsV2AuthHeader, "name", "TF Acceptance OAuth"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testOAuthClientCredentialsV2AuthHeader, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testOAuthClientCredentialsV2AuthHeader, "id"),
				),
			},
			// ImportState testing
			{ResourceName: testOAuthClientCredentialsV2AuthHeader, ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(testOAuthClientCredentialsV2AuthHeader, "name", "TF Acceptance OAuth - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderV2Resource_OAuthClientCredentialsPostBody(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/oauth/TestAccCredentialProviderV2Resource.tf")
	modifyFile, _ := os.ReadFile("../../tests/credential/oauth/TestAccCredentialProviderV2Resource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr(testOAuthClientCredentialsV2PostBody, "name", "TF Acceptance OAuth"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testOAuthClientCredentialsV2PostBody, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testOAuthClientCredentialsV2PostBody, "id"),
				),
			},
			// ImportState testing
			{ResourceName: testOAuthClientCredentialsV2PostBody, ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(testOAuthClientCredentialsV2PostBody, "name", "TF Acceptance OAuth - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderV2Resource_UsernamePassword(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/userpass/TestAccCredentialProviderV2Resource.tf")
	modifyFile, _ := os.ReadFile("../../tests/credential/userpass/TestAccCredentialProviderV2Resource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr(testCredentialProviderV2Username, "name", "TF Acceptance Username Password"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Username, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Username, "id"),
				),
			},
			// ImportState testing
			{ResourceName: testCredentialProviderV2Username, ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(testCredentialProviderV2Username, "name", "TF Acceptance Username Password - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderV2Resource_VaultClientToken(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/vault/TestAccCredentialProviderV2Resource.tf")
	modifyFile, _ := os.ReadFile("../../tests/credential/vault/TestAccCredentialProviderV2Resource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr(testCredentialProviderV2Vault, "name", "TF Acceptance Vault"),
					// Verify Tags.
					resource.TestCheckResourceAttr(testCredentialProviderV2Vault, "tags.%", "2"),
					resource.TestCheckResourceAttr(testCredentialProviderV2Vault, "tags.color", "blue"),
					resource.TestCheckResourceAttr(testCredentialProviderV2Vault, "tags.day", "Sunday"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Vault, "id"),
					resource.TestCheckResourceAttr(testCredentialProviderV2Vault, "vault_client_token.vault_forwarding", ""),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testCredentialProviderV2Vault, "id"),
				),
			},
			// ImportState testing
			{ResourceName: testCredentialProviderV2Vault, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(testCredentialProviderV2Vault, "name", "TF Acceptance Vault - Modified"),
					// Verify Tags.
					resource.TestCheckResourceAttr(testCredentialProviderV2Vault, "tags.%", "2"),
					resource.TestCheckResourceAttr(testCredentialProviderV2Vault, "tags.color", "orange"),
					resource.TestCheckResourceAttr(testCredentialProviderV2Vault, "tags.day", "Tuesday"),
					// Verify Vault_Forwarding update
					resource.TestCheckResourceAttr(testCredentialProviderV2Vault, "vault_client_token.vault_forwarding", "conditional"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
