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
		custom_claims = [
			{
				key = "aud"
				value = "audience"
				value_type = "literal"
			},
			{
				key = "sub"
				value = "subject"
				value_type = "dynamic"
			},
			{
				key = "other"
				value = "test"
				value_type = "literal"
			}
		]
	}
}