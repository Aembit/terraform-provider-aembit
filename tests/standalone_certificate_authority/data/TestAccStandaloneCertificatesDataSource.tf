resource "aembit_standalone_certificate_authority" "test" {
	name = "Unit Test 1"
    description = "Description"
	leaf_lifetime = 1440
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}

provider "aembit" {
}

data "aembit_standalone_certificate_authorities" "test" {
    depends_on = [ aembit_standalone_certificate_authority.test ]
}
