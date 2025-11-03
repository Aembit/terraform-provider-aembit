provider "aembit" {
}

resource "aembit_standalone_certificate_authority" "test" {
	name = "unittestname"
    description = "Description"
	leaf_lifetime = 1440
}
