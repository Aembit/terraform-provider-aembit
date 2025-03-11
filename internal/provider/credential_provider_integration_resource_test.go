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

func TestAccCredentialProviderIntegrationResource_GitLab(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential_provider_integration/gitlab/TestAccCredentialProviderIntegrationResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/credential_provider_integration/gitlab/TestAccCredentialProviderIntegrationResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Integration Name
					resource.TestCheckResourceAttr(testCredentialProviderIntegrationGitLab, "name", "TF Acceptance GitLab Credential Integration"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCredentialProviderIntegrationGitLab, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testCredentialProviderIntegrationGitLab, "id"),
				),
			},
			// Test Aembit API Removal causes re-create with non-empty plan
			{Config: string(createFile), Check: testDeleteCredentialProviderIntegration(testCredentialProviderIntegrationGitLab), ExpectNonEmptyPlan: true},
			// Recreate the resource from the first test step
			{Config: string(createFile)},
			// ImportState testing
			{ResourceName: testCredentialProviderIntegrationGitLab, ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(testCredentialProviderIntegrationGitLab, "name", "TF Acceptance GitLab Credential Integration - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
