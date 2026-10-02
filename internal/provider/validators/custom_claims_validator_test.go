package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createClaimObjectType returns the object type specification for a custom claim.
func createClaimObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"key":        types.StringType,
			"value":      types.StringType,
			"value_type": types.StringType,
		},
	}
}

// createClaimObject creates a custom claim object value from key, value, and value_type strings.
func createClaimObject(key, value, valueType string) types.Object {
	obj, diags := types.ObjectValue(
		createClaimObjectType().AttrTypes,
		map[string]attr.Value{
			"key":        types.StringValue(key),
			"value":      types.StringValue(value),
			"value_type": types.StringValue(valueType),
		},
	)
	if diags.HasError() {
		panic("failed to construct claim object value: " + diags.Errors()[0].Detail())
	}
	return obj
}

// createClaimObjectWithAttrs creates a custom claim object value from explicit attr.Value instances.
func createClaimObjectWithAttrs(key, value, valueType attr.Value) types.Object {
	obj, diags := types.ObjectValue(
		createClaimObjectType().AttrTypes,
		map[string]attr.Value{
			"key":        key,
			"value":      value,
			"value_type": valueType,
		},
	)
	if diags.HasError() {
		panic("failed to construct custom claim object value with raw attrs: " + diags.Errors()[0].Detail())
	}
	return obj
}

// TestCustomClaimsValidator_ValidDynamicExpressions verifies that valid dynamic template expressions pass validation.
func TestCustomClaimsValidator_ValidDynamicExpressions(t *testing.T) {
	t.Parallel()

	validExpressions := []string{
		"${client.executable.hash.sha256}",
		"${client.id}",
		"prefix-${client.executable.hash.sha256}-suffix",
		"${aembit:token['sub']}",
		"${client.path}/v1/${client.id}",
		"literal_prefix_${client.os}",
	}

	for _, expr := range validExpressions {
		expr := expr
		t.Run(expr, func(t *testing.T) {
			t.Parallel()
			v := NewCustomClaimsValidator()
			ctx := context.Background()

			claimObj := createClaimObject("workload_hash", expr, "dynamic")
			setVal, diags := types.SetValue(createClaimObjectType(), []attr.Value{claimObj})
			require.False(t, diags.HasError())

			req := validator.SetRequest{
				Path:        path.Root("custom_claims"),
				ConfigValue: setVal,
			}
			resp := &validator.SetResponse{
				Diagnostics: diag.Diagnostics{},
			}

			v.ValidateSet(ctx, req, resp)
			assert.False(t, resp.Diagnostics.HasError(), "expected valid dynamic expression to pass: %s", expr)
		})
	}
}

// TestCustomClaimsValidator_ValidLiteralStrings verifies that valid literal strings pass validation.
func TestCustomClaimsValidator_ValidLiteralStrings(t *testing.T) {
	t.Parallel()

	validLiterals := []string{
		"my-literal-value",
		"production",
		"static-guid-12345",
		"plain text with spaces",
		"https://example.com/api",
		"$not_a_template",
		"{still_not_template}",
	}

	for _, lit := range validLiterals {
		lit := lit
		t.Run(lit, func(t *testing.T) {
			t.Parallel()
			v := NewCustomClaimsValidator()
			ctx := context.Background()

			claimObj := createClaimObject("environment", lit, "literal")
			setVal, diags := types.SetValue(createClaimObjectType(), []attr.Value{claimObj})
			require.False(t, diags.HasError())

			req := validator.SetRequest{
				Path:        path.Root("custom_claims"),
				ConfigValue: setVal,
			}
			resp := &validator.SetResponse{
				Diagnostics: diag.Diagnostics{},
			}

			v.ValidateSet(ctx, req, resp)
			assert.False(t, resp.Diagnostics.HasError(), "expected valid literal string to pass: %s", lit)
		})
	}
}

// TestCustomClaimsValidator_InvalidDynamic_MissingTemplate verifies error when dynamic claim lacks template syntax.
func TestCustomClaimsValidator_InvalidDynamic_MissingTemplate(t *testing.T) {
	t.Parallel()

	v := NewCustomClaimsValidator()
	ctx := context.Background()

	claimObj := createClaimObject("test_claim", "value2", "dynamic")
	setVal, diags := types.SetValue(createClaimObjectType(), []attr.Value{claimObj})
	require.False(t, diags.HasError())

	req := validator.SetRequest{
		Path:        path.Root("custom_claims"),
		ConfigValue: setVal,
	}
	resp := &validator.SetResponse{
		Diagnostics: diag.Diagnostics{},
	}

	v.ValidateSet(ctx, req, resp)
	require.True(t, resp.Diagnostics.HasError())
	assert.Contains(
		t,
		resp.Diagnostics[0].Detail(),
		"Field 'CustomClaim 'test_claim'' must contain a valid template expression '${...}' when configured as dynamic.",
	)
}

