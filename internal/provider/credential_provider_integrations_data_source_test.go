package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testCredentialProviderIntegrationsDataSource string = "data.aembit_credential_provider_integrations.test"
	testCredentialProviderIntegrationResource    string = "aembit_credential_provider_integration.gitlab"
)

func testFindCredentialProviderIntegration(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var err error
		var ok bool
		var notFound bool
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if _, err, notFound = testClient.GetCredentialProviderIntegration(rs.Primary.ID, nil); notFound {
			return err
		}
		return nil
	}
}

func TestAccCredentialProviderIntegrationsDataSource(t *testing.T) {
	t.Skip("skipping test until we figure out a way to handle the GitLab tokens appropriately")

	createFile, _ := os.ReadFile(
		"../../tests/credential_provider_integration/data/TestAccCredentialProviderIntegrationsDataSource.tf",
	)
	createFileConfig, _, _ := randomizeFileConfigs(
		string(createFile),
		"",
		"TF Acceptance GitLab Credential Integration",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify non-zero number of Integrations returned
					resource.TestCheckResourceAttrSet(
						testCredentialProviderIntegrationsDataSource,
						"credential_provider_integrations.#",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(
						testCredentialProviderIntegrationsDataSource,
						"credential_provider_integrations.0.id",
					),
					// Find newly created entry
					testFindCredentialProviderIntegration(
						testCredentialProviderIntegrationResource,
					),
				),
			},
		},
	})
}
