provider "aembit" {
}

resource "aembit_client_workload" "test" {
    name = "TF Acceptance ProcessCommandLine"
    description = "Acceptance Test Client Workload"
    is_active = false
    identities = [
        {
            type = "processCommandLine"
            value = "*process command line*"
        },
    ]
}
