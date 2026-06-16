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

resource "aembit_server_workload" "test" {
    provider = aembit.rs_manager
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
    provider = aembit.rs_manager
    depends_on = [ aembit_server_workload.test ]
}
