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

resource "aembit_credential_provider" "oidc_id_token_empty_custom_claims" {
	name = "TF Acceptance OIDC ID Token - EmptyClaims"
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
		custom_claims = []
	}
}

resource "aembit_credential_provider" "oidc_id_token_null_custom_claims" {
	name = "TF Acceptance OIDC ID Token - NullClaims"
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