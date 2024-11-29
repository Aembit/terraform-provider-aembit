provider "aembit" {
}

resource "aembit_signin_policy" "test" {
	sso_required = true
	mfa_required = true
}