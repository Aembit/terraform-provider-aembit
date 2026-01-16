package validators

import (
	"testing"
)

func TestRegex_UUID_Valid(t *testing.T) {
	t.Parallel()
	validInputs := []string{
		"550e8400-e29b-41d4-a716-446655440000",
		"{550e8400-e29b-41d4-a716-446655440000}",
		"a3bb189e-8bf9-3888-9912-ace4e6543002",
		"{A3BB189E-8BF9-3888-9912-ACE4E6543002}",
	}
	for _, input := range validInputs {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if !UUIDRegex.MatchString(input) {
				t.Errorf("Expected valid UUID, but got invalid: %s", input)
			}
		})
	}
}

func TestRegex_UUID_Invalid(t *testing.T) {
	t.Parallel()
	invalidInputs := []string{
		"550e8400e29b41d4a716446655440000",      // Missing dashes
		"550e8400-e29b-41d4-a716-44665544",      // Too short
		"550e8400-e29b-41d4-a716-4466554400000", // Too long
		"550e8400-e29b-41d4-a716-44665544zzzz",  // Invalid hex characters
		"550e8400-e29b-41d4-a716-44665544-0000", // Extra dash
	}
	for _, input := range invalidInputs {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if UUIDRegex.MatchString(input) {
				t.Errorf("Expected invalid UUID, but got valid: %s", input)
			}
		})
	}
}

func TestRegex_Email_Valid(t *testing.T) {
	t.Parallel()
	validInputs := []string{
		"test@example.com",
		"user.name+tag+sorting@example.com",
		"x@example.com",
		"example-indeed@strange-example.com",
	}
	for _, input := range validInputs {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if !EmailRegex.MatchString(input) {
				t.Errorf("Expected valid email, but got invalid: %s", input)
			}
		})
	}
}

func TestRegex_Email_Invalid(t *testing.T) {
	t.Parallel()
	invalidInputs := []string{
		"plainaddress",
		"@missinguser.com",
		"user@.com",
		"user@com",
	}
	for _, input := range invalidInputs {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if EmailRegex.MatchString(input) {
				t.Errorf("Expected invalid email, but got valid: %s", input)
			}
		})
	}
}

func TestRegex_SnowflakeAccount_Valid(t *testing.T) {
	t.Parallel()
	validInputs := []string{
		"test-user",
		"Testuser",
		"TESTUSER",
		"TEST-USER",
		"TEST-USER-",
	}
	for _, input := range validInputs {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if !SnowflakeAccountRegex.MatchString(input) {
				t.Errorf("Expected valid snowflake account, but got invalid: %s", input)
			}
		})
	}
}

func TestRegex_SnowflakeAccount_Invalid(t *testing.T) {
	t.Parallel()
	invalidInputs := []string{
		"test user",
		"Testuser?",
		"TESTUSER_",
	}
	for _, input := range invalidInputs {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if SnowflakeAccountRegex.MatchString(input) {
				t.Errorf("Expected invalid snowflake account, but got valid: %s", input)
			}
		})
	}
}

func TestRegex_SnowflakeUserName_Valid(t *testing.T) {
	t.Parallel()
	validInputs := []string{
		"test_user",
		"Testuser",
		"TESTUSER",
		"TEST_USER@",
		"TEST_USER_",
	}
	for _, input := range validInputs {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if !SnowflakeUserNameRegex.MatchString(input) {
				t.Errorf("Expected valid snowflake username, but got invalid: %s", input)
			}
		})
	}
}

func TestRegex_SnowflakeUserName_Invalid(t *testing.T) {
	t.Parallel()
	invalidInputs := []string{
		"test user",
		"Testuser?",
		"TEST;USER",
		"TEST-USER",
	}
	for _, input := range invalidInputs {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if SnowflakeUserNameRegex.MatchString(input) {
				t.Errorf("Expected invalid snowflake username, but got valid: %s", input)
			}
		})
	}
}

func TestRegex_URLScheme_Valid(t *testing.T) {
	t.Parallel()
	validInputs := []string{
		"http://test.com",
		"https://test.com",
		"https://localhost",
	}
	for _, input := range validInputs {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if !UrlSchemeRegex.MatchString(input) {
				t.Errorf("Expected valid url scheme, but got invalid: %s", input)
			}
		})
	}
}

func TestRegex_URLScheme_Invalid(t *testing.T) {
	t.Parallel()
	invalidInputs := []string{
		"htttp://test.com",
		"https//test.com",
		"httpslocalhost",
	}
	for _, input := range invalidInputs {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if UrlSchemeRegex.MatchString(input) {
				t.Errorf("Expected invalid url scheme, but got valid: %s", input)
			}
		})
	}
}

