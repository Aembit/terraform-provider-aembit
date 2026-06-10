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

resource "aembit_client_workload" "test" {
    resource_set_id = aembit_resource_set.crs.id
    name = "Unit Test 1"
    description = "Acceptance Test client workload"
    is_active = false
    identities = [
        {
            type = "k8sNamespace"
            value = "unittest1namespace"
        },
    ]
}

data "aembit_client_workloads" "test" {
    depends_on = [ aembit_client_workload.test ]
}
