resource "aembit_credential_provider_integration" "awsiamrole_cpi" {
	name = "TF Acceptance Aws IAM Role Credential Provider Integration"
	description = "TF Acceptance Aws IAM Role Credential Provider Integration Description"
	aws_iam_role = {
		role_arn = "arn:aws:iam::123456789012:role/MyRole"
		lifetime_in_seconds = 3600
	}
}

resource "aembit_credential_provider" "aws_sm_value" {
	name = "TF Acceptance AWS Secrets Manager Value CP"
	description = "TF Acceptance AWS Secrets Manager Value CP Description"
	is_active = true
	aws_secrets_manager_value = {
		secret_arn = "arn:aws:secretsmanager:us-east-2:123456789012:secret:secretname-ABCDEF"
		secret_key_1 = "key1"
		secret_key_2 = "key2"
		private_network_access = false
		credential_provider_integration_id = aembit_credential_provider_integration.awsiamrole_cpi.id
	}
}
