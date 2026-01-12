provider "aembit" {
}

resource "aembit_credential_provider" "oauth_authorization_code" {
	id = "replace-with-uuid-first"
	name = "TF Acceptance OAuth Authorization Code"
	is_active = true
	oauth_authorization_code = {
		oauth_discovery_url = "https://aembit.io/.well-known/openid-configuration"
		oauth_authorization_url = "https://aembit.io/authorize"
		oauth_token_url = "https://aembit.io/token"
		client_id = "test_client_id"
		client_secret = "test_client_secret"
		scopes = "test_scopes"
		is_pkce_required = true
		lifetime = 31536000
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

resource "aembit_credential_provider" "oauth_authorization_code_empty_custom_parameters" {
	id = "replace-with-uuid-second"
	name = "TF Acceptance OAuth Authorization Code"
	is_active = true
	oauth_authorization_code = {
		oauth_discovery_url = "https://aembit.io/.well-known/openid-configuration"
		oauth_authorization_url = "https://aembit.io/authorize"
		oauth_token_url = "https://aembit.io/token"
		oauth_introspection_url = "https://aembit.io/introspect"
		client_id = "test_client_id"
		client_secret = "test_client_secret"
		scopes = "test_scopes"
		is_pkce_required = true
		lifetime = 31536000
	}
}

