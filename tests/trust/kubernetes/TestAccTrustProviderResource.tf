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
		jwks = <<-EOT
{
  "keys": [
    {
      "kty": "RSA",
      "use": "sig",
      "kid": "T41hVcPtA3ehDjSaZXSI9LKuanyTkBOf0YKlAM6gtNQ=",
      "e": "AQAB",
      "n": "vEJ5_IKWyoGjoB-Us5uooNWR0dvvTC_8eilRWPth1LxsbAahxORlOO8asmFc0C1pDwIo74XZlbwfLfet8Q0WzSre_8IJHDStiQUgiDPnh9Z5vDIH3HoSVQIOW9W4AIdYeQd5iW7hGVucwm6eal3jv3sF1CvvYZT77vf8bBFKl26xr_cIpsl77wECIFij6dR_dtE59g7etsz1EeDvwm75OOgNL7z-bCum149E7luyE5y7bNtpqtbthQK31vyaifrGABYXragi4vWcw7yWif1IV7M_smlZBHPeGbRZ4xCKiVkL7vtwz6AgW8BfhewGI4_qQfEONAXEJv70VK6OpJ5oZw",
      "alg": "RS256"
    },
    {
      "kty": "EC",
      "use": "sig",
      "kid": "iTYG7jb2cyaQ04cp69CpoMBzxNjRmixlGGxZTIHSpXg=",
      "alg": "ES256",
      "x": "ItvdSxTnkMqPq3kKeHYlAF1ArZqz4_CXjUmiPvHDQ08",
      "y": "C2z0b9zNhvywzboDt03F2xb_7fOaw8LWbakgudjN3kE",
      "crv": "P-256"
    }
  ]
}
EOT
	}
}

resource "aembit_trust_provider" "kubernetes_symmetric_key" {
	name = "TF Acceptance Kubernetes Symmetric Key"
	is_active = true
	kubernetes_service_account = {
		issuer = "issuer"
		namespace = "namespace"
		pod_name = "pod_name"
		service_account_name = "service_account_name"
		subject = "subject"
		symmetric_key = "a2V5MQ=="
	}
}