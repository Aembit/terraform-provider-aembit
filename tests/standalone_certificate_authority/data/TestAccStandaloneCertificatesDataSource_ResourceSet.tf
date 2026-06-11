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

resource "aembit_standalone_certificate_authority" "test" {
    resource_set_id = aembit_resource_set.crs.id
	name = "Unit Test 1"
    description = "Description"
	leaf_lifetime = 1440
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}

data "aembit_standalone_certificate_authorities" "test" {
    resource_set_id = aembit_resource_set.crs.id
    depends_on = [ aembit_standalone_certificate_authority.test ]
}
