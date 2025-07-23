package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testCredentialProviderIntegrationGitLab string = "aembit_credential_provider_integration.gitlab"

// resource id set in the test file
const testCredentialProviderIntegrationAwsIamRole string = "aembit_credential_provider_integration.awsiamrole_cpi"

func testDeleteCredentialProviderIntegration(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var ok bool
		var err error
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if ok, err = testClient.DeleteCredentialProviderIntegration(context.Background(), rs.Primary.ID, nil); !ok {
			return err
		}
		return nil
	}
}

func TestAccGitLabCPIResource(t *testing.T) {
	t.Skip("skipping test until we figure out a way to handle the GitLab tokens appropriately")

	createFile, _ := os.ReadFile(
		"../../tests/credential_provider_integration/gitlab/TestAccCredentialProviderIntegrationResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/credential_provider_integration/gitlab/TestAccCredentialProviderIntegrationResource.tfmod",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Integration Name
					resource.TestCheckResourceAttr(
						testCredentialProviderIntegrationGitLab,
						"name",
						"TF Acceptance GitLab Credential Integration",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(
						testCredentialProviderIntegrationGitLab,
						"id",
					),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(
						testCredentialProviderIntegrationGitLab,
						"id",
					),
				),
			},
			// Test Aembit API Removal causes re-create with non-empty plan
			{
				Config: string(createFile),
				Check: testDeleteCredentialProviderIntegration(
					testCredentialProviderIntegrationGitLab,
				),
				ExpectNonEmptyPlan: true,
			},
			// Recreate the resource from the first test step
			{Config: string(createFile)},
			// ImportState testing
			{
				ResourceName:      testCredentialProviderIntegrationGitLab,
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						testCredentialProviderIntegrationGitLab,
						"name",
						"TF Acceptance GitLab Credential Integration - Modified",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccAwsIamRoleCPIResource(t *testing.T) {
	createFile, _ := os.ReadFile(
		"../../tests/credential_provider_integration/awsiamrole/TestAccAwsIamRoleCpiResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/credential_provider_integration/awsiamrole/TestAccAwsIamRoleCpiResource.tfmod",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Integration Name
					resource.TestCheckResourceAttr(
						testCredentialProviderIntegrationAwsIamRole,
						"name",
						"TF Acceptance Aws IAM Role Credential Provider Integration",
					),
					resource.TestCheckResourceAttr(
						testCredentialProviderIntegrationAwsIamRole,
						"description",
						"TF Acceptance Aws IAM Role Credential Provider Integration Description",
					),
					resource.TestCheckResourceAttr(
						testCredentialProviderIntegrationAwsIamRole,
						"aws_iam_role.role_arn",
						"arn:aws:iam::123456789012:role/MyRole",
					),
					resource.TestCheckResourceAttr(
						testCredentialProviderIntegrationAwsIamRole,
						"aws_iam_role.lifetime_in_seconds",
						"3600",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(
						testCredentialProviderIntegrationAwsIamRole,
						"id",
					),
				),
			},
			// Test Aembit API Removal causes re-create with non-empty plan
			{
				Config: string(createFile),
				Check: testDeleteCredentialProviderIntegration(
					testCredentialProviderIntegrationAwsIamRole,
				),
				ExpectNonEmptyPlan: true,
			},
			// Recreate the resource from the first test step
			{Config: string(createFile)},
			// ImportState testing
			{
				ResourceName:      testCredentialProviderIntegrationAwsIamRole,
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						testCredentialProviderIntegrationAwsIamRole,
						"name",
						"TF Acceptance Aws IAM Role Credential Provider Integration - Updated",
					),
					resource.TestCheckResourceAttr(
						testCredentialProviderIntegrationAwsIamRole,
						"description",
						"TF Acceptance Aws IAM Role Credential Provider Integration Description - Updated",
					),
					resource.TestCheckResourceAttr(
						testCredentialProviderIntegrationAwsIamRole,
						"aws_iam_role.role_arn",
						"arn:aws:iam::123456789012:role/MyRoleUpdated",
					),
					resource.TestCheckResourceAttr(
						testCredentialProviderIntegrationAwsIamRole,
						"aws_iam_role.lifetime_in_seconds",
						"5000",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
