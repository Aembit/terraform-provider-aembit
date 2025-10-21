provider "aembit" {
}

data "aembit_roles" "test" {}

resource "aembit_identity_provider" "test_idp_oidc" {
	name = "Identity Provider OIDC for TF Acceptance Test"
	description = "Description of Identity Provider for TF Acceptance Test"
	is_active = true
    oidc = {
        oidc_base_url = "https://test.oidc.com"
        client_id = "test_client_id"
        scopes = "profile email"
        auth_type = "ClientSecret"
        pcke_required = true
    }
}