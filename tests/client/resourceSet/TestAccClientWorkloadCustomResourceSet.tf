provider "aembit" {
    resource_set_id = "4d51992e-744a-4171-a2a2-492af4ace7e6"
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

