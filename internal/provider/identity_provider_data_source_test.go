package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testIdpDataSource = "data.aembit_identity_providers.test_idps"

func TestAccIdentityProviderDataSource(t *testing.T) {
	t.Parallel()
	createFile, err := os.ReadFile(
		"../../tests/identity_provider/data/TestAccIdentityProviderDataSource.tf",
	)
	if err != nil {
		t.Fatalf("failed to read test config file: %v", err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testIdpDataSource, "identity_providers.#"),
					resource.TestCheckResourceAttrSet(testIdpDataSource, "identity_providers.0.id"),
					resource.TestCheckResourceAttr(
						testIdpDataSource,
						"identity_providers.0.name",
						"Identity Provider for TF Acceptance Test",
					),
				),
			},
		},
	})
}
