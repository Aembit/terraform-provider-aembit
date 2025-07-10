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
      "kty": "EC",
      "use": "sig",
      "kid": "iTYG7jb2cyaQ04cp69CpoMBzxNjRmixlGGxZTIHSpXg=",
      "alg": "ES256",
      "x": "ItvdSxTnkMqPq3kKeHYlAF1ArZqz4_CXjUmiPvHDQ08",
      "y": "C2z0b9zNhvywzboDt03F2xb_7fOaw8LWbakgudjN3kE",
    }
  ]
}
EOT	
	}
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}