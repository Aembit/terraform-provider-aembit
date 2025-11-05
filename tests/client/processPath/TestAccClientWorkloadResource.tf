provider "aembit" {
}

resource "aembit_client_workload" "test" {
    name = "TF Acceptance ProcessPath"
    description = "Acceptance Test Client Workload"
    is_active = false
    identities = [
        {
            type = "processPath"
            value = "/process/path"
        },
    ]
}
