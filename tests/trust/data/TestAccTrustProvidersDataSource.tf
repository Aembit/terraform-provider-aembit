provider "aembit" {
}

resource "aembit_trust_provider" "kubernetes" {
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
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}

resource "aembit_trust_provider" "kubernetes_key" {
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
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}

data "aembit_trust_providers" "test" {
    depends_on = [
        aembit_trust_provider.kubernetes,
        aembit_trust_provider.kubernetes_key
     ]
}
