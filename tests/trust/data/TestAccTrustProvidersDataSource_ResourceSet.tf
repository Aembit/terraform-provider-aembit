provider "aembit" {
}

data "aembit_roles" "test" {}

locals {
  role_ids_by_name = { for role in data.aembit_roles.test.roles : role.name => role.id }
}

resource "aembit_resource_set" "crs" {
	name = "TF Acceptance Custom ResourceSet"
	description = "TF Acceptance Custom ResourceSet"
	roles = [local.role_ids_by_name["SuperAdmin"], local.role_ids_by_name["Auditor"]]
}

resource "aembit_trust_provider" "kubernetes" {
	resource_set_id = aembit_resource_set.crs.id 
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
 	resource_set_id = aembit_resource_set.crs.id
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
    depends_on = [
        aembit_trust_provider.kubernetes,
        aembit_trust_provider.kubernetes_key
     ]
}
