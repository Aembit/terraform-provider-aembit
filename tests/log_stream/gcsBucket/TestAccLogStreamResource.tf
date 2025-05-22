provider "aembit" {
}

resource "aembit_log_stream" "gcs_bucket" {
	name = "TF Acceptance GCSBucket LogStream"
	description = "TF Acceptance GCSBucket LogStream"
	data_type = "AuditLogs"
	type = "GcsBucket"
	is_active = false
	gcs_bucket = {
		gcs_bucket_name = "test-bucket-name"
		gcs_path_prefix = "test/test"
		audience = "test-audience"
        service_account_email = "test@test.com"
        token_lifetime = 3600
	}
}
