provider "aembit" {
	default_tags {
		tags = {
			Name           = "Terraform"
			Owner          = "Aembit"    
		}	
	}
}

resource "aembit_standalone_certificate_authority" "test" {
	name = "unittestname"
    description = "Description"
	leaf_lifetime = 1440
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}
