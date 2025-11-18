provider "aembit" {
}

resource "aembit_trust_provider" "certificate_signed_attestation" {
	name = "TF Acceptance CertificateSignedAttestation"
	is_active = true
	certificate_signed_attestation = {		
	}
}

resource "aembit_trust_provider_secret" "secret" {
	trust_provider_id = aembit_trust_provider.certificate_signed_attestation.id
	secret      = <<-EOT
-----BEGIN CERTIFICATE-----
MIIC8jCCAdqgAwIBAgIJAMIlrM/czmegMA0GCSqGSIb3DQEBCwUAMBsxGTAXBgNVBAMTEFRlc3Qg
Q2VydGlmaWNhdGUwHhcNMjUxMTE2MTU0MTEyWhcNMjUxMjAyMTU0MTEyWjAbMRkwFwYDVQQDExBU
ZXN0IENlcnRpZmljYXRlMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsjBXXaG9PF3Q
54xh3xr39IoW42RM6YTHWm2QgcD56uZW2YXkhLXvbiDWXA1RPa4Tfhy8UheimNosE0t3/fudB3jL
J8S22k/60X9tuGosWU5rV68/Xnhsd+gFp1SwB2S1icRQKxiqxwtbH1ukYMSWcQ4kT/U0bbe/9G5V
QWr8NHpPgTGMSIEzSyRLaEBtAvZsqkv+4bGfXDBo/AFqwGyJavhrln2I/550QyjtlthhqaH/lrnJ
WirThM6M9s0RU7W2h6El37AT7/PUFd5nzi8JVFkerw8nzS6fLYcUgJb5CTVwIy/wKGf3mAKPlIN0
JpX2OVt9nMe0sLt6972fjH/CvQIDAQABozkwNzAJBgNVHRMEAjAAMAsGA1UdDwQEAwIFoDAdBgNV
HQ4EFgQUU+7wxT6cp2uaux6+/gBK3h1YOjswDQYJKoZIhvcNAQELBQADggEBAI2t0HxjH26BND6H
QbJCCu7cj/tuGAaFIuubeccrtnHsbMkKVmOCxfNm2PPDcxcGXQf1pGbRnYWFwR/UkJm+ZlGIZfXR
Z8F/+wor+hya0uSXppD1Qd3K/l1BiJAn5huMvlGppx56mzlqmB+h8ahZKeW0YvweJYpvuRFlL/p5
rwfsG/EGzCYaJeS7o/k7dxtaDw5XsHvNVZy72nboG9hEUMwyfeCW/rYHbNWZDe3g08NfwrDosuZm
ZdKgpyi48slfLVoiYRJPkZqF+TP8JBq3EIkH7xSdUtpkoLR6bgnTfYvWTW+4KS3vXJhHnlFYNBrk
Ph3gNoIYBZi/TerC1Sao0k8=
-----END CERTIFICATE-----
EOT	
}