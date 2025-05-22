provider "aembit" {
}

resource "aembit_client_workload" "test" {
    name = "Unit Test 1 - miscellaneous"
    description = "Acceptance Test client workload"
    is_active = false
    identities = [
        {
            type = "IDENTITY_TYPE_PLACEHOLDER"
            value = "IDENTITY_VALUE_PLACEHOLDER"
        },
    ]
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}

