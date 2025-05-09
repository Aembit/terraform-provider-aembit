provider "aembit" {
    resource_set_id = "7a14a7f5-ca61-4ef8-9337-2fb70299fe81"
}
resource "aembit_client_workload" "test" {
    name = "TF Acceptance RS"
    description = "TF Acceptance CW in custom workload"
    is_active = false
    identities = [
        {
            type = "k8sPodName"
            value = "custom-resource-set"
        },
    ]
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}

