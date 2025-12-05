provider "aembit" {
}

resource "aembit_trust_provider" "certificate_signed_attestation" {
	name = "TF Acceptance CertificateSignedAttestation"
	is_active = true
	certificate_signed_attestation = {		
	}
}

resource "aembit_trust_provider_secret" "secret_certificate" {
	trust_provider_id = aembit_trust_provider.certificate_signed_attestation.id
	secret      = <<-EOT
PEM_CERTIFICATE
EOT	
}