// TestCustomClaimsValidator_InvalidDynamic_EmptyString verifies error when dynamic claim value is empty.
func TestCustomClaimsValidator_InvalidDynamic_EmptyString(t *testing.T) {
	t.Parallel()

	v := NewCustomClaimsValidator()
	ctx := context.Background()

	claimObj := createClaimObject("empty_claim", "", "dynamic")
	setVal, diags := types.SetValue(createClaimObjectType(), []attr.Value{claimObj})
	require.False(t, diags.HasError())

	req := validator.SetRequest{
		Path:        path.Root("custom_claims"),
		ConfigValue: setVal,
	}
	resp := &validator.SetResponse{
		Diagnostics: diag.Diagnostics{},
	}

	v.ValidateSet(ctx, req, resp)
	require.True(t, resp.Diagnostics.HasError())
	assert.Contains(
		t,
		resp.Diagnostics[0].Detail(),
		"Field 'CustomClaim 'empty_claim'' must contain a valid template expression '${...}' when configured as dynamic.",
	)
}

// TestCustomClaimsValidator_InvalidDynamic_MalformedSyntax verifies error when dynamic claim has malformed template syntax.
func TestCustomClaimsValidator_InvalidDynamic_MalformedSyntax(t *testing.T) {
	t.Parallel()

	malformedCases := []struct {
		name  string
		value string
	}{
		{
			name:  "spaces inside expression",
			value: "${invalid with spaces}",
		},
		{
			name:  "unclosed expression",
			value: "${unclosed",
		},
		{
			name:  "empty brackets",
			value: "${}",
		},
		{
			name:  "valid followed by unclosed",
			value: "${client.id} and ${unclosed",
		},
		{
			name:  "unclosed with trailing braces",
			value: "${client.id",
		},
	}

	for _, tc := range malformedCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			v := NewCustomClaimsValidator()
			ctx := context.Background()

			claimObj := createClaimObject("bad_claim", tc.value, "dynamic")
			setVal, diags := types.SetValue(createClaimObjectType(), []attr.Value{claimObj})
			require.False(t, diags.HasError())

			req := validator.SetRequest{
				Path:        path.Root("custom_claims"),
				ConfigValue: setVal,
			}
			resp := &validator.SetResponse{
				Diagnostics: diag.Diagnostics{},
			}

			v.ValidateSet(ctx, req, resp)
			require.True(t, resp.Diagnostics.HasError(), "expected error for malformed case: %s", tc.name)
			assert.Contains(
				t,
				resp.Diagnostics[0].Detail(),
				"Field 'CustomClaim 'bad_claim'' has invalid template expression syntax.",
			)
		})
	}
}

// TestCustomClaimsValidator_InvalidLiteral_ContainsTemplate verifies error when literal claim contains template syntax.
func TestCustomClaimsValidator_InvalidLiteral_ContainsTemplate(t *testing.T) {
	t.Parallel()

	forbiddenLiterals := []struct {
		name  string
		value string
	}{
		{
			name:  "prefix with dynamic variable",
			value: "literal-${foo}",
		},
		{
			name:  "exact dynamic variable",
			value: "${client.executable.hash.sha256}",
		},
		{
			name:  "empty dynamic envelope in literal",
			value: "literal-${}",
		},
		{
			name:  "embedded dynamic variable",
			value: "val_${test}_suffix",
		},
	}

	for _, tc := range forbiddenLiterals {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			v := NewCustomClaimsValidator()
			ctx := context.Background()

			claimObj := createClaimObject("literal_claim", tc.value, "literal")
			setVal, diags := types.SetValue(createClaimObjectType(), []attr.Value{claimObj})
			require.False(t, diags.HasError())

			req := validator.SetRequest{
				Path:        path.Root("custom_claims"),
				ConfigValue: setVal,
			}
			resp := &validator.SetResponse{
				Diagnostics: diag.Diagnostics{},
			}

			v.ValidateSet(ctx, req, resp)
			require.True(t, resp.Diagnostics.HasError(), "expected error for literal containing template: %s", tc.name)
			assert.Contains(
				t,
				resp.Diagnostics[0].Detail(),
				"Field 'CustomClaim 'literal_claim'' cannot contain dynamic template syntax '${...}' when configured as literal.",
			)
		})
	}
}

