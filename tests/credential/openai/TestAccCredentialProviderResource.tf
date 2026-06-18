provider "aembit" {
}

resource "aembit_credential_provider" "openai" {
	name      = "TF Acceptance OpenAI Wif"
	is_active = true
	openai_wif = {
		identity_provider_id = "idp_test"
		service_account_id   = "svac_test"
		audience             = "aud_test"
	}
}
