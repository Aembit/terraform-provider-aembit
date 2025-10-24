provider "aembit" {
}

resource "aembit_credential_provider" "jwt_svid_token" {
	name = "TF Acceptance JWT-SVID Token"
	is_active = true
	jwt_svid_token = {
		subject = "test.com"
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
