package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testTrustProvidersDataSource string = "data.aembit_trust_providers.test"
	testTrustProviderResource    string = "aembit_trust_provider.kubernetes_key"
)

func testFindTrustProvider(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var err error
		var ok bool
		var notFound bool
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if _, err, notFound = testClient.GetTrustProvider(rs.Primary.ID, nil); notFound {
			return err
		}
		return nil
	}
}

func TestAccTrustProvidersDataSource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/trust/data/TestAccTrustProvidersDataSource.tf")
	createFileConfig, _, _ := randomizeFileConfigs(
		string(createFile),
		"",
		"TF Acceptance Kubernetes",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify non-zero number of Trust Providers returned
					resource.TestCheckResourceAttrSet(
						testTrustProvidersDataSource,
						"trust_providers.#",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(
						testTrustProvidersDataSource,
						"trust_providers.0.id",
					),
					// Find newly created entry
					testFindTrustProvider(testTrustProviderResource),
				),
			},
		},
	})
}
