provider "aembit" {
}

resource "aembit_credential_provider" "vault" {
	name = "TF Acceptance Vault"
	is_active = true
    tags = {
        color = "blue"
        day   = "Sunday"
    }
	vault_client_token = {
		subject = "subject"
		subject_type = "literal"
		lifetime = 60
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
		vault_host = "vault.aembit.io"
		vault_port = 8200
		vault_tls = true
		vault_namespace = "vault_namespace"
		vault_path = "vault_path"
		vault_role = "vault_role"
	}
}