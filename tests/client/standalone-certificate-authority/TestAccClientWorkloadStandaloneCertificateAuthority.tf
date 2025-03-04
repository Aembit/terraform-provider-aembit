provider "aembit" {
}

resource "aembit_standalone_certificate_authority" "first_ca" {
    name = "first terraform standalone certificate authority"
    is_active = true
    leaf_lifetime = 60
}

resource "aembit_client_workload" "test" {
    name = "TF Acceptance Standalone CA"
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