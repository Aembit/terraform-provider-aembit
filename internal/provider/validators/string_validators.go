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
