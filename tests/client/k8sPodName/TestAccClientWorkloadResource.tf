provider "aembit" {
    alias = "rs_loader"
}

data "aembit_resource_sets" "all" {
    provider = aembit.rs_loader
}

// Create a Provider and Resource in the second Resource Set
provider "aembit" {
    alias = "rs_manager"
    resource_set_id = data.aembit_resource_sets.all.resource_sets[1].id
}

resource "aembit_client_workload" "test" {
    provider = aembit.rs_manager

    name = "Unit Test 1 - In Resource Set"
    description = "Acceptance Test client workload"
    is_active = false
    identities = [
        {
            type = "k8sPodName"
            value = "unittest1podname"
        },
    ]
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}

