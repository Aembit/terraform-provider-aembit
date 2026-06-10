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

resource "aembit_credential_provider_integration" "aws_sm_secret" {
	provider = aembit.rs_manager
	name        = "AWS IAM Role Credential Provider Integration For Data Source Test"
  	aws_iam_role = {
    	role_arn            = "arn:aws:iam::123456789012:role/MyRole"
    	lifetime_in_seconds = 3600
  	}
}

data "aembit_credential_provider_integrations" "test" {
	provider = aembit.rs_manager
    depends_on = [ aembit_credential_provider_integration.aws_sm_secret ]
}
