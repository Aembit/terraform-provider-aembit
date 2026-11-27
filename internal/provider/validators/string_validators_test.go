package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestCimdUrlValidation(t *testing.T) {
	t.Parallel()

	v := CimdUrlValidation()
	require.NotNil(t, v)

	ctx := context.Background()

	// Verify Description and MarkdownDescription
	require.Contains(t, v.Description(ctx), "CIMD client ID URL")
	require.Contains(t, v.MarkdownDescription(ctx), "CIMD client ID URL")

	validInputs := []string{
		"https://mcpjam.com/.well-known/oauth/client-metadata.json",
		"https://example.com/some/path",
		"https://sub.example.com/path/to/resource",
		"https://a.b.co/d",
	}

	for _, input := range validInputs {
		t.Run("Valid_"+input, func(t *testing.T) {
			req := validator.StringRequest{
				ConfigValue: types.StringValue(input),
				Path:        path.Root("cimd_url"),
			}
			resp := &validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}
			v.ValidateString(ctx, req, resp)
			require.False(t, resp.Diagnostics.HasError(), "Expected valid CIMD URL, but got error: %s", resp.Diagnostics)
		})
	}

	invalidInputs := []string{
		"http://example.com/path",             // Not https
		"https://example.com",                 // No path
		"https://example.com/",                // Path is just /
		"https://example.com/path?query=1",    // Query string not allowed
		"https://example.com/path#fragment",   // Fragment not allowed
		"https://user@example.com/path",       // User info not allowed
		"https://example.com/path with space", // Spaces not allowed
	}

	for _, input := range invalidInputs {
		t.Run("Invalid_"+input, func(t *testing.T) {
			req := validator.StringRequest{
				ConfigValue: types.StringValue(input),
				Path:        path.Root("cimd_url"),
			}
			resp := &validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}
			v.ValidateString(ctx, req, resp)
			require.True(t, resp.Diagnostics.HasError(), "Expected invalid CIMD URL to fail validation, but it succeeded: %s", input)
		})
	}
}

func TestSpiffeSubjectValidation(t *testing.T) {
	t.Parallel()

	v := SpiffeSubjectValidation()
	require.NotNil(t, v)

	ctx := context.Background()

	// Verify Description and MarkdownDescription
	require.NotEmpty(t, v.Description(ctx))
	require.NotEmpty(t, v.MarkdownDescription(ctx))

	validInputs := []string{
		"spiffe://trust-domain-name/path",
		"spiffe://example.org/ns/prod/sa/app",
		"spiffe://foo/bar",
	}

	for _, input := range validInputs {
		t.Run("Valid_"+input, func(t *testing.T) {
			req := validator.StringRequest{
				ConfigValue: types.StringValue(input),
				Path:        path.Root("spiffe_subject"),
			}
			resp := &validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}
			v.ValidateString(ctx, req, resp)
			require.False(t, resp.Diagnostics.HasError(), "Expected valid SPIFFE subject, but got error: %s", resp.Diagnostics)
		})
	}

	invalidInputs := []string{
		"http://example.com",
		"https://example.com",
		"not-spiffe://domain/path",
		"spiffe:",
		"spiffe:/",
		"spiffe",
	}

	for _, input := range invalidInputs {
		t.Run("Invalid_"+input, func(t *testing.T) {
			req := validator.StringRequest{
				ConfigValue: types.StringValue(input),
				Path:        path.Root("spiffe_subject"),
			}
			resp := &validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}
			v.ValidateString(ctx, req, resp)
			require.True(t, resp.Diagnostics.HasError(), "Expected invalid SPIFFE subject to fail validation, but it succeeded: %s", input)
		})
	}
}

