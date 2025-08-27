provider "aembit" {
}

resource "aembit_credential_provider" "jwt_svid_token" {
	name = "TF Acceptance JWT-SVID Token"
	is_active = true
    tags = {
        color = "blue"
        day   = "Sunday"
    }
	jwt_svid_token = {
		subject = "spiffe://test.com"
		subject_type = "literal"
		lifetime_in_minutes = 60
		audience = "test.aembit.io"
		algorithm_type = "ES256"		
	}
}
