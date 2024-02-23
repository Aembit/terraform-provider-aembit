provider "aembit" {
}

resource "aembit_client_workload" "first_client" {
    name = "first terraform client workload"
    description = "new client workload for policy integration"
    is_active = false
    identities = [
        {
            type = "k8sNamespace"
            value = "clientworkload"
        },
    ]
}

resource "aembit_server_workload" "first_server" {
    name = "first terraform server workload"
    description = "new server workload for policy integration"
    is_active = false
    service_endpoint = {
        host = "myhost.unittest.com"
        port = 443
        app_protocol = "HTTP"
        transport_protocol = "TCP"
        requested_port = 80
        tls_verification = "full"
	    requested_tls = true
	    tls = true
    }
}

resource "aembit_access_policy" "first_policy" {
    name = "policy1"
    description = "new policy for test"
    is_active = false
    client_workload = aembit_client_workload.first_client.id
    server_workload = aembit_server_workload.first_server.id
    //server_workload = "b7b82d82-a0ba-4c1e-8d5e-a6b4e6533731"
}