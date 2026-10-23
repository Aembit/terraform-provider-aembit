package provider

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIdentityProviderResource_Saml(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
					resource.TestCheckResourceAttr(
						resourceName,
						"oidc.enable_userinfo_group_lookup",
						"true",
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"oidc.group_claim_attribute",
						"roles",
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"oidc.group_cache_ttl_minutes",
						"30",
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"oidc.token_url",
						"https://test.oidc.com/oauth/v2/token",
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
					resource.TestCheckResourceAttr(
						resourceName,
						"oidc.enable_userinfo_group_lookup",
						"false",
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"oidc.group_claim_attribute",
						"groups",
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"oidc.group_cache_ttl_minutes",
						"15",
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"oidc.token_url",
						"https://test.oidc.com/oauth/v2/token-updated",
					),
				),
			},
		},
	})
}

func TestAccIdentityProviderResource_Oidc_MissingSecret(t *testing.T) {
	t.Parallel()
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

func TestAccIdentityProviderResource_Oidc_InvalidTtlLow(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile(
		"../../tests/identity_provider/oidc/TestAccIdentityProviderResource_InvalidTtlLow.tf",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      string(createFile),
				ExpectError: regexp.MustCompile(`value must be between 5 and 1440`),
			},
		},
	})
}

func TestAccIdentityProviderResource_Oidc_InvalidTtlHigh(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile(
		"../../tests/identity_provider/oidc/TestAccIdentityProviderResource_InvalidTtlHigh.tf",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      string(createFile),
				ExpectError: regexp.MustCompile(`value must be between 5 and 1440`),
			},
		},
	})
}
