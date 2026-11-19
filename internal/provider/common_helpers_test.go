package provider

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
)

// skipNotCI can be used to skip tests which can ONLY run on GitHub.
//
//nolint:unused
func skipNotCI(t *testing.T) {
	if os.Getenv("CI") == "" {
		t.Skip("Skipping testing in non CI environment")
	}
}

func getTerraformVersion() string {
	cmd := exec.Command("terraform", "version")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		return "v1.6" // return the lowest version if something goes wrong
	}

	terraformVersion := strings.Split(strings.TrimSpace(string(output)), "\n")[0]
	fmt.Println(terraformVersion)

	return terraformVersion
}

func TestNewHTTPHeadersModel(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	testCases := []struct {
		name       string
		headers    []aembit.KeyValuePair
		expectNull bool
		expected   map[string]string
	}{
		{
			name:       "nil headers returns null map",
			headers:    nil,
			expectNull: true,
		},
		{
			name:       "empty headers returns null map",
			headers:    []aembit.KeyValuePair{},
			expectNull: true,
		},
		{
			name: "populated headers returns populated map",
			headers: []aembit.KeyValuePair{
				{Key: "X-Custom-A", Value: "val-a"},
				{Key: "Authorization", Value: "Bearer token"},
			},
			expectNull: false,
			expected: map[string]string{
				"X-Custom-A":    "val-a",
				"Authorization": "Bearer token",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := newHTTPHeadersModel(ctx, tc.headers)
			if tc.expectNull {
				require.True(t, got.IsNull())
				require.Equal(t, types.StringType, got.ElementType(ctx))
				return
			}

			require.False(t, got.IsNull())
			expectedMap, diags := types.MapValueFrom(ctx, types.StringType, tc.expected)
			require.False(t, diags.HasError())
			require.True(t, expectedMap.Equal(got), "expected %v, got %v", expectedMap, got)
		})
	}
}

func TestGetResourceSetId(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		resourceSetId types.String
		client        *aembit.CloudClient
		expected      string
	}{
		{
			name:          "null resource_set_id with nil client uses default",
			resourceSetId: types.StringNull(),
			client:        nil,
			expected:      DEFAULT_RESOURCESET_ID,
		},
		{
			name:          "unknown resource_set_id with empty client uses default",
			resourceSetId: types.StringUnknown(),
			client:        &aembit.CloudClient{},
			expected:      DEFAULT_RESOURCESET_ID,
		},
		{
			name:          "empty string resource_set_id uses default",
			resourceSetId: types.StringValue(""),
			client:        &aembit.CloudClient{ResourceSetId: DEFAULT_RESOURCESET_ID},
			expected:      DEFAULT_RESOURCESET_ID,
		},
		{
			name:          "explicit resource_set_id preserved when client has default",
			resourceSetId: types.StringValue("11111111-1111-1111-1111-111111111111"),
			client:        &aembit.CloudClient{ResourceSetId: DEFAULT_RESOURCESET_ID},
			expected:      "11111111-1111-1111-1111-111111111111",
		},
		{
			name:          "explicit resource_set_id preserved when client has empty",
			resourceSetId: types.StringValue("11111111-1111-1111-1111-111111111111"),
			client:        &aembit.CloudClient{ResourceSetId: ""},
			expected:      "11111111-1111-1111-1111-111111111111",
		},
		{
			name:          "non-default client resource_set_id overrides resource_set_id",
			resourceSetId: types.StringValue("11111111-1111-1111-1111-111111111111"),
			client:        &aembit.CloudClient{ResourceSetId: "22222222-2222-2222-2222-222222222222"},
			expected:      "22222222-2222-2222-2222-222222222222",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := getResourceSetId(tc.resourceSetId, tc.client)
			require.Equal(t, tc.expected, got)
		})
	}
}

func makeResourceSetTestSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"resource_set_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
		},
	}
}

func makeResourceSetTestObjectType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"resource_set_id": tftypes.String,
		},
	}
}

