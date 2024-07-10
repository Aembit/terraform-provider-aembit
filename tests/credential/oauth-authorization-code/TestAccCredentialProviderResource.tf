provider "aembit" {
}

resource "aembit_credential_provider" "oauth_authorization_code" {
	id = "a8f8ecdc-eb48-4a29-b854-d25745628b51"
	name = "TF Acceptance OAuth Authorization Code"
	is_active = true
	tags = {
		color = "blue"
        day   = "Sunday"
    }
	oauth_authorization_code = {
		oauth_url = "https://aembit.io/.well-known/openid-configuration"
		authorization_url = "https://aembit.io/authorize"
		token_url = "https://aembit.io/token"
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
				value_type = "dynamic"
			}
		]
	}
}
