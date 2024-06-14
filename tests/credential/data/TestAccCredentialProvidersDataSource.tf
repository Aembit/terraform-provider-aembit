provider "aembit" {
}

resource "aembit_credential_provider" "oauth" {
	name = "TF Acceptance OAuth"
	is_active = true
	oauth_client_credentials = {
		token_url = "https://aembit.io/token"
		client_id = "test_client_id"
		client_secret = "test_client_secret"
		scopes = "test_scopes"
		credential_style = "authHeader"
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

data "aembit_credential_providers" "test" {
    depends_on = [ aembit_credential_provider.oauth ]
}

output "test" {
  value     = data.aembit_credential_providers.test
  sensitive = true
}