provider "aembit" {
}

resource "aembit_credential_provider_integration" "aws_sm_secret" {
	name        = "AWS IAM Role Credential Provider Integration For Data Source Test"
  	aws_iam_role = {
    	role_arn            = "arn:aws:iam::123456789012:role/MyRole"
    	lifetime_in_seconds = 3600
  	}
}

data "aembit_credential_provider_integrations" "test" {
    depends_on = [ aembit_credential_provider_integration.aws_sm_secret ]
}
