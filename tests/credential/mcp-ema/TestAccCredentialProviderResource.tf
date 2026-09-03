provider "aembit" {
}

data "aembit_roles" "test" {}

locals {
  role_ids_by_name = { for role in data.aembit_roles.test.roles : role.name => role.id }
}

resource "aembit_identity_provider" "test_idp_oidc" {
	name        = "Identity Provider OIDC for TF Acceptance Test MCP EMA"
	description = "Description of Identity Provider for TF Acceptance Test MCP EMA"
	is_active   = true
	sso_statement_role_mappings = [
		{
			attribute_name  = "test-attribute-name"
			attribute_value = "test-attribute-value"
			roles           = [local.role_ids_by_name["SuperAdmin"], local.role_ids_by_name["Auditor"]]
		}
	]
	oidc = {
		oidc_base_url = "https://test.oidc.com"
		client_id     = "test_client_id"
		scopes        = "profile email"
		auth_type     = "ClientSecret"
		client_secret = "some_secret"
		pcke_required = true
	}
}

resource "aembit_credential_provider" "mcp_ema" {
	name      = "TF Acceptance MCP EMA"
	is_active = true
	mcp_ema = {
		issuer               = "https://idp.example.com"
		mcp_server_url       = "https://mcp.example.com"
		client_id            = "test-client-id"
		scopes               = "openid profile email"
		authorization_url   = "https://idp.example.com/oauth/authorize"
		token_url           = "https://idp.example.com/oauth/token"
		introspection_url   = "https://idp.example.com/oauth/introspect"
		is_corporate_idp     = true
		identity_provider_id = aembit_identity_provider.test_idp_oidc.id
	}
}
