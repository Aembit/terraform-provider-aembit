provider "aembit" {
}

data "aembit_roles" "test" {}

locals {
  role_ids_by_name = { for role in data.aembit_roles.test.roles : role.name => role.id }
}

resource "aembit_resource_set" "crs" {
	name = "TF Acceptance Custom ResourceSet"
	description = "TF Acceptance Custom ResourceSet"
	roles = [local.role_ids_by_name["SuperAdmin"], local.role_ids_by_name["Auditor"]]
}

resource "aembit_credential_provider_integration" "aws_sm_secret" {
	resource_set_id = aembit_resource_set.crs.id
	name        = "AWS IAM Role Credential Provider Integration For Data Source Test"
  	aws_iam_role = {
    	role_arn            = "arn:aws:iam::123456789012:role/MyRole"
    	lifetime_in_seconds = 3600
  	}
}

data "aembit_credential_provider_integrations" "test" {
    depends_on = [ aembit_credential_provider_integration.aws_sm_secret ]
}
