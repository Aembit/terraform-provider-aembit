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
MIIC5DCCAcygAwIBAgIUSq9DjW7fVmSbfiVT0loPvJWQWxUwDQYJKoZIhvcNAQEL
BQAwGzEZMBcGA1UEAwwQVGVzdCBDZXJ0aWZpY2F0ZTAeFw0yNTEyMDMwOTM0MTRa
Fw0zMDEyMDMwOTM0MTRaMBsxGTAXBgNVBAMMEFRlc3QgQ2VydGlmaWNhdGUwggEi
MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDQbCe7czS7oSMx9fw7UTLWngKw
B9JusJju/awXCAh9ogHWbiwy70YZu8eby37ND1N3mY9T5Gb7s/AJTunbx8WC0H86
AZilakVO++IO5zf6AuBWCYxMA3Z9rlpqvfXreQ6RDBEZaXHTE+6Srv8eiVghnsu6
8aBmxicnyiSecUyrOt0/CryncwfNgKV0vtnlU8j/8FksGL4t7o/bA3hWxfLKsDe9
zS/BbSrjpyV84yKalBom3AILfiMI435gFZNUzQcpYDSyWaH1o2soyommZYe011rx
nwHrTVceolIiGuUVADYLibsA5ZMkshAtr1uFbOJFpZy9mYtXhi98EpG8IMAlAgMB
AAGjIDAeMAwGA1UdEwEB/wQCMAAwDgYDVR0PAQH/BAQDAgWgMA0GCSqGSIb3DQEB
CwUAA4IBAQA3GkiPASfwRFfnnIYFw8C41wd4I6cYeU+8CxeoLSgKY2tIvgfJvo9h
3Np5H3YesqZQDAL68jBJVO3xQ2D3ZJOL9sgMJC+4DVRYjgvUdKdoIe0krz0mw8Pp
476a8+6IRUZW1p6Z1cE+uDeiX1kEa5+ezq6J6ileXGkzFCwpShzrTcCVP2X3CXnR
VrumnNr19wn5A+w4wpcFMqIyE/2Z9v/yvY5/0ac37j0gURURWXleIfr8zarXYCqf
vCRP/IkVBQOxwhPu1VUIaTP2uup+FTP8Dwn9uWWEvAFuBQkmWuFbdLcaUZLiXVuI
aTOJ++hdzYrm4Wa+bIGqiMpQDUjwz3aC
-----END CERTIFICATE-----
EOT	
}