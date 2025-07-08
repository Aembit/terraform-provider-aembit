provider "aembit" {
}

resource "aembit_trust_provider" "oidcidtoken_jwks" {
	name = "TF Acceptance OIDC ID Token JWKS"
	is_active = true
	oidc_id_token = {
		issuer = "issuer"
		subject = "subject"
		audience = "audience"
		jwks = {
			keys = [
				{
					kid = "Tbm3LtlhYlNObRdRc+Tz3mEo2SASPbfR03HI4dmoUkg="
					kty = "RSA"
					use = "sig"
					alg = "RS256"
					e = "AQAB"
				}
			]			
		}
	}
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}