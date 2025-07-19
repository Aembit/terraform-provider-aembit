package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"terraform-provider-aembit/internal/provider/models"
)

type CredentialProviderMappingValidator struct{}

func (v CredentialProviderMappingValidator) Description(ctx context.Context) string {
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

	uniqueMap := make(map[string]bool)
	for _, cp := range credentialProviders {
		mapValue := cp.AccountName.ValueString() + cp.HeaderName.ValueString() + cp.HeaderValue.ValueString() + cp.HttpbodyFieldPath.ValueString() + cp.HttpbodyFieldValue.ValueString()
		_, found := uniqueMap[mapValue]

		if found {
			resp.Diagnostics.AddError(
				"Error validating access policy credential provider mapping",
				fmt.Errorf("duplicate credential provider mapping already exists").Error(),
			)
			return
		}
		uniqueMap[mapValue] = true
	}
}

func NewCredentialProviderMappingValidator() validator.Set {
	return CredentialProviderMappingValidator{}
}
