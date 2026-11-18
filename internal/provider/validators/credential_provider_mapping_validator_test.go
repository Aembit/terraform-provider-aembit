package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func mappingObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"credential_provider_id": types.StringType,
			"mapping_type":           types.StringType,
			"account_name":           types.StringType,
			"access_key_id":          types.StringType,
			"header_name":            types.StringType,
			"header_value":           types.StringType,
			"httpbody_field_path":    types.StringType,
			"httpbody_field_value":   types.StringType,
		},
	}
}

func makeMappingObject(
	cpID, mappingType, accountName, accessKeyID, headerName, headerValue, bodyPath, bodyValue string,
) types.Object {
	obj, diags := types.ObjectValue(
		mappingObjectType().AttrTypes,
		map[string]attr.Value{
			"credential_provider_id": types.StringValue(cpID),
			"mapping_type":           types.StringValue(mappingType),
			"account_name":           types.StringValue(accountName),
			"access_key_id":          types.StringValue(accessKeyID),
			"header_name":            types.StringValue(headerName),
			"header_value":           types.StringValue(headerValue),
			"httpbody_field_path":    types.StringValue(bodyPath),
			"httpbody_field_value":   types.StringValue(bodyValue),
		},
	)
	if diags.HasError() {
		panic("failed to build mapping object: " + diags.Errors()[0].Detail())
	}
	return obj
}

func TestCredentialProviderMappingValidator(t *testing.T) {
	t.Parallel()

	v := NewCredentialProviderMappingValidator()
	require.NotNil(t, v)

	ctx := context.Background()
	require.Equal(t, "Ensure there are no duplicate mapping values", v.Description(ctx))
	require.Equal(t, v.Description(ctx), v.MarkdownDescription(ctx))

	testCases := []struct {
		name        string
		configValue types.Set
		expectError bool
		errContains string
	}{
		{
			name:        "null set",
			configValue: types.SetNull(mappingObjectType()),
			expectError: false,
		},
		{
			name:        "unknown set",
			configValue: types.SetUnknown(mappingObjectType()),
			expectError: false,
		},
		{
			name:        "empty set",
			configValue: types.SetValueMust(mappingObjectType(), []attr.Value{}),
			expectError: false,
		},
		{
			name: "distinct mappings pass",
			configValue: types.SetValueMust(mappingObjectType(), []attr.Value{
				makeMappingObject("cp-1", "HttpHeader", "", "", "X-Tenant", "alpha", "", ""),
				makeMappingObject("cp-2", "HttpHeader", "", "", "X-Tenant", "beta", "", ""),
				makeMappingObject("cp-3", "SnowflakeAccount", "acct-1", "", "", "", "", ""),
				makeMappingObject("cp-4", "AwsAccessKeyId", "", "AKIA123", "", "", "", ""),
				makeMappingObject("cp-5", "HttpBody", "", "", "", "", "$.role", "admin"),
			}),
			expectError: false,
		},
		{
			name: "distinct mappings with adjacent concatenated fields do not collide",
			configValue: types.SetValueMust(mappingObjectType(), []attr.Value{
				makeMappingObject("cp-1", "HttpHeader", "ab", "", "c", "", "", ""),
				makeMappingObject("cp-2", "HttpHeader", "a", "", "bc", "", "", ""),
			}),
			expectError: false,
		},
		{
			name: "duplicate header mapping fails",
			configValue: types.SetValueMust(mappingObjectType(), []attr.Value{
				makeMappingObject("cp-1", "HttpHeader", "", "", "X-Tenant", "alpha", "", ""),
				makeMappingObject("cp-2", "HttpHeader", "", "", "X-Tenant", "alpha", "", ""),
			}),
			expectError: true,
			errContains: "duplicate credential provider mapping already exists",
		},
		{
			name: "duplicate account name mapping fails",
			configValue: types.SetValueMust(mappingObjectType(), []attr.Value{
				makeMappingObject("cp-1", "SnowflakeAccount", "acct-1", "", "", "", "", ""),
				makeMappingObject("cp-2", "SnowflakeAccount", "acct-1", "", "", "", "", ""),
			}),
			expectError: true,
			errContains: "duplicate credential provider mapping already exists",
		},
		{
			name: "mismatched element type triggers ElementsAs diagnostic error",
			configValue: types.SetValueMust(types.StringType, []attr.Value{
				types.StringValue("not-an-object"),
			}),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := validator.SetRequest{
				Path:        path.Root("credential_providers"),
				ConfigValue: tc.configValue,
			}
			resp := &validator.SetResponse{
				Diagnostics: diag.Diagnostics{},
			}

			v.ValidateSet(ctx, req, resp)

			if tc.expectError {
				require.True(t, resp.Diagnostics.HasError(), "expected validation error")
				if tc.errContains != "" {
					require.Contains(t, resp.Diagnostics.Errors()[0].Detail(), tc.errContains)
				}
			} else {
				require.False(t, resp.Diagnostics.HasError(), "unexpected validation error: %v", resp.Diagnostics)
			}
		})
	}
}
