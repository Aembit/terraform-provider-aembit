package validators

import (
	"context"
	"fmt"

	"terraform-provider-aembit/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// mappingUniquenessKey uniquely identifies a credential provider mapping entry
// across its individual mapping fields without concatenation collisions.
type mappingUniquenessKey struct {
	// accountName is the Snowflake account name for the mapping.
	accountName string
	// headerName is the HTTP header name for the mapping.
	headerName string
	// headerValue is the HTTP header value for the mapping.
	headerValue string
	// httpbodyFieldPath is the HTTP body field path for the mapping.
	httpbodyFieldPath string
	// httpbodyFieldValue is the HTTP body field value for the mapping.
	httpbodyFieldValue string
	// accessKeyID is the AWS access key ID for the mapping.
	accessKeyID string
}

type CredentialProviderMappingValidator struct{}

func (v CredentialProviderMappingValidator) Description(_ context.Context) string {
	return "Ensure there are no duplicate mapping values"
}

func (v CredentialProviderMappingValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v CredentialProviderMappingValidator) ValidateSet(
	ctx context.Context,
	req validator.SetRequest,
	resp *validator.SetResponse,
) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var credentialProviders []models.PolicyCredentialMappingModel
	diags := req.ConfigValue.ElementsAs(ctx, &credentialProviders, false)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	uniqueMap := make(map[mappingUniquenessKey]bool)
	for _, cp := range credentialProviders {
		key := mappingUniquenessKey{
			accountName:        cp.AccountName.ValueString(),
			headerName:         cp.HeaderName.ValueString(),
			headerValue:        cp.HeaderValue.ValueString(),
			httpbodyFieldPath:  cp.HttpbodyFieldPath.ValueString(),
			httpbodyFieldValue: cp.HttpbodyFieldValue.ValueString(),
			accessKeyID:        cp.AccessKeyId.ValueString(),
		}
		_, found := uniqueMap[key]

		if found {
			resp.Diagnostics.AddError(
				"Error validating access policy credential provider mapping",
				fmt.Errorf("duplicate credential provider mapping already exists").Error(),
			)
			return
		}
		uniqueMap[key] = true
	}
}

func NewCredentialProviderMappingValidator() validator.Set {
	return CredentialProviderMappingValidator{}
}
