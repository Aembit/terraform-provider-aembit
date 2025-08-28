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
		vault_host = "vault.aembit.io"
		vault_port = 8200
		vault_tls = true
		vault_namespace = "vault_namespace"
		vault_path = "vault_path"
		vault_role = "vault_role"
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

resource "aembit_credential_provider" "vault_empty_custom_claims" {
	name = "TF Acceptance Vault - EmptyClaims"
	is_active = true
    tags = {
        color = "blue"
        day   = "Sunday"
    }
	vault_client_token = {
		subject = "subject"
		subject_type = "literal"
		lifetime = 60
		vault_host = "vault.aembit.io"
		vault_port = 8200
		vault_tls = true
		vault_namespace = "vault_namespace"
		vault_path = "vault_path"
		vault_role = "vault_role"
		custom_claims = []
	}
}

resource "aembit_credential_provider" "vault_null_custom_claims" {
	name = "TF Acceptance Vault - NullClaims"
	is_active = true
    tags = {
        color = "blue"
        day   = "Sunday"
    }
	vault_client_token = {
		subject = "subject"
		subject_type = "literal"
		lifetime = 60
		vault_host = "vault.aembit.io"
		vault_port = 8200
		vault_tls = true
		vault_namespace = "vault_namespace"
		vault_path = "vault_path"
		vault_role = "vault_role"
	}
}

