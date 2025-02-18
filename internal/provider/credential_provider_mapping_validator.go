package provider

import (
	"context"
	"fmt"

	"aembit.io/aembit"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type CredentialProviderMappingValidator struct{}

func (v CredentialProviderMappingValidator) Description(ctx context.Context) string {
	return "Ensures there is no duplicate in mapping values"
}

func (v CredentialProviderMappingValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v CredentialProviderMappingValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var plan accessPolicyResourceModel
	diags := req.Config.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var policy aembit.CreatePolicyDTO = convertAccessPolicyModelToPolicyDTO(plan, nil)

	uniqueMap := make(map[string]bool)
	for _, cp := range policy.CredentialProviders {
		mapValue := cp.AccountName + cp.HeaderName + cp.HeaderValue + cp.HttpbodyFieldPath + cp.HttpbodyFieldValue
		_, found := uniqueMap[mapValue]

		if found {
			resp.Diagnostics.AddError(
				"Error validating access policy",
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
