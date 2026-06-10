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

resource "aembit_trust_provider" "kubernetes" {
	provider = aembit.rs_manager
	name = "TF Acceptance Kubernetes"
	is_active = true
	kubernetes_service_account = {
		issuer = "issuer"
		namespace = "namespace"
		pod_name = "pod_name"
		service_account_name = "service_account_name"
		subject = "subject"
		oidc_endpoint = "https://3a3b5d.id.devbroadangle.aembit-eng.com/"
	}
}

resource "aembit_trust_provider" "kubernetes_key" {
	provider = aembit.rs_manager
	name = "TF Acceptance Kubernetes Key"
	is_active = true
	kubernetes_service_account = {
		issuer = "issuer"
		namespace = "namespace"
		pod_name = "pod_name"
		service_account_name = "service_account_name"
		subject = "subject"
		public_key = <<-EOT
-----BEGIN PUBLIC KEY-----
MFswDQYJKoZIhvcNAQEBBQADSgAwRwJAXWRPQyGlEY+SXz8Uslhe+MLjTgWd8lf/
nA0hgCm9JFKC1tq1S73cQ9naClNXsMqY7pwPt1bSY8jYRqHHbdoUvwIDAQAB
-----END PUBLIC KEY-----
EOT
	}
}

data "aembit_trust_providers" "test" {
	provider = aembit.rs_manager
    depends_on = [
        aembit_trust_provider.kubernetes,
        aembit_trust_provider.kubernetes_key
     ]
}
