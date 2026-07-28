provider "aembit" {
}

resource "aembit_trust_provider" "aws_alb_jwt" {
	name = "TF Acceptance AWS ALB JWT"
	is_active = true
	aws_alb_jwt = {
		issuer = "issuer"
		subject = "subject"
		audience = "audience"
		email = "test@example.com"
	}
}

resource "aembit_trust_provider" "aws_alb_jwt_multi" {
	name = "TF Acceptance AWS ALB JWT Multi"
	is_active = true
	aws_alb_jwt = {
		issuers = ["issuer1", "issuer2"]
		subjects = ["subject1", "subject2"]
		audiences = ["audience1", "audience2"]
		emails = ["test1@example.com", "test2@example.com"]
	}
}

resource "aembit_trust_provider" "aws_alb_jwt_customclaims" {
	name = "TF Acceptance AWS ALB JWT Custom Claims"
	is_active = true
	aws_alb_jwt = {
		custom_claims = [
			{
				claim_key   = "key1"
				claim_value = "value1"
			},
			{
				claim_key   = "key2"
				claim_value = "value2"
			}
		]
	}
}
