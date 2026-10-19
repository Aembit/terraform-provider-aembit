provider "aembit" {
}

data "aembit_roles" "test" {}

resource "aembit_identity_provider" "test_idp_oidc_invalid_ttl_high" {
	name = "Identity Provider OIDC High TTL"
	description = "Description"
	is_active = true
	oidc = {
		oidc_base_url = "https://test.oidc.com"
		client_id = "test_client_id"
		scopes = "profile email"
		auth_type = "ClientSecret"
		client_secret = "some_secret"
		pcke_required = true
		group_cache_ttl_minutes = 2000
	}
}
