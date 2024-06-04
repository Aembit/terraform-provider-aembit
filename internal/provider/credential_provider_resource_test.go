package provider

import (
	"fmt"
	"os"
	"testing"

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
		if ok, err = testClient.DeleteCredentialProvider(rs.Primary.ID, nil); !ok {
			return err
		}
		return nil
	}
}

func TestAccCredentialProviderResource_AembitToken(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/aembit/TestAccCredentialProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/credential/aembit/TestAccCredentialProviderResource.tfmod")
	createFileConfig, modifyFileConfig, _ := randomizeFileConfigs(string(createFile), string(modifyFile), "IGNORE")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider set values
					resource.TestCheckResourceAttr(testCredentialProviderAembit, "name", "TF Acceptance Aembit Token"),
					resource.TestCheckResourceAttr(testCredentialProviderAembit, "aembit_access_token.lifetime", "1800"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderAembit, "id"),
					resource.TestCheckResourceAttrSet(testCredentialProviderAembit, "aembit_access_token.audience"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testCredentialProviderAembit, "id"),
				),
			},
			// Test Aembit API Removal causes re-create with non-empty plan
			{Config: createFileConfig, Check: testDeleteCredentialProvider(testCredentialProviderAembit), ExpectNonEmptyPlan: true},
			// Recreate the resource from the first test step
			{Config: createFileConfig},
			// ImportState testing
			{ResourceName: testCredentialProviderAembit, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(testCredentialProviderAembit, "name", "TF Acceptance Aembit Token - Modified"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderAembit, "aembit_access_token.audience"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderResource_ApiKey(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/apikey/TestAccCredentialProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/credential/apikey/TestAccCredentialProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
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
			{ResourceName: "aembit_credential_provider.api_key", ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_credential_provider.api_key", "name", "TF Acceptance API Key - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderResource_AwsSTS(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/aws/TestAccCredentialProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/credential/aws/TestAccCredentialProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential set values
					resource.TestCheckResourceAttr("aembit_credential_provider.aws", "name", "TF Acceptance AWS STS"),
					resource.TestCheckResourceAttr("aembit_credential_provider.aws", "aws_sts.lifetime", "1800"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.aws", "id"),
					resource.TestCheckResourceAttrSet("aembit_credential_provider.aws", "aws_sts.oidc_issuer"),
					resource.TestCheckResourceAttrSet("aembit_credential_provider.aws", "aws_sts.token_audience"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_credential_provider.aws", "id"),
				),
			},
			// ImportState testing
			{ResourceName: "aembit_credential_provider.aws", ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_credential_provider.aws", "name", "TF Acceptance AWS STS - Modified"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.aws", "id"),
					resource.TestCheckResourceAttrSet("aembit_credential_provider.aws", "aws_sts.oidc_issuer"),
					resource.TestCheckResourceAttrSet("aembit_credential_provider.aws", "aws_sts.token_audience"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderResource_GoogleWorkload(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/gcp/TestAccCredentialProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/credential/gcp/TestAccCredentialProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider set values
					resource.TestCheckResourceAttr("aembit_credential_provider.gcp", "name", "TF Acceptance GCP Workload"),
					resource.TestCheckResourceAttr("aembit_credential_provider.gcp", "google_workload_identity.lifetime", "1800"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.gcp", "id"),
					resource.TestCheckResourceAttrSet("aembit_credential_provider.gcp", "google_workload_identity.oidc_issuer"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_credential_provider.gcp", "id"),
				),
			},
			// ImportState testing
			{ResourceName: "aembit_credential_provider.gcp", ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_credential_provider.gcp", "name", "TF Acceptance GCP Workload - Modified"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.gcp", "id"),
					resource.TestCheckResourceAttrSet("aembit_credential_provider.gcp", "google_workload_identity.oidc_issuer"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderResource_SnowflakeToken(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/snowflake/TestAccCredentialProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/credential/snowflake/TestAccCredentialProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr("aembit_credential_provider.snowflake", "name", "TF Acceptance Snowflake Token"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.snowflake", "id"),
					resource.TestCheckResourceAttrSet("aembit_credential_provider.snowflake", "snowflake_jwt.alter_user_command"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_credential_provider.snowflake", "id"),
				),
			},
			// ImportState testing
			{ResourceName: "aembit_credential_provider.snowflake", ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_credential_provider.snowflake", "name", "TF Acceptance Snowflake Token - Modified"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.snowflake", "id"),
					resource.TestCheckResourceAttrSet("aembit_credential_provider.snowflake", "snowflake_jwt.alter_user_command"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderResource_OAuthClientCredentialsAuthHeader(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/oauth/TestAccCredentialProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/credential/oauth/TestAccCredentialProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr("aembit_credential_provider.oauth_authHeader", "name", "TF Acceptance OAuth"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.oauth_authHeader", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_credential_provider.oauth_authHeader", "id"),
				),
			},
			// ImportState testing
			{ResourceName: "aembit_credential_provider.oauth_authHeader", ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_credential_provider.oauth_authHeader", "name", "TF Acceptance OAuth - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderResource_OAuthClientCredentialsPostBody(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/oauth/TestAccCredentialProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/credential/oauth/TestAccCredentialProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr("aembit_credential_provider.oauth_postBody", "name", "TF Acceptance OAuth"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.oauth_postBody", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_credential_provider.oauth_postBody", "id"),
				),
			},
			// ImportState testing
			{ResourceName: "aembit_credential_provider.oauth_postBody", ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_credential_provider.oauth_postBody", "name", "TF Acceptance OAuth - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderResource_UsernamePassword(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/userpass/TestAccCredentialProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/credential/userpass/TestAccCredentialProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr("aembit_credential_provider.userpass", "name", "TF Acceptance Username Password"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.userpass", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_credential_provider.userpass", "id"),
				),
			},
			// ImportState testing
			{ResourceName: "aembit_credential_provider.userpass", ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_credential_provider.userpass", "name", "TF Acceptance Username Password - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderResource_VaultClientToken(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/vault/TestAccCredentialProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/credential/vault/TestAccCredentialProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr("aembit_credential_provider.vault", "name", "TF Acceptance Vault"),
					// Verify Tags.
					resource.TestCheckResourceAttr("aembit_credential_provider.vault", "tags.%", "2"),
					resource.TestCheckResourceAttr("aembit_credential_provider.vault", "tags.color", "blue"),
					resource.TestCheckResourceAttr("aembit_credential_provider.vault", "tags.day", "Sunday"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.vault", "id"),
					resource.TestCheckResourceAttr("aembit_credential_provider.vault", "vault_client_token.vault_forwarding", ""),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_credential_provider.vault", "id"),
				),
			},
			// ImportState testing
			{ResourceName: "aembit_credential_provider.vault", ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_credential_provider.vault", "name", "TF Acceptance Vault - Modified"),
					// Verify Tags.
					resource.TestCheckResourceAttr("aembit_credential_provider.vault", "tags.%", "2"),
					resource.TestCheckResourceAttr("aembit_credential_provider.vault", "tags.color", "orange"),
					resource.TestCheckResourceAttr("aembit_credential_provider.vault", "tags.day", "Tuesday"),
					// Verify Vault_Forwarding update
					resource.TestCheckResourceAttr("aembit_credential_provider.vault", "vault_client_token.vault_forwarding", "conditional"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