func TestModifyPlanForResourceSetId(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	testSchema := makeResourceSetTestSchema()
	objType := makeResourceSetTestObjectType()

	makeRaw := func(val tftypes.Value) tftypes.Value {
		return tftypes.NewValue(objType, map[string]tftypes.Value{
			"resource_set_id": val,
		})
	}

	nullObj := tftypes.NewValue(objType, nil)

	testCases := []struct {
		name            string
		configVal       tftypes.Value
		stateRaw        tftypes.Value
		planRaw         tftypes.Value
		schemaOverride  *schema.Schema
		client          *aembit.CloudClient
		expectError     bool
		expectPlanUnset bool
		expectedPlanVal attr.Value
	}{
		{
			name:            "destroy plan (null raw plan) exits early",
			configVal:       tftypes.NewValue(tftypes.String, nil),
			stateRaw:        makeRaw(tftypes.NewValue(tftypes.String, DEFAULT_RESOURCESET_ID)),
			planRaw:         nullObj,
			client:          &aembit.CloudClient{},
			expectPlanUnset: true,
		},
		{
			name:        "missing attribute in schema adds diagnostic error",
			configVal:   tftypes.NewValue(tftypes.String, nil),
			stateRaw:    nullObj,
			planRaw:     makeRaw(tftypes.NewValue(tftypes.String, tftypes.UnknownValue)),
			client:      &aembit.CloudClient{},
			expectError: true,
			schemaOverride: &schema.Schema{
				Attributes: map[string]schema.Attribute{},
			},
		},
		{
			name:            "unknown config value exits early without modifying plan",
			configVal:       tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			stateRaw:        nullObj,
			planRaw:         makeRaw(tftypes.NewValue(tftypes.String, tftypes.UnknownValue)),
			client:          &aembit.CloudClient{ResourceSetId: "22222222-2222-2222-2222-222222222222"},
			expectedPlanVal: types.StringUnknown(),
		},
		{
			name:            "explicit known config value leaves plan untouched",
			configVal:       tftypes.NewValue(tftypes.String, "11111111-1111-1111-1111-111111111111"),
			stateRaw:        nullObj,
			planRaw:         makeRaw(tftypes.NewValue(tftypes.String, "11111111-1111-1111-1111-111111111111")),
			client:          &aembit.CloudClient{ResourceSetId: "22222222-2222-2222-2222-222222222222"},
			expectedPlanVal: types.StringValue("11111111-1111-1111-1111-111111111111"),
		},
		{
			name:            "null config on create with nil client defaults to DEFAULT_RESOURCESET_ID",
			configVal:       tftypes.NewValue(tftypes.String, nil),
			stateRaw:        nullObj,
			planRaw:         makeRaw(tftypes.NewValue(tftypes.String, tftypes.UnknownValue)),
			client:          nil,
			expectedPlanVal: types.StringValue(DEFAULT_RESOURCESET_ID),
		},
		{
			name:            "null config on create with client resource_set_id uses client id",
			configVal:       tftypes.NewValue(tftypes.String, nil),
			stateRaw:        nullObj,
			planRaw:         makeRaw(tftypes.NewValue(tftypes.String, tftypes.UnknownValue)),
			client:          &aembit.CloudClient{ResourceSetId: "22222222-2222-2222-2222-222222222222"},
			expectedPlanVal: types.StringValue("22222222-2222-2222-2222-222222222222"),
		},
		{
			name:            "null config with matching state preserves state value",
			configVal:       tftypes.NewValue(tftypes.String, nil),
			stateRaw:        makeRaw(tftypes.NewValue(tftypes.String, "22222222-2222-2222-2222-222222222222")),
			planRaw:         makeRaw(tftypes.NewValue(tftypes.String, tftypes.UnknownValue)),
			client:          &aembit.CloudClient{ResourceSetId: "22222222-2222-2222-2222-222222222222"},
			expectedPlanVal: types.StringValue("22222222-2222-2222-2222-222222222222"),
		},
		{
			name:            "null config with different state updates plan to target id",
			configVal:       tftypes.NewValue(tftypes.String, nil),
			stateRaw:        makeRaw(tftypes.NewValue(tftypes.String, DEFAULT_RESOURCESET_ID)),
			planRaw:         makeRaw(tftypes.NewValue(tftypes.String, DEFAULT_RESOURCESET_ID)),
			client:          &aembit.CloudClient{ResourceSetId: "22222222-2222-2222-2222-222222222222"},
			expectedPlanVal: types.StringValue("22222222-2222-2222-2222-222222222222"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			activeSchema := testSchema
			if tc.schemaOverride != nil {
				activeSchema = *tc.schemaOverride
			}

			configRaw := makeRaw(tc.configVal)
			if tc.schemaOverride != nil {
				emptyObj := tftypes.Object{AttributeTypes: map[string]tftypes.Type{}}
				configRaw = tftypes.NewValue(emptyObj, map[string]tftypes.Value{})
				tc.planRaw = tftypes.NewValue(emptyObj, map[string]tftypes.Value{})
			}

			req := resource.ModifyPlanRequest{
				Config: tfsdk.Config{Raw: configRaw, Schema: activeSchema},
				State:  tfsdk.State{Raw: tc.stateRaw, Schema: activeSchema},
				Plan:   tfsdk.Plan{Raw: tc.planRaw, Schema: activeSchema},
			}
			resp := &resource.ModifyPlanResponse{
				Plan: tfsdk.Plan{Raw: tc.planRaw, Schema: activeSchema},
			}

			modifyPlanForResourceSetId(ctx, req, resp, tc.client)

			if tc.expectError {
				require.True(t, resp.Diagnostics.HasError())
				return
			}
			require.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)

			if tc.expectPlanUnset {
				require.True(t, resp.Plan.Raw.IsNull())
				return
			}

			var got types.String
			diags := resp.Plan.GetAttribute(ctx, path.Root("resource_set_id"), &got)
			require.False(t, diags.HasError())
			require.Equal(t, tc.expectedPlanVal, got)
		})
	}
}
