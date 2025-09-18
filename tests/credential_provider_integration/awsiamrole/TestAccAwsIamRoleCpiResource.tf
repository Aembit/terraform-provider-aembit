resource "aembit_credential_provider_integration" "awsiamrole_cpi" {
	name = "TF Acceptance Aws IAM Role Credential Provider Integration"
	description = "TF Acceptance Aws IAM Role Credential Provider Integration Description"
	aws_iam_role = {
		role_arn = "arn:aws:iam::123456789012:role/MyRole"
		lifetime_in_seconds = 3600
		fetch_secret_arns = false
	}
}
