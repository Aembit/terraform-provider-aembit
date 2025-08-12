provider "aembit" {
}

resource "aembit_standalone_certificate_authority" "first_ca" {
    name = "unittestname"
    leaf_lifetime = 60
}

resource "aembit_client_workload" "test" {
    name = "Unit Test 1"
    description = "Acceptance Test client workload"
    is_active = false
    identities = [
        {
            type = "k8sNamespace"
            value = "unittest1namespace"
        },
    ]
    tags = {
        color = "blue"
        day   = "Sunday"
    }
    standalone_certificate_authority = aembit_standalone_certificate_authority.first_ca.id
}

data "aembit_client_workloads" "test" {
    depends_on = [ aembit_client_workload.test ]
}

