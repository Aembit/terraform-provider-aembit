package validators

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var UUIDRegex = regexp.MustCompile(`^[{]?[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}[}]?$`)
var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)
var SnowflakeAccountRegex = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
var SnowflakeUserNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_@]+$`)
var UrlSchemeRegex = regexp.MustCompile(`^http(s)?:\/\/.*$`)
var SecureURLRegex = regexp.MustCompile(`^https:\/\/[a-zA-Z0-9\-.]+\.[a-zA-Z]{2,}(\/\S*)?$`)
var HostNameRegex = regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)*(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)$`)
var S3BucketRegionRegex = regexp.MustCompile(`^[a-z\-\d]+$`)
var S3BucketNameRegex = regexp.MustCompile(`^[a-z\d][a-z\-\.\d]{1,61}?[a-z\d]$`)
var S3PathPrefixRegex = regexp.MustCompile(`^[^\\\{\}\^\%` + "`" + `""'~#\[\]\>\<\|\x80-\xff]*$`)
var GCSBucketNameSimpleRegex = regexp.MustCompile(`^[a-z0-9-_]{3,63}$`)
var GCSBucketNameWithPeriodRegex = regexp.MustCompile(`^[a-z0-9-\._]{3,222}$`)
var GCSPathPrefixRegex = regexp.MustCompile(`^[^#\[\]*?:\"<>|]{0,256}$`)
var HecHostPortRegex = regexp.MustCompile(`^([a-zA-Z0-9.-]+):(\d{2,5})$`)
var AuthenticationTokenRegex = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)

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
	return stringvalidator.RegexMatches(SnowflakeAccountRegex, "must be a valid Snowflake Account ID")
}

func SnowflakeUserNameValidation() validator.String {
	return stringvalidator.RegexMatches(SnowflakeUserNameRegex, "must be a valid Snowflake Username")
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
	return stringvalidator.RegexMatches(GCSBucketNameWithPeriodRegex, "must be a valid GCS Bucket name")
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
	return stringvalidator.RegexMatches(AuthenticationTokenRegex, "must be a valid authentication token")
}
