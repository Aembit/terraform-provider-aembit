provider "aembit" {
}

data "aembit_roles" "test" {}

locals {
  role_ids_by_name = { for role in data.aembit_roles.test.roles : role.name => role.id }
}

resource "aembit_identity_provider" "test_idp_oidc" {
	name = "Identity Provider OIDC for TF Acceptance Test"
	description = "Description of Identity Provider for TF Acceptance Test"
	is_active = true
    sso_statement_role_mappings = [
        {
            attribute_name = "test-attribute-name"
            attribute_value = "test-attribute-value"
            roles = [local.role_ids_by_name["SuperAdmin"], local.role_ids_by_name["Auditor"]]
        }
    ]    
    oidc = {
        oidc_base_url = "https://test.oidc.com"
        client_id = "test_client_id"
        scopes = "profile email"
        auth_type = "ClientSecret"
        client_secret = "some_secret"
        pcke_required = true
    }
}