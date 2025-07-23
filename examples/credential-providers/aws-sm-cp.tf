resource "aembit_credential_provider" "aws_sm_value" {
  name        = "AWS Secrets Manager Value Credential Provider"
  description = "Credential Provider Description"
  is_active   = true
  aws_secrets_manager_value = {
    secret_arn                         = "arn:aws:secretsmanager:us-east-1:123456789012:secret:samplename-ABCDEF"
    secret_key_1                       = "key1"
    secret_key_2                       = "key2"
    private_network_access             = false
    credential_provider_integration_id = aembit_credential_provider_integration.aws_iam_role_cpi.id
  }
}