// TestCustomClaimsValidator_CaseInsensitivity verifies that value_type matching is case-insensitive.
func TestCustomClaimsValidator_CaseInsensitivity(t *testing.T) {
	t.Parallel()

	v := NewCustomClaimsValidator()
	ctx := context.Background()

	// Dynamic in uppercase
	dynamicObj := createClaimObject("dyn_claim", "${client.id}", "DYNAMIC")
	// Literal in uppercase with invalid dynamic template
	literalObj := createClaimObject("lit_claim", "foo-${bar}", "LITERAL")

	setVal, diags := types.SetValue(createClaimObjectType(), []attr.Value{dynamicObj, literalObj})
	require.False(t, diags.HasError())

	req := validator.SetRequest{
		Path:        path.Root("custom_claims"),
		ConfigValue: setVal,
	}
	resp := &validator.SetResponse{
		Diagnostics: diag.Diagnostics{},
	}

	v.ValidateSet(ctx, req, resp)
	require.True(t, resp.Diagnostics.HasError())
	assert.Contains(
		t,
		resp.Diagnostics[0].Detail(),
		"Field 'CustomClaim 'lit_claim'' cannot contain dynamic template syntax '${...}' when configured as literal.",
	)
}

// TestCustomClaimsValidator_EmptyKeyFieldName verifies diagnostic formatting when claim key is empty.
func TestCustomClaimsValidator_EmptyKeyFieldName(t *testing.T) {
	t.Parallel()

	v := NewCustomClaimsValidator()
	ctx := context.Background()

	claimObj := createClaimObject("", "literal-${foo}", "literal")
	setVal, diags := types.SetValue(createClaimObjectType(), []attr.Value{claimObj})
	require.False(t, diags.HasError())

	req := validator.SetRequest{
		Path:        path.Root("custom_claims"),
		ConfigValue: setVal,
	}
	resp := &validator.SetResponse{
		Diagnostics: diag.Diagnostics{},
	}

	v.ValidateSet(ctx, req, resp)
	require.True(t, resp.Diagnostics.HasError())
	assert.Contains(
		t,
		resp.Diagnostics[0].Detail(),
		"Field 'CustomClaim' cannot contain dynamic template syntax '${...}' when configured as literal.",
	)
}

// TestCustomClaimsValidator_NullAndUnknownValues verifies that null and unknown values safely pass validation.
func TestCustomClaimsValidator_NullAndUnknownValues(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		setValue types.Set
	}{
		{
			name:     "null set value",
			setValue: types.SetNull(createClaimObjectType()),
		},
		{
			name:     "unknown set value",
			setValue: types.SetUnknown(createClaimObjectType()),
		},
		{
			name: "set with unknown object element",
			setValue: func() types.Set {
				s, diags := types.SetValue(
					createClaimObjectType(),
					[]attr.Value{types.ObjectUnknown(createClaimObjectType().AttrTypes)},
				)
				require.False(t, diags.HasError())
				return s
			}(),
		},
		{
			name: "set with unknown claim value",
			setValue: func() types.Set {
				claim := createClaimObjectWithAttrs(
					types.StringValue("key1"),
					types.StringUnknown(),
					types.StringValue("dynamic"),
				)
				s, diags := types.SetValue(createClaimObjectType(), []attr.Value{claim})
				require.False(t, diags.HasError())
				return s
			}(),
		},
		{
			name: "set with unknown value_type",
			setValue: func() types.Set {
				claim := createClaimObjectWithAttrs(
					types.StringValue("key1"),
					types.StringValue("literal-${val}"),
					types.StringUnknown(),
				)
				s, diags := types.SetValue(createClaimObjectType(), []attr.Value{claim})
				require.False(t, diags.HasError())
				return s
			}(),
		},
		{
			name: "set with null claim value",
			setValue: func() types.Set {
				claim := createClaimObjectWithAttrs(
					types.StringValue("key1"),
					types.StringNull(),
					types.StringValue("dynamic"),
				)
				s, diags := types.SetValue(createClaimObjectType(), []attr.Value{claim})
				require.False(t, diags.HasError())
				return s
			}(),
		},
		{
			name: "set with null value_type",
			setValue: func() types.Set {
				claim := createClaimObjectWithAttrs(
					types.StringValue("key1"),
					types.StringValue("literal-${val}"),
					types.StringNull(),
				)
				s, diags := types.SetValue(createClaimObjectType(), []attr.Value{claim})
				require.False(t, diags.HasError())
				return s
			}(),
		},
		{
			name: "set with unknown claim key but valid dynamic value",
			setValue: func() types.Set {
				claim := createClaimObjectWithAttrs(
					types.StringUnknown(),
					types.StringValue("${client.executable.hash.sha256}"),
					types.StringValue("dynamic"),
				)
				s, diags := types.SetValue(createClaimObjectType(), []attr.Value{claim})
				require.False(t, diags.HasError())
				return s
			}(),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			v := NewCustomClaimsValidator()
			ctx := context.Background()

			req := validator.SetRequest{
				Path:        path.Root("custom_claims"),
				ConfigValue: tc.setValue,
			}
			resp := &validator.SetResponse{
				Diagnostics: diag.Diagnostics{},
			}

			v.ValidateSet(ctx, req, resp)
			assert.False(t, resp.Diagnostics.HasError(), "expected null/unknown test to pass: %s", tc.name)
		})
	}
}

