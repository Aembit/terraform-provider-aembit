package validators

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	UUIDRegex = regexp.MustCompile(
		`^[{]?[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}[}]?$`,
	)
	EmailRegex = regexp.MustCompile(
		`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`,
	)
	SnowflakeAccountRegex  = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
	SnowflakeUserNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_@]+$`)
	UrlSchemeRegex         = regexp.MustCompile(`^http(s)?:\/\/.*$`)
	SecureURLRegex         = regexp.MustCompile(
		`^https:\/\/[a-zA-Z0-9\-.]+\.[a-zA-Z]{2,}(\/\S*)?$`,
	)
	HostNameRegex = regexp.MustCompile(
		`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)*(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)$`,
	)
	HostIPRegex = regexp.MustCompile(
		`^((?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)\.){3}(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)$`,
	)
	StructuralHostRegex = regexp.MustCompile(
		`^(?:(?:[a-zA-Z\d*](?:[a-zA-Z\d*-]*[a-zA-Z\d*])?\.)*(?:[a-zA-Z\d](?:[a-zA-Z\d-]*[a-zA-Z\d])?)\.(?:[a-zA-Z\d](?:[a-zA-Z\d-]*[a-zA-Z\d])?)|[a-zA-Z\d](?:[a-zA-Z\d-]*[a-zA-Z\d])?)$`,
	)
	S3BucketRegionRegex = regexp.MustCompile(`^[a-z\-\d]+$`)
	S3BucketNameRegex   = regexp.MustCompile(`^[a-z\d][a-z\-\.\d]{1,61}?[a-z\d]$`)
	S3PathPrefixRegex   = regexp.MustCompile(
		`^[^\\\{\}\^\%` + "`" + `""'~#\[\]\>\<\|\x80-\xff]*$`,
	)
	GCSBucketNameSimpleRegex     = regexp.MustCompile(`^[a-z0-9-_]{3,63}$`)
	GCSBucketNameWithPeriodRegex = regexp.MustCompile(`^[a-z0-9-\._]{3,222}$`)
	GCSPathPrefixRegex           = regexp.MustCompile(`^[^#\[\]*?:\"<>|]{0,256}$`)
	HecHostPortRegex             = regexp.MustCompile(`^([a-zA-Z0-9.-]+):(\d{2,5})$`)
	AuthenticationTokenRegex     = regexp.MustCompile(
		`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`,
	)
	ApiKeyRegex = regexp.MustCompile(`^[a-fA-F0-9-]{32,40}$`)
	Base64Regex = regexp.MustCompile(
		`^(?:[A-Za-z0-9+/]{4})*` +
			`(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$`,
	)
	AwsIamRoleArnRegex = regexp.MustCompile(`^arn:aws:iam::\d{12}:role/[\w+=,.@/-]+$`)
	AwsSecretArnRegex  = regexp.MustCompile(
		`^arn:aws:secretsmanager:[a-z0-9-]+:\d{12}:secret:[A-Za-z0-9/_+=.@-]+-[A-Za-z0-9]{6}$`,
	)
	SpiffeRegex = regexp.MustCompile(
		`^spiffe:\/\/`,
	)
)

func NameLengthValidation() validator.String {
	return stringvalidator.LengthBetween(1, 128)
}

func UUIDRegexValidation() validator.String {
	return stringvalidator.RegexMatches(UUIDRegex, "must be a valid uuid")
}

func EmailValidation() validator.String {
	return stringvalidator.RegexMatches(EmailRegex, "must be a valid email")
}

func SnowflakeAccountValidation() validator.String {
	return stringvalidator.RegexMatches(
		SnowflakeAccountRegex,
		"must be a valid Snowflake Account ID",
	)
}

func SnowflakeUserNameValidation() validator.String {
	return stringvalidator.RegexMatches(
		SnowflakeUserNameRegex,
		"must be a valid Snowflake Username",
	)
}

func UrlSchemeValidation() validator.String {
	return stringvalidator.RegexMatches(UrlSchemeRegex, "must be a valid http(s) scheme")
}

func SecureURLValidation() validator.String {
	return stringvalidator.RegexMatches(SecureURLRegex, "must be a valid https url")
}

func HostValidation() validator.String {
	return stringvalidator.RegexMatches(HostNameRegex, "must be a valid hostname")
}

func SafeWildcardHostNameValidation() validator.String {
	return stringvalidator.Any(
		stringvalidator.RegexMatches(HostIPRegex, "must be a valid hostname or IP address"),
		stringvalidator.RegexMatches(StructuralHostRegex, "must be a valid hostname or IP address"),
	)
}

func S3BucketRegionLengthValidation() validator.String {
	return stringvalidator.LengthBetween(8, 30)
}

func S3BucketRegionValidation() validator.String {
	return stringvalidator.RegexMatches(S3BucketRegionRegex, "must be a valid AWS S3 Bucket region")
}

func S3BucketNameLengthValidation() validator.String {
	return stringvalidator.LengthBetween(3, 63)
}

func S3BucketNameValidation() validator.String {
	return stringvalidator.RegexMatches(S3BucketNameRegex, "must be a valid AWS S3 Bucket name")
}

func S3PathPrefixLengthValidation() validator.String {
	return stringvalidator.LengthBetween(3, 800)
}

func S3PathPrefixValidation() validator.String {
	return stringvalidator.RegexMatches(S3PathPrefixRegex, "must be a valid AWS S3 path prefix")
}

func GCSBucketNameSimpleValidation() validator.String {
	return stringvalidator.RegexMatches(GCSBucketNameSimpleRegex, "must be a valid GCS Bucket name")
}

func GCSBucketNameWithPeriodValidation() validator.String {
	return stringvalidator.RegexMatches(
		GCSBucketNameWithPeriodRegex,
		"must be a valid GCS Bucket name",
	)
}

func GCSPathPrefixValidation() validator.String {
	return stringvalidator.RegexMatches(GCSPathPrefixRegex, "must be a valid GCS path prefix")
}

func AudienceValidation() validator.String {
	return stringvalidator.LengthBetween(0, 256)
}

func HecHostPortValidation() validator.String {
	return stringvalidator.RegexMatches(HecHostPortRegex, "must be a valid HEC host port")
}

func AuthenticationTokenValidation() validator.String {
	return stringvalidator.RegexMatches(
		AuthenticationTokenRegex,
		"must be a valid authentication token",
	)
}

func CrowdstrikeApiKeyValidation() validator.String {
	return stringvalidator.RegexMatches(ApiKeyRegex, "must be a valid API Key")
}

func Base64Validation() validator.String {
	return stringvalidator.RegexMatches(
		Base64Regex,
		"must be a valid base64-encoded string",
	)
}

func AwsIamRoleArnValidation() validator.String {
	return stringvalidator.RegexMatches(AwsIamRoleArnRegex, "must be a valid AWS ARN")
}

func AwsSecretArnValidation() validator.String {
	return stringvalidator.RegexMatches(
		AwsSecretArnRegex,
		"must be a valid AWS Secrets Manager secret ARN",
	)
}

func SpiffeSubjectValidation() validator.String {
	return stringvalidator.RegexMatches(
		SpiffeRegex,
		"must be in the format: spiffe://trust-domain-name/path",
	)
}
