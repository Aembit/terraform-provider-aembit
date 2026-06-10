provider "aembit" {
    alias = "rs_loader"
}

data "aembit_resource_sets" "all" {
    provider = aembit.rs_loader
}

locals {
    tf_testing_rs_id = [for rs in data.aembit_resource_sets.all.resource_sets : rs.id if rs.name == "TF Testing"][0]
}

// Create a Provider and Resource in the TF Testing Resource Set
provider "aembit" {
    alias = "rs_manager"
    resource_set_id = local.tf_testing_rs_id
}

resource "aembit_credential_provider" "oauth" {
	provider = aembit.rs_manager
	name = "TF Acceptance OAuth"
	is_active = true
	oauth_client_credentials = {
		token_url = "https://aembit.io/token"
		client_id = "test_client_id"
		client_secret = "test_client_secret"
		scopes = "test_scopes"
		credential_style = "authHeader"
		custom_parameters = [
			{
				key = "key"
				value = "value"
				value_type = "literal"
			}
		]
	}
}

data "aembit_credential_providers" "test" {
	provider = aembit.rs_manager
    depends_on = [ aembit_credential_provider.oauth ]
}