// TestCustomClaimsValidator_Descriptions verifies that validator descriptions are populated.
func TestCustomClaimsValidator_Descriptions(t *testing.T) {
	t.Parallel()

	v := NewCustomClaimsValidator()
	ctx := context.Background()

	assert.NotEmpty(t, v.Description(ctx))
	assert.Equal(t, v.Description(ctx), v.MarkdownDescription(ctx))
}

// TestCustomClaimsValidator_UninitializedStructDefaults verifies that an uninitialized struct compiles regexes lazily.
func TestCustomClaimsValidator_UninitializedStructDefaults(t *testing.T) {
	t.Parallel()

	v := CustomClaimsValidator{}
	ctx := context.Background()

	claimObj := createClaimObject("workload_hash", "${client.executable.hash.sha256}", "dynamic")
	setVal, diags := types.SetValue(createClaimObjectType(), []attr.Value{claimObj})
	require.False(t, diags.HasError())

	req := validator.SetRequest{
		Path:        path.Root("custom_claims"),
		ConfigValue: setVal,
	}
	resp := &validator.SetResponse{
		Diagnostics: diag.Diagnostics{},
	}

	v.ValidateSet(ctx, req, resp)
	assert.False(t, resp.Diagnostics.HasError())
}

// TestCustomClaimsValidator_NonObjectElement verifies that non-object elements in the set are safely skipped.
func TestCustomClaimsValidator_NonObjectElement(t *testing.T) {
	t.Parallel()

	v := NewCustomClaimsValidator()
	ctx := context.Background()

	setVal, diags := types.SetValue(types.StringType, []attr.Value{types.StringValue("not an object")})
	require.False(t, diags.HasError())

	req := validator.SetRequest{
		Path:        path.Root("custom_claims"),
		ConfigValue: setVal,
	}
	resp := &validator.SetResponse{
		Diagnostics: diag.Diagnostics{},
	}

	v.ValidateSet(ctx, req, resp)
	assert.False(t, resp.Diagnostics.HasError())
}

// TestCustomClaimsValidator_InvalidObjectSchema verifies that object conversion errors are reported and skipped.
func TestCustomClaimsValidator_InvalidObjectSchema(t *testing.T) {
	t.Parallel()

	v := NewCustomClaimsValidator()
	ctx := context.Background()

	mismatchedType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"key":        types.Int64Type,
			"value":      types.StringType,
			"value_type": types.StringType,
		},
	}
	obj, diags := types.ObjectValue(
		mismatchedType.AttrTypes,
		map[string]attr.Value{
			"key":        types.Int64Value(123),
			"value":      types.StringValue("val"),
			"value_type": types.StringValue("literal"),
		},
	)
	require.False(t, diags.HasError())

	setVal, diags := types.SetValue(mismatchedType, []attr.Value{obj})
	require.False(t, diags.HasError())

	req := validator.SetRequest{
		Path:        path.Root("custom_claims"),
		ConfigValue: setVal,
	}
	resp := &validator.SetResponse{
		Diagnostics: diag.Diagnostics{},
	}

	v.ValidateSet(ctx, req, resp)
	assert.True(t, resp.Diagnostics.HasError())
}
