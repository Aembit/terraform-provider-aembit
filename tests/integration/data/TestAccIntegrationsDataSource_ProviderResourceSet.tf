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

resource "aembit_integration" "wiz" {
	provider = aembit.rs_manager
	name = "TF Acceptance Wiz"
	is_active = false
	type = "WizIntegrationApi"
	sync_frequency = 3600
	endpoint = "https://api.us17.app.wiz.io/graphql"
	oauth_client_credentials = {
		token_url = "https://auth.app.wiz.io/oauth/token"
		client_id = "client_id"
		client_secret = "client_secret"
		audience = "audience"
	}
}

data "aembit_integrations" "test" {
	provider = aembit.rs_manager
    depends_on = [ aembit_integration.wiz ]
}
