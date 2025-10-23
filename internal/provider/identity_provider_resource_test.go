package provider

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIdentityProviderResource_Saml(t *testing.T) {
	createFile, err := os.ReadFile(
		"../../tests/identity_provider/saml/TestAccIdentityProviderResource.tf",
	)
	if err != nil {
		t.Fatalf("failed to read test config file: %v", err)
	}
	modifiedFile, err := os.ReadFile(
		"../../tests/identity_provider/saml/TestAccIdentityProviderResource.tfmod",
	)
	if err != nil {
		t.Fatalf("failed to read test config file: %v", err)
	}
	resourceName := "aembit_identity_provider.test_idp_saml"
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
						"Identity Provider SAML for TF Acceptance Test",
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
					// Verify Tags.
					resource.TestCheckResourceAttr(resourceName, tagsCount, "2"),
					resource.TestCheckResourceAttr(
						resourceName,
						tagsAllCount,
						"4",
					),
					resource.TestCheckResourceAttr(resourceName, tagsColor, "blue"),
					resource.TestCheckResourceAttr(resourceName, tagsDay, "Sunday"),
					resource.TestCheckResourceAttr(
						resourceName,
						tagsAllName,
						"Terraform",
					),
					resource.TestCheckResourceAttr(
						resourceName,
						tagsAllOwner,
						"Aembit",
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
						"Identity Provider SAML for TF Acceptance Test - Updated",
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

func TestAccIdentityProviderResource_Oidc(t *testing.T) {
	createFile, err := os.ReadFile(
		"../../tests/identity_provider/oidc/TestAccIdentityProviderResource.tf",
	)
	if err != nil {
		t.Fatalf("failed to read test config file: %v", err)
	}
	modifiedFile, err := os.ReadFile(
		"../../tests/identity_provider/oidc/TestAccIdentityProviderResource.tfmod",
	)
	if err != nil {
		t.Fatalf("failed to read test config file: %v", err)
	}
	resourceName := "aembit_identity_provider.test_idp_oidc"
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
						"Identity Provider OIDC for TF Acceptance Test",
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
						"Identity Provider OIDC for TF Acceptance Test - Updated",
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

func TestAccIdentityProviderResource_Oidc_MissingSecret(t *testing.T) {
	createFile, _ := os.ReadFile(
		"../../tests/identity_provider/oidc/TestAccIdentityProviderResource_MissingSecret.tf",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      string(createFile),
				ExpectError: regexp.MustCompile(`attribute must be set.`),
			},
		},
	})
}