func TestIdentityProviderIDPrefixValidation(t *testing.T) {
	t.Parallel()

	v := IdentityProviderIDPrefixValidation()
	require.NotNil(t, v)

	ctx := context.Background()

	// Verify Description and MarkdownDescription
	require.Contains(t, v.Description(ctx), "must be prefixed with \"idp_\"")
	require.Contains(t, v.MarkdownDescription(ctx), "must be prefixed with \"idp_\"")

	validInputs := []string{
		"idp_test",
		"idp_12345",
		"idp_some-idp-id",
	}

	for _, input := range validInputs {
		t.Run("Valid_"+input, func(t *testing.T) {
			req := validator.StringRequest{
				ConfigValue: types.StringValue(input),
				Path:        path.Root("identity_provider_id"),
			}
			resp := &validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}
			v.ValidateString(ctx, req, resp)
			require.False(t, resp.Diagnostics.HasError(), "Expected valid Identity Provider ID, but got error: %s", resp.Diagnostics)
		})
	}

	invalidInputs := []string{
		"id_test",
		"test_idp_",
		"idp",
		"123-idp_",
		"idp_",
	}

	for _, input := range invalidInputs {
		t.Run("Invalid_"+input, func(t *testing.T) {
			req := validator.StringRequest{
				ConfigValue: types.StringValue(input),
				Path:        path.Root("identity_provider_id"),
			}
			resp := &validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}
			v.ValidateString(ctx, req, resp)
			require.True(t, resp.Diagnostics.HasError(), "Expected invalid Identity Provider ID to fail validation, but it succeeded: %s", input)
		})
	}
}

func TestServiceAccountIDPrefixValidation(t *testing.T) {
	t.Parallel()

	v := ServiceAccountIDPrefixValidation()
	require.NotNil(t, v)

	ctx := context.Background()

	// Verify Description and MarkdownDescription
	require.Contains(t, v.Description(ctx), "must be prefixed with \"user-\"")
	require.Contains(t, v.MarkdownDescription(ctx), "must be prefixed with \"user-\"")

	validInputs := []string{
		"user-test",
		"user-12345",
		"user-some-service-account-id",
	}

	for _, input := range validInputs {
		t.Run("Valid_"+input, func(t *testing.T) {
			req := validator.StringRequest{
				ConfigValue: types.StringValue(input),
				Path:        path.Root("service_account_id"),
			}
			resp := &validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}
			v.ValidateString(ctx, req, resp)
			require.False(t, resp.Diagnostics.HasError(), "Expected valid Service Account ID, but got error: %s", resp.Diagnostics)
		})
	}

	invalidInputs := []string{
		"usr-test",
		"test-user-",
		"user",
		"123-user-",
		"user-",
	}

	for _, input := range invalidInputs {
		t.Run("Invalid_"+input, func(t *testing.T) {
			req := validator.StringRequest{
				ConfigValue: types.StringValue(input),
				Path:        path.Root("service_account_id"),
			}
			resp := &validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}
			v.ValidateString(ctx, req, resp)
			require.True(t, resp.Diagnostics.HasError(), "Expected invalid Service Account ID to fail validation, but it succeeded: %s", input)
		})
	}
}

func TestUrlPathValidation(t *testing.T) {
	t.Parallel()

	v := UrlPathValidation()
	require.NotNil(t, v)

	ctx := context.Background()

	require.Contains(t, v.Description(ctx), "must be a valid URL path")
	require.Contains(t, v.MarkdownDescription(ctx), "must be a valid URL path")

	validInputs := []string{
		"/",
		"/api",
		"/api/",
		"/api/v1/resource",
		"/api/v1/resource/",
		"/api/v3/agentspace/mcp/",
		"/path-with_special~chars.123",
		"/path/with%20encoded%2Fvalues/",
		"/sub/dir;param=value:123@host+test$other,test*ok",
	}

	for _, input := range validInputs {
		t.Run("Valid_"+input, func(t *testing.T) {
			req := validator.StringRequest{
				ConfigValue: types.StringValue(input),
				Path:        path.Root("url_path"),
			}
			resp := &validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}
			v.ValidateString(ctx, req, resp)
			require.False(t, resp.Diagnostics.HasError(), "Expected valid URL path, but got error: %s", resp.Diagnostics)
		})
	}

	invalidInputs := []string{
		"",
		"relative/path",
		"relative/path/",
		"//",
		"/api//mcp",
		"/api/v1//",
		"/path with unencoded spaces",
		"/path/with space/",
		"/path?query=1",
		"/path#fragment",
	}

	for _, input := range invalidInputs {
		t.Run("Invalid_"+input, func(t *testing.T) {
			req := validator.StringRequest{
				ConfigValue: types.StringValue(input),
				Path:        path.Root("url_path"),
			}
			resp := &validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}
			v.ValidateString(ctx, req, resp)
			require.True(t, resp.Diagnostics.HasError(), "Expected invalid URL path to fail validation, but it succeeded: %s", input)
		})
	}
}