func TestRegex_SecureURL_Valid(t *testing.T) {
	t.Parallel()
	validInputs := []string{
		"https://test.com",
	}
	for _, input := range validInputs {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if !SecureURLRegex.MatchString(input) {
				t.Errorf("Expected a valid secure url, but got invalid: %s", input)
			}
		})
	}
}

func TestRegex_SecureURL_Invalid(t *testing.T) {
	t.Parallel()
	invalidInputs := []string{
		"http://test.com",
		"http://localhost",
	}
	for _, input := range invalidInputs {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if SecureURLRegex.MatchString(input) {
				t.Errorf("Expected invalid secure url, but got valid: %s", input)
			}
		})
	}
}

func TestRegex_Hostname_Valid(t *testing.T) {
	t.Parallel()
	validInputs := []string{
		"example.com",
		"sub.domain.com",
		"my-site.co.uk",
		"localhost",
		"example.travel",
	}
	for _, input := range validInputs {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if !HostNameRegex.MatchString(input) {
				t.Errorf("Expected valid hostname, but got invalid: %s", input)
			}
		})
	}
}

func TestRegex_Hostname_Invalid(t *testing.T) {
	t.Parallel()
	validInputs := []string{
		"-invalid.com",
		"example..com",
		"example_com",
		".example.com",
		"example.com-",
	}
	for _, input := range validInputs {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if HostNameRegex.MatchString(input) {
				t.Errorf("Expected invalid hostname, but got valid: %s", input)
			}
		})
	}
}

func TestRegex_SafeWildcardHostname_Valid(t *testing.T) {
	t.Parallel()
	validInputs := []string{
		"example.com",
		"sub.domain.com",
		"my-site.co.uk",
		"localhost",
		"example.travel",
		"aembit.io",
		"host",
		"test.com",
		"*.amazonaws.com",
		"bucket.*.amazonaws.com",
		"bucket.*.microsoft.com",
		"wild*card.amazonaws.com",
	}
	for _, input := range validInputs {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if !StructuralHostRegex.MatchString(input) {
				t.Errorf("Expected valid structural hostname, but got invalid: %s", input)
			}
		})
	}
}

func TestRegex_SafeWildcardHostname_Invalid(t *testing.T) {
	t.Parallel()
	validInputs := []string{
		"-invalid.com",
		"example..com",
		"example_com",
		".example.com",
		"example.com-",
		"google.*",
		"google.n*t",
		"*.com",
		"goo*gle.com",
		"", // Empty string
	}
	for _, input := range validInputs {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if StructuralHostRegex.MatchString(input) {
				t.Errorf("Expected invalid structural hostname, but got valid: %s", input)
			}
		})
	}
}

func TestRegex_IPHost_Valid(t *testing.T) {
	t.Parallel()
	validInputs := []string{
		"192.168.1.1", // NOSONAR
		"10.0.0.1",    // NOSONAR
		"172.16.0.1",  // NOSONAR
	}
	for _, input := range validInputs {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if !HostIPRegex.MatchString(input) {
				t.Errorf("Expected valid IP host, but got invalid: %s", input)
			}
		})
	}
}

func TestRegex_IPHost_Invalid(t *testing.T) {
	t.Parallel()
	validInputs := []string{
		"192.*.86.1", // IP with wildcard
		"",           // Empty string
	}
	for _, input := range validInputs {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if HostIPRegex.MatchString(input) {
				t.Errorf("Expected invalid IP host, but got valid: %s", input)
			}
		})
	}
}

func TestAwsIamRoleArnRegex(t *testing.T) {
	t.Parallel()
	valid := "arn:aws:iam::123456789012:role/MyRole"
	invalid := "arn:aws:iam::123456789012:role/"

	if !AwsIamRoleArnRegex.MatchString(valid) {
		t.Errorf("expected valid AWS IAM Role ARN to match regex")
	}
	if AwsIamRoleArnRegex.MatchString(invalid) {
		t.Errorf("expected invalid AWS IAM Role ARN to not match regex")
	}
}

func TestAwsSecretArnRegex(t *testing.T) {
	t.Parallel()
	valid := "arn:aws:secretsmanager:us-west-2:123456789012:secret:mysecret-ABC123"
	invalid := "arn:aws:secretsmanager:us-west-2:123456789012:secret:mysecret"

	if !AwsSecretArnRegex.MatchString(valid) {
		t.Errorf("expected valid AWS Secret ARN to match regex")
	}
	if AwsSecretArnRegex.MatchString(invalid) {
		t.Errorf("expected invalid AWS Secret ARN to not match regex")
	}
}
