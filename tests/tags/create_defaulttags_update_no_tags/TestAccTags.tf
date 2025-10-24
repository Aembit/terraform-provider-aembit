provider "aembit" {
	default_tags {
		tags = {
			Name           = "Terraform"
			Owner          = "Aembit"    
		}	
	}
}

resource "aembit_client_workload" "test_tags" {
    name = "Unit Test Tags"
    is_active = false
    identities = [
        {
            type = "k8sNamespace"
            value = "unittest1namespace"
        },
    ]
}

