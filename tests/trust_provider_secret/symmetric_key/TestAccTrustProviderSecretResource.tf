provider "aembit" {
}

resource "aembit_trust_provider" "oidcidtoken_symmetric_key" {
	name = "TF Acceptance OIDC ID Token Symmetric Key"
	is_active = true
	oidc_id_token = {
		issuer = "issuer"
		subject = "subject"
		audience = "audience"
		oidc_endpoint = "https://3a3b5d.id.devbroadangle.aembit-eng.com/"
	}
}

resource "aembit_trust_provider_secret" "symmetric_key_secret1" {
	trust_provider_id = aembit_trust_provider.oidcidtoken_symmetric_key.id
	secret      = "c2VjcmV0MQ=="
	type = "SymmetricKey"
}

resource "aembit_trust_provider_secret" "symmetric_key_secret2" {
	trust_provider_id = aembit_trust_provider.oidcidtoken_symmetric_key.id
	secret      = "c2VjcmV0Mg=="
	type = "SymmetricKey"
}