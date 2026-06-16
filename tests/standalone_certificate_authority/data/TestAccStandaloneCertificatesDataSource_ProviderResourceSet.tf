provider "aembit" {
    alias = "rs_loader"
}

data "aembit_resource_sets" "all" {
    provider = aembit.rs_loader
}

locals {
    tf_testing_rs_id = [for rs in data.aembit_resource_sets.all.resource_sets : rs.id if rs.name == "TF Testing"][0]
}

// Create a Provider and Resource in the TF Testing Resource Set
provider "aembit" {
    alias = "rs_manager"
    resource_set_id = local.tf_testing_rs_id
}

resource "aembit_standalone_certificate_authority" "test" {
    provider = aembit.rs_manager
	name = "Unit Test 1"
    description = "Description"
	leaf_lifetime = 1440
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}

data "aembit_standalone_certificate_authorities" "test" {
    provider = aembit.rs_manager
    depends_on = [ aembit_standalone_certificate_authority.test ]
}
