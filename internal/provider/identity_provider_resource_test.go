package provider

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccIdentityProviderResource_Saml verifies the end-to-end lifecycle of SAML identity providers.
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

// TestAccIdentityProviderResource_Oidc verifies create, read, import, and update lifecycle of OIDC identity providers.
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
					resource.TestCheckResourceAttr(
						resourceName,
						"oidc.prompt",
						"login",
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
					resource.TestCheckResourceAttr(
						resourceName,
						"oidc.prompt",
						"login consent",
					),
				),
			},
		},
	})
}

// TestAccIdentityProviderResource_Oidc_MissingSecret verifies plan rejection when auth_type is ClientSecret without a secret.
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

// TestAccIdentityProviderResource_Oidc_InvalidTtlLow verifies plan-time rejection when group_cache_ttl_minutes is below range.
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

// TestAccIdentityProviderResource_Oidc_InvalidTtlHigh verifies plan-time rejection when group_cache_ttl_minutes is above range.
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

// TestAccIdentityProviderResource_Oidc_InvalidPromptToken verifies plan-time rejection of an unrecognized prompt token.
func TestAccIdentityProviderResource_Oidc_InvalidPromptToken(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "aembit_identity_provider" "test_idp_oidc" {
	name        = "Identity Provider OIDC Invalid Prompt"
	description = "Testing invalid prompt token"
	is_active   = true
	oidc = {
		oidc_base_url = "https://test.oidc.com"
		client_id     = "test_client_id"
		scopes        = "profile email"
		auth_type     = "ClientSecret"
		client_secret = "some_secret"
		prompt        = "invalid_prompt"
	}
}
`,
				ExpectError: regexp.MustCompile(`Invalid OIDC prompt value "invalid_prompt"`),
			},
		},
	})
}

// TestAccIdentityProviderResource_Oidc_DuplicatePromptToken verifies plan-time rejection of duplicate prompt tokens.
func TestAccIdentityProviderResource_Oidc_DuplicatePromptToken(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "aembit_identity_provider" "test_idp_oidc" {
	name        = "Identity Provider OIDC Duplicate Prompt"
	description = "Testing duplicate prompt token"
	is_active   = true
	oidc = {
		oidc_base_url = "https://test.oidc.com"
		client_id     = "test_client_id"
		scopes        = "profile email"
		auth_type     = "ClientSecret"
		client_secret = "some_secret"
		prompt        = "login login"
	}
}
`,
				ExpectError: regexp.MustCompile(`Duplicate OIDC prompt value "login" found`),
			},
		},
	})
}

// TestAccIdentityProviderResource_Oidc_MutuallyExclusiveNonePrompt verifies plan-time rejection of combining 'none' with other prompt tokens.
func TestAccIdentityProviderResource_Oidc_MutuallyExclusiveNonePrompt(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "aembit_identity_provider" "test_idp_oidc" {
	name        = "Identity Provider OIDC Mutually Exclusive None Prompt"
	description = "Testing mutually exclusive none prompt token"
	is_active   = true
	oidc = {
		oidc_base_url = "https://test.oidc.com"
		client_id     = "test_client_id"
		scopes        = "profile email"
		auth_type     = "ClientSecret"
		client_secret = "some_secret"
		prompt        = "none login"
	}
}
`,
				ExpectError: regexp.MustCompile(`OIDC prompt 'none' cannot be combined with any other prompt values`),
			},
		},
	})
}