func TestStringValidators_TableDriven(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type stringValidatorCase struct {
		name      string
		validator validator.String
		descPart  string
		valid     []string
		invalid   []string
	}

	cases := []stringValidatorCase{
		{
			name:      "NameLengthValidation",
			validator: NameLengthValidation(),
			descPart:  "string length must be between 1 and 128",
			valid: []string{
				"a",
				"my-resource-name",
				"12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678", // 128 chars
			},
			invalid: []string{
				"",
				"123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789", // 129 chars
			},
		},
		{
			name:      "UUIDRegexValidation",
			validator: UUIDRegexValidation(),
			descPart:  "must be a valid uuid",
			valid: []string{
				"550e8400-e29b-41d4-a716-446655440000",
				"{550e8400-e29b-41d4-a716-446655440000}",
			},
			invalid: []string{
				"not-a-uuid",
				"550e8400e29b41d4a716446655440000",
			},
		},
		{
			name:      "EmailValidation",
			validator: EmailValidation(),
			descPart:  "must be a valid email",
			valid: []string{
				"user@example.com",
				"user.name+tag@sub.example.co.uk",
			},
			invalid: []string{
				"not-an-email",
				"user@localhost",
			},
		},
		{
			name:      "SnowflakeAccountValidation",
			validator: SnowflakeAccountValidation(),
			descPart:  "must be a valid Snowflake Account ID",
			valid: []string{
				"acct123",
				"my-account-1",
			},
			invalid: []string{
				"acct_123",
				"acct.123",
			},
		},
		{
			name:      "SnowflakeUserNameValidation",
			validator: SnowflakeUserNameValidation(),
			descPart:  "must be a valid Snowflake Username",
			valid: []string{
				"svc_user",
				"user@example",
			},
			invalid: []string{
				"user-name",
				"user.name",
			},
		},
		{
			name:      "UrlSchemeValidation",
			validator: UrlSchemeValidation(),
			descPart:  "must be a valid http(s) scheme",
			valid: []string{
				"http://example.com",
				"https://example.com/path",
			},
			invalid: []string{
				"ftp://example.com",
				"example.com",
			},
		},
		{
			name:      "SecureURLValidation",
			validator: SecureURLValidation(),
			descPart:  "must be a valid https url",
			valid: []string{
				"https://example.com",
				"https://sub.example.com/path?q=1",
			},
			invalid: []string{
				"http://example.com",
				"https://localhost",
			},
		},
		{
			name:      "HostValidation",
			validator: HostValidation(),
			descPart:  "must be a valid hostname",
			valid: []string{
				"example.com",
				"api-1.internal",
			},
			invalid: []string{
				"-invalid.com",
				"invalid-.com",
			},
		},
		{
			name:      "SafeWildcardHostNameValidation",
			validator: SafeWildcardHostNameValidation(),
			descPart:  "must be a valid hostname or IP address",
			valid: []string{
				"192.168.1.1", // NOSONAR
				"*.example.com",
				"service.internal",
			},
			invalid: []string{
				"google.*",
				"*",
				"bad host",
			},
		},
		{
			name:      "S3BucketRegionLengthValidation",
			validator: S3BucketRegionLengthValidation(),
			descPart:  "string length must be between 8 and 30",
			valid: []string{
				"us-east1",
				"ap-southeast-2",
			},
			invalid: []string{
				"us-east",
				"region-name-exceeding-thirty-characters",
			},
		},
		{
			name:      "S3BucketRegionValidation",
			validator: S3BucketRegionValidation(),
			descPart:  "must be a valid AWS S3 Bucket region",
			valid: []string{
				"us-east-2",
				"eu-central-1",
			},
			invalid: []string{
				"US-EAST-2",
				"us_east_2",
			},
		},
		{
			name:      "S3BucketNameLengthValidation",
			validator: S3BucketNameLengthValidation(),
			descPart:  "string length must be between 3 and 63",
			valid: []string{
				"abc",
				"my-valid-bucket-name-123",
			},
			invalid: []string{
				"ab",
				"a123456789012345678901234567890123456789012345678901234567890123", // 64 chars
			},
		},
		{
			name:      "S3BucketNameValidation",
			validator: S3BucketNameValidation(),
			descPart:  "must be a valid AWS S3 Bucket name",
			valid: []string{
				"my-bucket.1",
				"abc",
			},
			invalid: []string{
				"-mybucket",
				"MyBucket",
			},
		},
		{
			name:      "S3PathPrefixLengthValidation",
			validator: S3PathPrefixLengthValidation(),
			descPart:  "string length must be between 3 and 800",
			valid: []string{
				"abc",
				"logs/prod/",
			},
			invalid: []string{
				"ab",
			},
		},
		{
			name:      "S3PathPrefixValidation",
			validator: S3PathPrefixValidation(),
			descPart:  "must be a valid AWS S3 path prefix",
			valid: []string{
				"logs/aembit/prod/",
				"prefix-1_2",
			},
			invalid: []string{
				"logs/{bad}",
				"logs/#hash",
			},
		},
		{
			name:      "GCSBucketNameSimpleValidation",
			validator: GCSBucketNameSimpleValidation(),
			descPart:  "must be a valid GCS Bucket name",
			valid: []string{
				"my_gcs-bucket1",
				"abc",
			},
			invalid: []string{
				"ab",
				"my.gcs.bucket",
			},
		},
		{
			name:      "GCSBucketNameWithPeriodValidation",
			validator: GCSBucketNameWithPeriodValidation(),
			descPart:  "must be a valid GCS Bucket name",
			valid: []string{
				"my.gcs_bucket-1",
				"abc",
			},
			invalid: []string{
				"ab",
				"MyBucket",
			},
		},
		{
			name:      "GCSPathPrefixValidation",
			validator: GCSPathPrefixValidation(),
			descPart:  "must be a valid GCS path prefix",
			valid: []string{
				"",
				"logs/aembit/",
			},
			invalid: []string{
				"logs/*wild",
				"logs/bad?prefix",
			},
		},
		{
			name:      "AudienceValidation",
			validator: AudienceValidation(),
			descPart:  "string length must be between 0 and 256",
			valid: []string{
				"",
				"https://api.example.com",
			},
			invalid: []string{
				string(make([]byte, 257)),
			},
		},
		{
			name:      "HecHostPortValidation",
			validator: HecHostPortValidation(),
			descPart:  "must be a valid HEC host port",
			valid: []string{
				"hec.example.com:8088",
				"10.0.0.1:443",
			},
			invalid: []string{
				"hec.example.com",
				"hec.example.com:8",
			},
		},
		{
			name:      "AuthenticationTokenValidation",
			validator: AuthenticationTokenValidation(),
			descPart:  "must be a valid authentication token",
			valid: []string{
				"550e8400-e29b-41d4-a716-446655440000",
			},
			invalid: []string{
				"{550e8400-e29b-41d4-a716-446655440000}",
				"invalid-token",
			},
		},
		{
			name:      "CrowdstrikeApiKeyValidation",
			validator: CrowdstrikeApiKeyValidation(),
			descPart:  "must be a valid API Key",
			valid: []string{
				"0123456789abcdef0123456789abcdef",
				"01234567-89ab-cdef-0123-456789abcdef",
			},
			invalid: []string{
				"short-key",
				"zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz",
			},
		},
		{
			name:      "Base64Validation",
			validator: Base64Validation(),
			descPart:  "must be a valid base64-encoded string",
			valid: []string{
				"YWVtYml0",
				"YWI=",
			},
			invalid: []string{
				"not-base64!",
				"abc",
			},
		},
		{
			name:      "AwsIamRoleArnValidation",
			validator: AwsIamRoleArnValidation(),
			descPart:  "must be a valid AWS ARN",
			valid: []string{
				"arn:aws:iam::123456789012:role/MyRole",
				"arn:aws:iam::123456789012:role/service-role/MyRole",
			},
			invalid: []string{
				"arn:aws:iam::12345:role/MyRole",
				"arn:aws:s3:::my-bucket",
			},
		},
		{
			name:      "AwsSecretArnValidation",
			validator: AwsSecretArnValidation(),
			descPart:  "must be a valid AWS Secrets Manager secret ARN",
			valid: []string{
				"arn:aws:secretsmanager:us-east-2:123456789012:secret:my-secret-AbCdEf",
			},
			invalid: []string{
				"arn:aws:secretsmanager:us-east-2:123456789012:secret:nosuffix",
				"not-an-arn",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			require.NotNil(t, tc.validator)
			require.Contains(t, tc.validator.Description(ctx), tc.descPart)
			require.Contains(t, tc.validator.MarkdownDescription(ctx), tc.descPart)

			for _, val := range tc.valid {
				t.Run("Valid_"+val, func(t *testing.T) {
					req := validator.StringRequest{
						ConfigValue: types.StringValue(val),
						Path:        path.Root("test_field"),
					}
					resp := &validator.StringResponse{
						Diagnostics: diag.Diagnostics{},
					}
					tc.validator.ValidateString(ctx, req, resp)
					require.False(t, resp.Diagnostics.HasError(), "expected %q to pass %s, got: %s", val, tc.name, resp.Diagnostics)
				})
			}

			for _, val := range tc.invalid {
				t.Run("Invalid_"+val, func(t *testing.T) {
					req := validator.StringRequest{
						ConfigValue: types.StringValue(val),
						Path:        path.Root("test_field"),
					}
					resp := &validator.StringResponse{
						Diagnostics: diag.Diagnostics{},
					}
					tc.validator.ValidateString(ctx, req, resp)
					require.True(t, resp.Diagnostics.HasError(), "expected %q to fail %s", val, tc.name)
				})
			}

			for _, skipped := range []struct {
				label string
				val   types.String
			}{
				{label: "Null", val: types.StringNull()},
				{label: "Unknown", val: types.StringUnknown()},
			} {
				t.Run(skipped.label, func(t *testing.T) {
					req := validator.StringRequest{
						ConfigValue: skipped.val,
						Path:        path.Root("test_field"),
					}
					resp := &validator.StringResponse{
						Diagnostics: diag.Diagnostics{},
					}
					tc.validator.ValidateString(ctx, req, resp)
					require.False(t, resp.Diagnostics.HasError(), "expected %s value to skip validation in %s", skipped.label, tc.name)
				})
			}
		})
	}
}

