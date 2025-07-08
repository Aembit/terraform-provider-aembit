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
					kid = "mxAhc1VybhA8LT2jHRQFEzcWSoLbFmnDWGYoViS/aKw="
					kty = "EC"
					use = "sig"
					alg = "ES256"
					x: "x-pRNOyN2BwmgvPuLTOEJMLB1vcc4vljjU41W0jz5Sw",
					crv: "P-256"				
				}
			]		
		}
	}
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}