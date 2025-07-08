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

resource "aembit_trust_provider" "kubernetes_jwks" {
	name = "TF Acceptance Kubernetes JWKS"
	is_active = true
	kubernetes_service_account = {
		issuer = "issuer"
		namespace = "namespace"
		pod_name = "pod_name"
		service_account_name = "service_account_name"
		subject = "subject"
		jwks = {
			keys = [
				{
					kid = "Tbm3LtlhYlNObRdRc+Tz3mEo2SASPbfR03HI4dmoUkg="
					kty = "RSA"
					use = "sig"
					alg = "RS256"
					e = "AQAB"
					n = "z5DRwd-vq2lhTMklKDezYv9L1pBhG1qVzzMK0vFZB-QqUeYgZ5Ky3Ie74xOJzIfRrbCAhamVNZKR_H-4YpokLVTnw_Wu118EonMWSeuUvUsEIpV7EFzpu4H-JSqa3Ynq9cYE7MvgC5nXkdivQCuze3ZFuO9RZatWGtzhGyFsU7JEYTnEf0Hues3JH4Hk5-Crux5KX7KBGu1-ecTL_cXO1swx2Q5bAX1knsNptE7c-hLn0kTzZWFsY6l880G_AotVbrNgFuKc9JRFWGTroVbpd4JK2vASPbIpQuVnMXAhtLih7_YKyLng5dbrVRwwOYxSQ_Tn1OesE9XnYaJts8dnVw"
				},
				{
					kid = "mxAhc1VybhA8LT2jHRQFEzcWSoLbFmnDWGYoViS/aKw="
					kty = "EC"
					use = "sig"
					alg = "ES256"
					x: "x-pRNOyN2BwmgvPuLTOEJMLB1vcc4vljjU41W0jz5Sw",
					y: "7beLCvWqWomVizZAhrxR2vsttzD3owKnE__ZADccuyk",
					crv: "P-256"				
				}
			]			
		}
	}
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}