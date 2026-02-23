provider "aembit" {
}

resource "aembit_credential_provider" "mcp_user_based_access_token" {
	id = "replace-with-uuid-first"
	name = "TF Acceptance McpUserBasedAccessToken"
	is_active = true
	mcp_user_based_access_token = {
		mcp_server_url = "https://aembit.io/.well-known/openid-configuration"
		oauth_authorization_url = "https://aembit.io/authorize"
		oauth_token_url = "https://aembit.io/token"
		client_id = "test_client_id"
		client_secret = "test_client_secret"
		scopes = "test_scopes"
		is_pkce_required = true
		custom_parameters = [
			{
				key = "key"
				value = "value"
				value_type = "literal"
			},
			{
				key = "key2"
				value = "value2"
				value_type = "literal"
			}
		]
	}
}

resource "aembit_credential_provider" "mcp_user_based_access_token_empty_custom_parameters" {
	id = "replace-with-uuid-second"
	name = "TF Acceptance McpUserBasedAccessToken"
	is_active = true
	mcp_user_based_access_token = {
		mcp_server_url = "https://aembit.io/.well-known/openid-configuration"
		oauth_authorization_url = "https://aembit.io/authorize"
		oauth_token_url = "https://aembit.io/token"
		oauth_introspection_url = "https://aembit.io/introspect"
		client_id = "test_client_id"
		client_secret = "test_client_secret"
		scopes = "test_scopes"
		resource = "resourceid"
		is_pkce_required = true
	}
}

