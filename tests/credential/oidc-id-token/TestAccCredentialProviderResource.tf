provider "aembit" {
}

resource "aembit_credential_provider" "oidc_id_token" {
	name = "TF Acceptance OIDC ID Token"
	is_active = true
    tags = {
        color = "blue"
        day   = "Sunday"
    }
	oidc_id_token = {
		subject = "subject"
		subject_type = "literal"
		lifetime_in_minutes = 60
		audience = "test.aembit.io"
		algorithm_type = "ES256"
	}
}