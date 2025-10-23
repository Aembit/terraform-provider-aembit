provider "aembit" {
	default_tags {
		tags = {
			Name           = "Terraform"
			Owner          = "Aembit"    
		}	
	}
}

resource "aembit_server_workload" "test" {
	name = "Unit Test 1"
    description = "Description"
    is_active = false
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
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}

resource "aembit_server_workload" "test_wildcard" {
	name = "Unit Test 1 Wildcard"
    description = "Description"
    is_active = false
	service_endpoint = {
		host = "*.testhost.com"
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
	}
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}
