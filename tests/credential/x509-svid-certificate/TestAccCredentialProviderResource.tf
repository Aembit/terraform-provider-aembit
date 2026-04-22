provider "aembit" {
}

resource "aembit_standalone_certificate_authority" "first_ca" {
    name = "unittestname"
    leaf_lifetime = 60
}

resource "aembit_credential_provider" "x509_svid_certificate" {
	name = "TF Acceptance X.509-SVID Certificate"
	is_active = true
	x509_svid_certificate = {
		subject = "subject"
		subject_type = "literal"
		spiffe_id = "spiffe://test.com/path"
		lifetime_in_minutes = 120
		algorithm_type = "ES256"
		id_kp_client_auth = true
		id_kp_server_auth = false
		standalone_certificate_authority = aembit_standalone_certificate_authority.first_ca.id
	}
}