func TestOptionalSecureURLValidation(t *testing.T) {
	t.Parallel()

	v := OptionalSecureURLValidation()
	require.NotNil(t, v)

	ctx := context.Background()

	require.Contains(t, v.Description(ctx), "must be a valid https url")
	require.Contains(t, v.Description(ctx), "or be empty")
	require.Contains(t, v.MarkdownDescription(ctx), "must be a valid https url")
	require.Contains(t, v.MarkdownDescription(ctx), "or be empty")

	validCases := []struct {
		name  string
		value types.String
	}{
		{name: "Null", value: types.StringNull()},
		{name: "Unknown", value: types.StringUnknown()},
		{name: "Empty", value: types.StringValue("")},
		{name: "ValidHTTPS", value: types.StringValue("https://example.com")},
		{name: "ValidHTTPSWithPathAndQuery", value: types.StringValue("https://sub.example.com/oauth/revoke?hint=token")},
	}

	for _, tc := range validCases {
		t.Run("Valid_"+tc.name, func(t *testing.T) {
			req := validator.StringRequest{
				ConfigValue: tc.value,
				Path:        path.Root("oauth_revocation_url"),
			}
			resp := &validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}
			v.ValidateString(ctx, req, resp)
			require.False(t, resp.Diagnostics.HasError(), "expected %v to be valid, got: %s", tc.value, resp.Diagnostics)
		})
	}

	invalidCases := []struct {
		name  string
		value string
	}{
		{name: "HTTP", value: "http://example.com/revoke"},
		{name: "WhitespaceOnly", value: "   "},
		{name: "MalformedURL", value: "not-a-valid-url"},
		{name: "LocalhostWithoutTLD", value: "https://localhost/revoke"},
		{name: "FTPScheme", value: "ftp://example.com/revoke"},
	}

	for _, tc := range invalidCases {
		t.Run("Invalid_"+tc.name, func(t *testing.T) {
			req := validator.StringRequest{
				ConfigValue: types.StringValue(tc.value),
				Path:        path.Root("oauth_revocation_url"),
			}
			resp := &validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}
			v.ValidateString(ctx, req, resp)
			require.True(t, resp.Diagnostics.HasError(), "expected %q to fail validation", tc.value)
			errs := resp.Diagnostics.Errors()
			require.Len(t, errs, 1, "expected exactly 1 error diagnostic for %q", tc.value)
			require.Contains(t, errs[0].Detail(), "must be a valid https url")
		})
	}
}
