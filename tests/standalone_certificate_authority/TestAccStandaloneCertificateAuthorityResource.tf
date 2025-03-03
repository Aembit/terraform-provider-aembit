provider "aembit" {
}

resource "aembit_standalone_certificate_authority" "test" {
	name = "Unit Test 1"
    description = "Description"
    is_active = true
	leaf_lifetime = 1440
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}
