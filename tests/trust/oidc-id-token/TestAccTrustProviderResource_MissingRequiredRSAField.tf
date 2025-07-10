provider "aembit" {
}

resource "aembit_trust_provider" "oidcidtoken_jwks" {
	name = "TF Acceptance OIDC ID Token JWKS"
	is_active = true
	oidc_id_token = {
		issuer = "issuer"
		subject = "subject"
		audience = "audience"
		jwks = <<-EOT
{
  "keys": [
    {
      "kty": "RSA",
      "use": "sig",
      "kid": "T41hVcPtA3ehDjSaZXSI9LKuanyTkBOf0YKlAM6gtNQ=",
      "e": "AQAB",
      "alg": "RS256"
    },    
  ]
}		
EOT		
	}
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}