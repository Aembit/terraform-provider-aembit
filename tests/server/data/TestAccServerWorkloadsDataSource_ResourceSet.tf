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

resource "aembit_server_workload" "test" {
    resource_set_id = aembit_resource_set.crs.id
	name = "Unit Test 1"
    description = "Description"
    is_active = true
	service_endpoint = {
		host = "unittest.testhost.com"
		port = 443
        tls = true
		app_protocol = "HTTP"
		transport_protocol = "TCP"
		requested_port = 443
        requested_tls = true
		tls_verification = "full"
		authentication_config = {
			"method" = "HTTP Authentication"
			"scheme" = "Bearer"
		}
		http_headers = {
			host = "graph.microsoft.com"
			user-agent = "curl/7.64.1"
			accept = "*/*"
		}
	}
}

data "aembit_server_workloads" "test" {
	resource_set_id = aembit_resource_set.crs.id
    depends_on = [ aembit_server_workload.test ]
}
