resource "aembit_credential_provider_integration" "aws_iam_role_cpi" {
  name        = "AWS IAM Role Credential Provider Integration"
  description = "Detailed description of AWS IAM Role Credential Provider Integration"
  aws_iam_role = {
    role_arn            = "arn:aws:iam::123456789012:role/MyRole"
    lifetime_in_seconds = 3600
  }
}
