provider "aembit" {
}

resource "aembit_credential_provider" "oauth_authHeader" {
	name = "TF Acceptance OAuth"
	is_active = true
	tags = {
		color = "blue"
        day   = "Sunday"
    }
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
				value_type = "literal"
			}
		]
	}
}

resource "aembit_credential_provider" "oauth_postBody" {
	name = "TF Acceptance OAuth"
	is_active = true
	tags = {
		color = "blue"
        day   = "Sunday"
    }
	oauth_client_credentials = {
		token_url = "https://aembit.io/token"
		client_id = "test_client_id"
		client_secret = "test_client_secret"
		scopes = "test_scopes"
		credential_style = "postBody"
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
