provider "aembit" {
}

resource "aembit_credential_provider_v2" "userpass" {
	name = "TF Acceptance Username Password"
	is_active = true
	username_password = {
		username = "username"
		password = "password"
	}
}
