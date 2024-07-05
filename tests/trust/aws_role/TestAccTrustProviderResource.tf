provider "aembit" {
}

resource "aembit_trust_provider" "aws_role" {
	name = "TF Acceptance AWS Role"
	is_active = true
	aws_role = {
		account_id = "account_id"
		assumed_role = "assumed_role"
		role_arn = "role_arn"
		username = "username"
	}
}
