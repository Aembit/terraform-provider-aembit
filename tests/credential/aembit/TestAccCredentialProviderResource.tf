provider "aembit" {
}

resource "aembit_credential_provider" "aembit" {
	name = "TF Acceptance Aembit Token"
	is_active = true
	aembit_access_token = {
		role_id = "87e3c07e-afa2-49cf-b585-769c71557d02"
		lifetime = 900
	}
}
