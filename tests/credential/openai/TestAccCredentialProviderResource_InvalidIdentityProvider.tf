provider "aembit" {
}

resource "aembit_credential_provider" "openai" {
	name      = "TF Acceptance OpenAI Wif Invalid IDP"
	is_active = true
	openai_wif = {
		identity_provider_id = "invalid_idp"
		service_account_id   = "user-test"
		audience             = "aud_test"
	}
}
