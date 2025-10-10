package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIdentityProviderResource(t *testing.T) {
	createFile, err := os.ReadFile(
		"../../tests/identity_provider/TestAccIdentityProviderResource.tf",
	)
	if err != nil {
		t.Fatalf("failed to read test config file: %v", err)
	}
	modifiedFile, err := os.ReadFile(
		"../../tests/identity_provider/TestAccIdentityProviderResource.tfmod",
	)
	if err != nil {
		t.Fatalf("failed to read test config file: %v", err)
	}
	resourceName := "aembit_identity_provider.test_idp"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: string(createFile),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName,
						"name",
						"Identity Provider for TF Acceptance Test",
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"description",
						"Description of Identity Provider for TF Acceptance Test",
					),
					resource.TestCheckResourceAttr(resourceName, "is_active", "true"),
					resource.TestCheckResourceAttr(
						resourceName,
						"sso_statement_role_mappings.#",
						"1",
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"sso_statement_role_mappings.0.attribute_name",
						"test-attribute-name",
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"sso_statement_role_mappings.0.attribute_value",
						"test-attribute-value",
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"sso_statement_role_mappings.0.roles.#",
						"2",
					),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
			},
			{
				Config: string(modifiedFile),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName,
						"name",
						"Identity Provider for TF Acceptance Test - Updated",
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"description",
						"Description of Identity Provider for TF Acceptance Test - Updated",
					),
					resource.TestCheckResourceAttr(resourceName, "is_active", "true"),
					resource.TestCheckResourceAttr(
						resourceName,
						"sso_statement_role_mappings.#",
						"1",
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"sso_statement_role_mappings.0.attribute_name",
						"test-attribute-name-updated",
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"sso_statement_role_mappings.0.attribute_value",
						"test-attribute-value-updated",
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"sso_statement_role_mappings.0.roles.#",
						"1",
					),
				),
			},
		},
	})
}
