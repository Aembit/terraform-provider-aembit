provider "aembit" {
}

resource "aembit_standalone_certificate_authority" "first_ca" {
    name = "unittestname"
    is_active = true
    leaf_lifetime = 60
}

resource "aembit_client_workload" "test" {
    name = "unittestname"
    description = "TF Acceptance CW with Standalone CA"
    is_active = false
    identities = [
        {
            type = "k8sPodName"
            value = "sample_ca"
        },
    ]
    tags = {
        color = "blue"
        day   = "Sunday"
    }
    standalone_certificate_authority = aembit_standalone_certificate_authority.first_ca.id
}