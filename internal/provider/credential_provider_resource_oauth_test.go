package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"testing"

	"terraform-provider-aembit/internal/provider/models"

	"aembit.io/aembit"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertToOAuthAuthorizationCodeDTO_IncludesFinalCallbackUrl(t *testing.T) {
	model := models.CredentialProviderResourceModel{
		Tags: types.MapNull(types.StringType),
		OAuthAuthorizationCode: &models.CredentialProviderOAuthAuthorizationCodeModel{
			OAuthDiscoveryUrl:    types.StringValue("https://aembit.io/.well-known/openid-configuration"),
			UserAuthorizationUrl: types.StringValue("https://aembit.io/user-authorize"),
			FinalCallbackUrl:     types.StringValue("https://aembit.io/final-callback"),
			State:                types.StringValue("state"),
			Lifetime:             31536000,
			OAuthCodeModel: models.OAuthCodeModel{
				OAuthAuthorizationUrl: types.StringValue("https://aembit.io/authorize"),
				OAuthTokenUrl:         types.StringValue("https://aembit.io/token"),
				OAuthIntrospectionUrl: types.StringValue("https://aembit.io/introspect"),
				ClientID:              types.StringValue("client-id"),
				ClientSecret:          types.StringValue("client-secret"),
				Scopes:                types.StringValue("scope"),
				IsPkceRequired:        types.BoolValue(true),
				CallBackUrl:           types.StringValue("https://aembit.io/callback"),
			},
		},
	}

	dto := convertCredentialProviderModelToV2DTO(
		context.Background(),
		model,
		nil,
		"tenant",
		"stack.example",
		nil,
	)

	if dto.FinalCallbackUrl != "https://aembit.io/final-callback" {
		t.Fatalf("expected final callback url to be mapped, got %q", dto.FinalCallbackUrl)
	}
}

func TestConvertCredentialProviderV2DTOToModel_OAuthAuthorizationCode_IncludesFinalCallbackUrl(
	t *testing.T,
) {
	dto := aembit.CredentialProviderV2DTO{
		Type:             "oauth-authorization-code",
		ClientID:         "client-id",
		Scope:            "scope",
		CustomParameters: []aembit.CustomClaimsDTO{},
		AuthorizationUrl: "https://aembit.io/authorize",
		TokenUrl:         "https://aembit.io/token",
		IntrospectionUrl: "https://aembit.io/introspect",
		CredentialOAuthAuthorizationCodeV2DTO: aembit.CredentialOAuthAuthorizationCodeV2DTO{
			IsPkceRequired:   true,
			CallBackUrl:      "https://aembit.io/callback",
			FinalCallbackUrl: "https://aembit.io/final-callback",
		},
	}

	model := convertCredentialProviderV2DTOToModel(
		context.Background(),
		dto,
		&models.CredentialProviderResourceModel{Tags: types.MapNull(types.StringType)},
		"tenant",
		"stack.example",
	)

	if model.OAuthAuthorizationCode == nil {
		t.Fatal("expected oauth authorization code model to be set")
	}

	if model.OAuthAuthorizationCode.FinalCallbackUrl.ValueString() != "https://aembit.io/final-callback" {
		t.Fatalf(
			"expected final callback url to be mapped, got %q",
			model.OAuthAuthorizationCode.FinalCallbackUrl.ValueString(),
		)
	}
}

func TestConvertToMcpUserBasedAccessTokenDTO_RevocationUrl(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name               string
		oauthRevocationUrl types.String
		expectedDTOValue   string
	}{
		{
			name:               "valid https revocation url",
			oauthRevocationUrl: types.StringValue("https://aembit.io/revoke"),
			expectedDTOValue:   "https://aembit.io/revoke",
		},
		{
			name:               "empty string revocation url",
			oauthRevocationUrl: types.StringValue(""),
			expectedDTOValue:   "",
		},
		{
			name:               "null omitted revocation url",
			oauthRevocationUrl: types.StringNull(),
			expectedDTOValue:   "",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			model := models.CredentialProviderResourceModel{
				Tags: types.MapNull(types.StringType),
				McpUserBasedAccessToken: &models.CredentialProviderMcpUserBasedAccessTokenModel{
					McpServerUrl:       types.StringValue("https://aembit.io/mcp"),
					Resource:           types.StringValue("resource-id"),
					OAuthRevocationUrl: tc.oauthRevocationUrl,
					OAuthCodeModel: models.OAuthCodeModel{
						OAuthAuthorizationUrl: types.StringValue("https://aembit.io/authorize"),
						OAuthTokenUrl:         types.StringValue("https://aembit.io/token"),
						OAuthIntrospectionUrl: types.StringValue("https://aembit.io/introspect"),
						ClientID:              types.StringValue("client-id"),
						ClientSecret:          types.StringValue("client-secret"),
						Scopes:                types.StringValue("scope"),
						IsPkceRequired:        types.BoolValue(true),
						CallBackUrl:           types.StringValue("https://aembit.io/callback"),
					},
				},
			}

			dto := convertCredentialProviderModelToV2DTO(
				context.Background(),
				model,
				nil,
				"tenant",
				"stack.example",
				nil,
			)

			assert.Equal(t, "mcp-user-based-access-token", dto.Type)
			assert.Equal(t, tc.expectedDTOValue, dto.RevocationUrl)
		})
	}
}

func TestConvertCredentialProviderV2DTOToModel_McpUserBasedAccessToken_RevocationUrl(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name             string
		dtoRevocationUrl string
		expectedModelVal string
	}{
		{
			name:             "populated revocation url",
			dtoRevocationUrl: "https://aembit.io/revoke",
			expectedModelVal: "https://aembit.io/revoke",
		},
		{
			name:             "empty revocation url defaults to empty string",
			dtoRevocationUrl: "",
			expectedModelVal: "",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dto := aembit.CredentialProviderV2DTO{
				Type:             "mcp-user-based-access-token",
				ClientID:         "client-id",
				Scope:            "scope",
				CustomParameters: []aembit.CustomClaimsDTO{},
				McpServerUrl:     "https://aembit.io/mcp",
				AuthorizationUrl: "https://aembit.io/authorize",
				TokenUrl:         "https://aembit.io/token",
				IntrospectionUrl: "https://aembit.io/introspect",
				RevocationUrl:    tc.dtoRevocationUrl,
				CredentialOAuthAuthorizationCodeV2DTO: aembit.CredentialOAuthAuthorizationCodeV2DTO{
					IsPkceRequired: true,
					CallBackUrl:    "https://aembit.io/callback",
					Resource:       "resource-id",
				},
			}

			model := convertCredentialProviderV2DTOToModel(
				context.Background(),
				dto,
				&models.CredentialProviderResourceModel{
					Tags: types.MapNull(types.StringType),
					McpUserBasedAccessToken: &models.CredentialProviderMcpUserBasedAccessTokenModel{
						OAuthCodeModel: models.OAuthCodeModel{
							ClientSecret: types.StringValue("retained-secret"),
						},
					},
				},
				"tenant",
				"stack.example",
			)

			require.NotNil(t, model.McpUserBasedAccessToken)
			assert.Equal(t, tc.expectedModelVal, model.McpUserBasedAccessToken.OAuthRevocationUrl.ValueString())
			assert.Equal(t, "retained-secret", model.McpUserBasedAccessToken.ClientSecret.ValueString())
		})
	}
}

func TestConvertOAuthAuthorizationCode_UnaffectedByRevocationUrl(t *testing.T) {
	t.Parallel()

	model := models.CredentialProviderResourceModel{
		Tags: types.MapNull(types.StringType),
		OAuthAuthorizationCode: &models.CredentialProviderOAuthAuthorizationCodeModel{
			OAuthDiscoveryUrl:    types.StringValue("https://aembit.io/.well-known/openid-configuration"),
			UserAuthorizationUrl: types.StringValue("https://aembit.io/user-authorize"),
			FinalCallbackUrl:     types.StringValue("https://aembit.io/final-callback"),
			State:                types.StringValue("state"),
			Lifetime:             31536000,
			OAuthCodeModel: models.OAuthCodeModel{
				OAuthAuthorizationUrl: types.StringValue("https://aembit.io/authorize"),
				OAuthTokenUrl:         types.StringValue("https://aembit.io/token"),
				OAuthIntrospectionUrl: types.StringValue("https://aembit.io/introspect"),
				ClientID:              types.StringValue("client-id"),
				ClientSecret:          types.StringValue("client-secret"),
				Scopes:                types.StringValue("scope"),
				IsPkceRequired:        types.BoolValue(true),
				CallBackUrl:           types.StringValue("https://aembit.io/callback"),
			},
		},
	}

	dto := convertCredentialProviderModelToV2DTO(
		context.Background(),
		model,
		nil,
		"tenant",
		"stack.example",
		nil,
	)
	assert.Empty(t, dto.RevocationUrl)

	dto.RevocationUrl = "https://aembit.io/revoke"
	converted := convertCredentialProviderV2DTOToModel(
		context.Background(),
		dto,
		&models.CredentialProviderResourceModel{Tags: types.MapNull(types.StringType)},
		"tenant",
		"stack.example",
	)
	require.NotNil(t, converted.OAuthAuthorizationCode)
	assert.Nil(t, converted.McpUserBasedAccessToken)
}

func TestCredentialProviderResourceSchema_McpUserBasedAccessToken_OAuthRevocationUrl(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	res := NewCredentialProviderResource()
	schemaResp := &fwresource.SchemaResponse{}
	res.Schema(ctx, fwresource.SchemaRequest{}, schemaResp)
	require.False(t, schemaResp.Diagnostics.HasError())

	oauthAuthCodeAttr, ok := schemaResp.Schema.Attributes["oauth_authorization_code"].(schema.SingleNestedAttribute)
	require.True(t, ok)
	_, hasRevocationInOAuthAuthCode := oauthAuthCodeAttr.Attributes["oauth_revocation_url"]
	assert.False(t, hasRevocationInOAuthAuthCode, "oauth_authorization_code must not expose oauth_revocation_url")

	mcpAttr, ok := schemaResp.Schema.Attributes["mcp_user_based_access_token"].(schema.SingleNestedAttribute)
	require.True(t, ok)
	revocationAttr, ok := mcpAttr.Attributes["oauth_revocation_url"].(schema.StringAttribute)
	require.True(t, ok, "mcp_user_based_access_token must expose oauth_revocation_url")

	assert.True(t, revocationAttr.Optional)
	assert.True(t, revocationAttr.Computed)
	assert.False(t, revocationAttr.Required)
	require.NotNil(t, revocationAttr.Default)

	defaultResp := &defaults.StringResponse{}
	revocationAttr.Default.DefaultString(ctx, defaults.StringRequest{}, defaultResp)
	require.False(t, defaultResp.Diagnostics.HasError())
	assert.Equal(t, types.StringValue(""), defaultResp.PlanValue)

	require.NotEmpty(t, revocationAttr.Validators)

	validCases := []struct {
		name  string
		value types.String
	}{
		{name: "valid https url", value: types.StringValue("https://aembit.io/revoke")},
		{name: "valid https url with path and query", value: types.StringValue("https://auth.example.com/oauth2/revoke?hint=access_token")},
		{name: "omitted null value", value: types.StringNull()},
		{name: "unknown value", value: types.StringUnknown()},
		{name: "explicit empty string", value: types.StringValue("")},
	}
	for _, tc := range validCases {
		tc := tc
		t.Run("Valid_"+tc.name, func(t *testing.T) {
			t.Parallel()
			for _, v := range revocationAttr.Validators {
				resp := &validator.StringResponse{Diagnostics: diag.Diagnostics{}}
				v.ValidateString(ctx, validator.StringRequest{
					Path:        path.Root("mcp_user_based_access_token").AtName("oauth_revocation_url"),
					ConfigValue: tc.value,
				}, resp)
				require.False(t, resp.Diagnostics.HasError(), "expected %v to be valid, got: %v", tc.value, resp.Diagnostics)
			}
		})
	}

	invalidCases := []struct {
		name  string
		value string
	}{
		{name: "non-https http url", value: "http://aembit.io/revoke"},
		{name: "ftp scheme", value: "ftp://aembit.io/revoke"},
		{name: "localhost without tld", value: "https://localhost/revoke"},
		{name: "malformed url", value: "not-a-valid-url"},
		{name: "missing scheme", value: "aembit.io/revoke"},
		{name: "whitespace only", value: "   "},
	}
	for _, tc := range invalidCases {
		tc := tc
		t.Run("Invalid_"+tc.name, func(t *testing.T) {
			t.Parallel()
			hasError := false
			for _, v := range revocationAttr.Validators {
				resp := &validator.StringResponse{Diagnostics: diag.Diagnostics{}}
				v.ValidateString(ctx, validator.StringRequest{
					Path:        path.Root("mcp_user_based_access_token").AtName("oauth_revocation_url"),
					ConfigValue: types.StringValue(tc.value),
				}, resp)
				if resp.Diagnostics.HasError() {
					hasError = true
					require.Len(t, resp.Diagnostics.Errors(), 1)
					assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "must be a valid https url")
				}
			}
			assert.True(t, hasError, "expected %q to fail validation", tc.value)
		})
	}

	t.Run("ValidateResourceConfig", func(t *testing.T) {
		t.Parallel()

		server, err := providerserver.NewProtocol6WithError(New("test", "unittest")())()
		require.NoError(t, err)

		schemaType := schemaResp.Schema.Type().TerraformType(ctx)

		validateCases := []struct {
			name          string
			revocationUrl types.String
			expectError   bool
		}{
			{name: "explicit empty string", revocationUrl: types.StringValue(""), expectError: false},
			{name: "omitted null value", revocationUrl: types.StringNull(), expectError: false},
			{name: "valid https url", revocationUrl: types.StringValue("https://aembit.io/revoke"), expectError: false},
			{name: "whitespace only", revocationUrl: types.StringValue("   "), expectError: true},
			{name: "non-https url", revocationUrl: types.StringValue("http://aembit.io/revoke"), expectError: true},
		}

		for _, tc := range validateCases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				cfgPlan := tfsdk.Plan{Schema: schemaResp.Schema}
				diags := cfgPlan.Set(ctx, models.CredentialProviderResourceModel{
					Name:     types.StringValue("MCP Validate Config Test"),
					IsActive: types.BoolValue(true),
					Tags:     types.MapNull(types.StringType),
					TagsAll:  types.MapNull(types.StringType),
					McpUserBasedAccessToken: &models.CredentialProviderMcpUserBasedAccessTokenModel{
						McpServerUrl:       types.StringValue("https://aembit.io/mcp"),
						OAuthRevocationUrl: tc.revocationUrl,
						OAuthCodeModel: models.OAuthCodeModel{
							OAuthAuthorizationUrl: types.StringValue("https://aembit.io/authorize"),
							OAuthTokenUrl:         types.StringValue("https://aembit.io/token"),
							ClientID:              types.StringValue("client-id"),
						},
					},
				})
				require.False(t, diags.HasError(), "unexpected plan.Set diagnostics: %v", diags)

				dynVal, dynErr := tfprotov6.NewDynamicValue(schemaType, cfgPlan.Raw)
				require.NoError(t, dynErr)

				valResp, valErr := server.ValidateResourceConfig(ctx, &tfprotov6.ValidateResourceConfigRequest{
					TypeName: "aembit_credential_provider",
					Config:   &dynVal,
				})
				require.NoError(t, valErr)

				var hasError bool
				for _, d := range valResp.Diagnostics {
					if d.Severity == tfprotov6.DiagnosticSeverityError {
						hasError = true
						assert.Contains(t, d.Detail, "must be a valid https url")
					}
				}
				assert.Equal(t, tc.expectError, hasError, "diagnostics: %+v", valResp.Diagnostics)
			})
		}
	})
}

func newMockCredentialProviderCloudClient(
	t *testing.T,
	externalID string,
	listedProviders []aembit.CredentialProviderV2DTO,
) (*aembit.CloudClient, func() []aembit.CredentialProviderV2DTO) {
	t.Helper()

	var (
		mu       sync.Mutex
		stored   *aembit.CredentialProviderV2DTO
		requests []aembit.CredentialProviderV2DTO
	)

	token := "unit-test-token"
	mockClient, err := aembit.NewClient(aembit.URLBuilder{}, &token, "", "unittest")
	require.NoError(t, err)
	mockClient.Tenant = "unit-tenant"
	mockClient.StackDomain = "aembit.io"
	mockClient.HTTPClient = &http.Client{
		Transport: &mockRoundTripper{
			fn: func(req *http.Request) (*http.Response, error) {
				mu.Lock()
				defer mu.Unlock()

				switch req.Method {
				case http.MethodPost:
					bodyBytes, readErr := io.ReadAll(req.Body)
					if readErr != nil {
						return nil, readErr
					}
					var incoming aembit.CredentialProviderV2DTO
					if unmarshalErr := json.Unmarshal(bodyBytes, &incoming); unmarshalErr != nil {
						return nil, unmarshalErr
					}
					requests = append(requests, incoming)
					incoming.ExternalID = externalID
					incoming.CallBackUrl = "https://unit-tenant.mcp.aembit.io/userauth/" + externalID + "/callback"
					storedCopy := incoming
					stored = &storedCopy

					respBytes, marshalErr := json.Marshal(stored)
					if marshalErr != nil {
						return nil, marshalErr
					}
					return &http.Response{
						StatusCode: http.StatusCreated,
						Body:       io.NopCloser(bytes.NewReader(respBytes)),
						Header:     make(http.Header),
					}, nil

				case http.MethodPut:
					bodyBytes, readErr := io.ReadAll(req.Body)
					if readErr != nil {
						return nil, readErr
					}
					var incoming aembit.CredentialProviderV2DTO
					if unmarshalErr := json.Unmarshal(bodyBytes, &incoming); unmarshalErr != nil {
						return nil, unmarshalErr
					}
					requests = append(requests, incoming)
					incoming.ExternalID = externalID
					incoming.CallBackUrl = "https://unit-tenant.mcp.aembit.io/userauth/" + externalID + "/callback"
					storedCopy := incoming
					stored = &storedCopy

					respBytes, marshalErr := json.Marshal(stored)
					if marshalErr != nil {
						return nil, marshalErr
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(respBytes)),
						Header:     make(http.Header),
					}, nil

				case http.MethodGet:
					if listedProviders != nil {
						resp := aembit.CredentialProvidersV2Response{
							CredentialProviders: listedProviders,
						}
						respBytes, marshalErr := json.Marshal(resp)
						if marshalErr != nil {
							return nil, marshalErr
						}
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewReader(respBytes)),
							Header:     make(http.Header),
						}, nil
					}
					if stored == nil {
						return &http.Response{
							StatusCode: http.StatusNoContent,
							Body:       io.NopCloser(bytes.NewReader(nil)),
							Header:     make(http.Header),
						}, nil
					}
					respBytes, marshalErr := json.Marshal(stored)
					if marshalErr != nil {
						return nil, marshalErr
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(respBytes)),
						Header:     make(http.Header),
					}, nil

				default:
					return &http.Response{
						StatusCode: http.StatusMethodNotAllowed,
						Body:       io.NopCloser(bytes.NewReader(nil)),
						Header:     make(http.Header),
					}, nil
				}
			},
		},
	}

	getRequests := func() []aembit.CredentialProviderV2DTO {
		mu.Lock()
		defer mu.Unlock()
		out := make([]aembit.CredentialProviderV2DTO, len(requests))
		copy(out, requests)
		return out
	}

	return mockClient, getRequests
}

func TestCredentialProviderResource_McpUserBasedAccessToken_Lifecycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	const externalID = "11111111-2222-3333-4444-555555555555"
	mockClient, getRequests := newMockCredentialProviderCloudClient(t, externalID, nil)

	res, ok := NewCredentialProviderResource().(*credentialProviderResource)
	require.True(t, ok)
	res.client = mockClient

	var schemaResp fwresource.SchemaResponse
	res.Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)
	require.False(t, schemaResp.Diagnostics.HasError())
	resSchema := schemaResp.Schema

	buildModel := func(id types.String, revocationUrl types.String) models.CredentialProviderResourceModel {
		return models.CredentialProviderResourceModel{
			ID:            id,
			ResourceSetId: types.StringValue(""),
			Name:          types.StringValue("MCP Revocation Lifecycle"),
			Description:   types.StringValue("Testing revocation url lifecycle"),
			IsActive:      types.BoolValue(true),
			Tags:          types.MapNull(types.StringType),
			TagsAll:       types.MapNull(types.StringType),
			McpUserBasedAccessToken: &models.CredentialProviderMcpUserBasedAccessTokenModel{
				McpServerUrl:       types.StringValue("https://aembit.io/mcp"),
				Resource:           types.StringValue("resource-id"),
				OAuthRevocationUrl: revocationUrl,
				OAuthCodeModel: models.OAuthCodeModel{
					OAuthAuthorizationUrl: types.StringValue("https://aembit.io/authorize"),
					OAuthTokenUrl:         types.StringValue("https://aembit.io/token"),
					OAuthIntrospectionUrl: types.StringValue(""),
					ClientID:              types.StringValue("client-id"),
					ClientSecret:          types.StringValue("client-secret"),
					Scopes:                types.StringValue("openid"),
					IsPkceRequired:        types.BoolValue(true),
					CallBackUrl:           types.StringUnknown(),
					CustomParameters:      []*models.CredentialProviderOAuthClientCustomParametersModel{},
				},
			},
		}
	}

	// 1. Create with oauth_revocation_url set.
	createPlan := tfsdk.Plan{Schema: resSchema}
	require.False(t, createPlan.Set(ctx, buildModel(types.StringUnknown(), types.StringValue("https://aembit.io/revoke"))).HasError())

	createResp := fwresource.CreateResponse{
		State: tfsdk.State{Schema: resSchema, Raw: tftypes.NewValue(resSchema.Type().TerraformType(ctx), nil)},
	}
	res.Create(ctx, fwresource.CreateRequest{Plan: createPlan}, &createResp)
	require.False(t, createResp.Diagnostics.HasError(), "Create diagnostics: %v", createResp.Diagnostics)

	var stateAfterCreate models.CredentialProviderResourceModel
	require.False(t, createResp.State.Get(ctx, &stateAfterCreate).HasError())
	require.NotNil(t, stateAfterCreate.McpUserBasedAccessToken)
	assert.Equal(t, types.StringValue("https://aembit.io/revoke"), stateAfterCreate.McpUserBasedAccessToken.OAuthRevocationUrl)

	// 2. Update oauth_revocation_url to a new HTTPS URL.
	updatePlan := tfsdk.Plan{Schema: resSchema}
	require.False(t, updatePlan.Set(ctx, buildModel(types.StringValue(externalID), types.StringValue("https://aembit.io/revoke-updated"))).HasError())

	updateResp := fwresource.UpdateResponse{State: createResp.State}
	res.Update(ctx, fwresource.UpdateRequest{
		State: createResp.State,
		Plan:  updatePlan,
	}, &updateResp)
	require.False(t, updateResp.Diagnostics.HasError(), "Update diagnostics: %v", updateResp.Diagnostics)

	var stateAfterUpdate models.CredentialProviderResourceModel
	require.False(t, updateResp.State.Get(ctx, &stateAfterUpdate).HasError())
	require.NotNil(t, stateAfterUpdate.McpUserBasedAccessToken)
	assert.Equal(t, types.StringValue("https://aembit.io/revoke-updated"), stateAfterUpdate.McpUserBasedAccessToken.OAuthRevocationUrl)

	// 3. Clear oauth_revocation_url back to "".
	clearPlan := tfsdk.Plan{Schema: resSchema}
	require.False(t, clearPlan.Set(ctx, buildModel(types.StringValue(externalID), types.StringValue(""))).HasError())

	clearResp := fwresource.UpdateResponse{State: updateResp.State}
	res.Update(ctx, fwresource.UpdateRequest{
		State: updateResp.State,
		Plan:  clearPlan,
	}, &clearResp)
	require.False(t, clearResp.Diagnostics.HasError(), "Clear diagnostics: %v", clearResp.Diagnostics)

	var stateAfterClear models.CredentialProviderResourceModel
	require.False(t, clearResp.State.Get(ctx, &stateAfterClear).HasError())
	require.NotNil(t, stateAfterClear.McpUserBasedAccessToken)
	assert.Equal(t, types.StringValue(""), stateAfterClear.McpUserBasedAccessToken.OAuthRevocationUrl)

	// 4. Read refreshes state and preserves cleared "".
	readResp := fwresource.ReadResponse{State: clearResp.State}
	res.Read(ctx, fwresource.ReadRequest{State: clearResp.State}, &readResp)
	require.False(t, readResp.Diagnostics.HasError(), "Read diagnostics: %v", readResp.Diagnostics)

	var stateAfterRead models.CredentialProviderResourceModel
	require.False(t, readResp.State.Get(ctx, &stateAfterRead).HasError())
	require.NotNil(t, stateAfterRead.McpUserBasedAccessToken)
	assert.Equal(t, types.StringValue(""), stateAfterRead.McpUserBasedAccessToken.OAuthRevocationUrl)

	reqs := getRequests()
	require.Len(t, reqs, 3)
	assert.Equal(t, "https://aembit.io/revoke", reqs[0].RevocationUrl)
	assert.Equal(t, "https://aembit.io/revoke-updated", reqs[1].RevocationUrl)
	assert.Equal(t, "", reqs[2].RevocationUrl)
}

func TestCredentialProviderResource_McpUserBasedAccessToken_PlanOmissionAfterNonEmpty(t *testing.T) {
	t.Setenv("AEMBIT_CLIENT_ID", "")
	t.Setenv("AEMBIT_RESOURCE_SET_ID", "")

	ctx := context.Background()
	const externalID = "11111111-2222-3333-4444-555555555555"

	p := New("test", "unittest")()
	providerSchemaResp := &fwprovider.SchemaResponse{}
	p.Schema(ctx, fwprovider.SchemaRequest{}, providerSchemaResp)
	require.False(t, providerSchemaResp.Diagnostics.HasError())

	providerCfgPlan := tfsdk.Plan{Schema: providerSchemaResp.Schema}
	require.False(t, providerCfgPlan.Set(ctx, aembitProviderModel{
		Tenant:      types.StringValue("unit-tenant"),
		Token:       types.StringValue("unit-test-token"),
		ClientID:    types.StringNull(),
		ResourceSet: types.StringNull(),
	}).HasError())

	providerDynVal, err := tfprotov6.NewDynamicValue(
		providerSchemaResp.Schema.Type().TerraformType(ctx),
		providerCfgPlan.Raw,
	)
	require.NoError(t, err)

	server, err := providerserver.NewProtocol6WithError(New("test", "unittest")())()
	require.NoError(t, err)

	cfgResp, err := server.ConfigureProvider(ctx, &tfprotov6.ConfigureProviderRequest{
		Config: &providerDynVal,
	})
	require.NoError(t, err)
	for _, d := range cfgResp.Diagnostics {
		require.NotEqual(t, tfprotov6.DiagnosticSeverityError, d.Severity, "ConfigureProvider error: %s", d.Detail)
	}

	mockClient, getRequests := newMockCredentialProviderCloudClient(t, externalID, nil)
	res, ok := NewCredentialProviderResource().(*credentialProviderResource)
	require.True(t, ok)
	res.client = mockClient

	var schemaResp fwresource.SchemaResponse
	res.Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)
	require.False(t, schemaResp.Diagnostics.HasError())
	resSchema := schemaResp.Schema
	schemaType := resSchema.Type().TerraformType(ctx)

	buildModel := func(id types.String, revocationUrl types.String) models.CredentialProviderResourceModel {
		return models.CredentialProviderResourceModel{
			ID:            id,
			ResourceSetId: types.StringValue(DEFAULT_RESOURCESET_ID),
			Name:          types.StringValue("MCP Revocation Plan Omission"),
			Description:   types.StringValue("Testing plan omission after non-empty"),
			IsActive:      types.BoolValue(true),
			Tags:          types.MapNull(types.StringType),
			TagsAll:       types.MapNull(types.StringType),
			McpUserBasedAccessToken: &models.CredentialProviderMcpUserBasedAccessTokenModel{
				McpServerUrl:       types.StringValue("https://aembit.io/mcp"),
				Resource:           types.StringValue("resource-id"),
				OAuthRevocationUrl: revocationUrl,
				OAuthCodeModel: models.OAuthCodeModel{
					OAuthAuthorizationUrl: types.StringValue("https://aembit.io/authorize"),
					OAuthTokenUrl:         types.StringValue("https://aembit.io/token"),
					OAuthIntrospectionUrl: types.StringValue(""),
					ClientID:              types.StringValue("client-id"),
					ClientSecret:          types.StringValue("client-secret"),
					Scopes:                types.StringValue("openid"),
					IsPkceRequired:        types.BoolValue(true),
					CallBackUrl:           types.StringValue("https://unit-tenant.mcp.aembit.io/userauth/" + externalID + "/callback"),
					CustomParameters:      []*models.CredentialProviderOAuthClientCustomParametersModel{},
				},
			},
		}
	}

	// 1. Seed initial resource with non-empty oauth_revocation_url.
	createPlan := tfsdk.Plan{Schema: resSchema}
	require.False(t, createPlan.Set(ctx, buildModel(types.StringUnknown(), types.StringValue("https://aembit.io/revoke-updated"))).HasError())
	createResp := fwresource.CreateResponse{
		State: tfsdk.State{Schema: resSchema, Raw: tftypes.NewValue(schemaType, nil)},
	}
	res.Create(ctx, fwresource.CreateRequest{Plan: createPlan}, &createResp)
	require.False(t, createResp.Diagnostics.HasError(), "Create diagnostics: %v", createResp.Diagnostics)

	// 2. Plan a config where oauth_revocation_url is omitted (null) immediately after non-empty state.
	// ProposedNewState carries forward the prior non-empty value as Terraform Core does for Optional+Computed.
	omittedCfgPlan := tfsdk.Plan{Schema: resSchema}
	omittedModel := buildModel(types.StringNull(), types.StringNull())
	omittedModel.ResourceSetId = types.StringNull()
	omittedModel.McpUserBasedAccessToken.CallBackUrl = types.StringNull()
	require.False(t, omittedCfgPlan.Set(ctx, omittedModel).HasError())

	proposedState := tfsdk.Plan{Schema: resSchema}
	require.False(t, proposedState.Set(ctx, buildModel(types.StringValue(externalID), types.StringValue("https://aembit.io/revoke-updated"))).HasError())

	priorDynVal, err := tfprotov6.NewDynamicValue(schemaType, createResp.State.Raw)
	require.NoError(t, err)
	configDynVal, err := tfprotov6.NewDynamicValue(schemaType, omittedCfgPlan.Raw)
	require.NoError(t, err)
	proposedDynVal, err := tfprotov6.NewDynamicValue(schemaType, proposedState.Raw)
	require.NoError(t, err)

	planResp, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
		TypeName:         "aembit_credential_provider",
		PriorState:       &priorDynVal,
		Config:           &configDynVal,
		ProposedNewState: &proposedDynVal,
	})
	require.NoError(t, err)
	for _, d := range planResp.Diagnostics {
		require.NotEqual(t, tfprotov6.DiagnosticSeverityError, d.Severity, "PlanResourceChange error: %s", d.Detail)
	}

	plannedTfVal, err := planResp.PlannedState.Unmarshal(schemaType)
	require.NoError(t, err)
	plannedPlan := tfsdk.Plan{Schema: resSchema, Raw: plannedTfVal}

	var plannedModel models.CredentialProviderResourceModel
	require.False(t, plannedPlan.Get(ctx, &plannedModel).HasError())
	require.NotNil(t, plannedModel.McpUserBasedAccessToken)
	assert.Equal(t, types.StringValue(""), plannedModel.McpUserBasedAccessToken.OAuthRevocationUrl)

	// 3. Apply the planned state via Update and verify PUT sends "" and Read returns "".
	updateResp := fwresource.UpdateResponse{State: createResp.State}
	res.Update(ctx, fwresource.UpdateRequest{
		State: createResp.State,
		Plan:  plannedPlan,
	}, &updateResp)
	require.False(t, updateResp.Diagnostics.HasError(), "Update diagnostics: %v", updateResp.Diagnostics)

	readResp := fwresource.ReadResponse{State: updateResp.State}
	res.Read(ctx, fwresource.ReadRequest{State: updateResp.State}, &readResp)
	require.False(t, readResp.Diagnostics.HasError(), "Read diagnostics: %v", readResp.Diagnostics)

	var stateAfterRead models.CredentialProviderResourceModel
	require.False(t, readResp.State.Get(ctx, &stateAfterRead).HasError())
	require.NotNil(t, stateAfterRead.McpUserBasedAccessToken)
	assert.Equal(t, types.StringValue(""), stateAfterRead.McpUserBasedAccessToken.OAuthRevocationUrl)

	reqs := getRequests()
	require.Len(t, reqs, 2)
	assert.Equal(t, "https://aembit.io/revoke-updated", reqs[0].RevocationUrl)
	assert.Equal(t, "", reqs[1].RevocationUrl)
}

func TestCredentialProvidersDataSourceSchema_OAuthRevocationUrl(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ds := NewCredentialProvidersDataSource()
	schemaResp := &datasource.SchemaResponse{}
	ds.Schema(ctx, datasource.SchemaRequest{}, schemaResp)
	require.False(t, schemaResp.Diagnostics.HasError())

	listAttr, ok := schemaResp.Schema.Attributes["credential_providers"].(dsschema.ListNestedAttribute)
	require.True(t, ok)

	oauthAuthCodeAttr, ok := listAttr.NestedObject.Attributes["oauth_authorization_code"].(dsschema.SingleNestedAttribute)
	require.True(t, ok)
	_, hasRevocationInOAuthAuthCode := oauthAuthCodeAttr.Attributes["oauth_revocation_url"]
	assert.False(t, hasRevocationInOAuthAuthCode, "data source oauth_authorization_code must not expose oauth_revocation_url")

	mcpAttr, ok := listAttr.NestedObject.Attributes["mcp_user_based_access_token"].(dsschema.SingleNestedAttribute)
	require.True(t, ok)
	revocationAttr, ok := mcpAttr.Attributes["oauth_revocation_url"].(dsschema.StringAttribute)
	require.True(t, ok, "data source mcp_user_based_access_token must expose oauth_revocation_url")
	assert.True(t, revocationAttr.Computed)
	assert.False(t, revocationAttr.Optional)
	assert.False(t, revocationAttr.Required)
}

func TestCredentialProvidersDataSource_Read_McpUserBasedAccessTokenRevocationUrl(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	listed := []aembit.CredentialProviderV2DTO{
		{
			AccessEntityDTO: aembit.AccessEntityDTO{
				EntityDTO: aembit.EntityDTO{
					ExternalID:  "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
					Name:        "MCP With Revocation",
					Description: "With revocation URL",
					IsActive:    true,
				},
			},
			Type:             "mcp-user-based-access-token",
			ClientID:         "client-1",
			Scope:            "openid",
			McpServerUrl:     "https://aembit.io/mcp",
			AuthorizationUrl: "https://aembit.io/authorize",
			TokenUrl:         "https://aembit.io/token",
			IntrospectionUrl: "https://aembit.io/introspect",
			RevocationUrl:    "https://aembit.io/revoke",
			CustomParameters: []aembit.CustomClaimsDTO{},
			CredentialOAuthAuthorizationCodeV2DTO: aembit.CredentialOAuthAuthorizationCodeV2DTO{
				IsPkceRequired: true,
				CallBackUrl:    "https://aembit.io/callback",
				Resource:       "res-1",
			},
		},
		{
			AccessEntityDTO: aembit.AccessEntityDTO{
				EntityDTO: aembit.EntityDTO{
					ExternalID:  "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
					Name:        "MCP Default Empty Revocation",
					Description: "Empty revocation URL",
					IsActive:    true,
				},
			},
			Type:             "mcp-user-based-access-token",
			ClientID:         "client-2",
			Scope:            "openid",
			McpServerUrl:     "https://aembit.io/mcp",
			AuthorizationUrl: "https://aembit.io/authorize",
			TokenUrl:         "https://aembit.io/token",
			RevocationUrl:    "",
			CustomParameters: []aembit.CustomClaimsDTO{},
			CredentialOAuthAuthorizationCodeV2DTO: aembit.CredentialOAuthAuthorizationCodeV2DTO{
				IsPkceRequired: false,
				CallBackUrl:    "https://aembit.io/callback2",
				Resource:       "res-2",
			},
		},
		{
			AccessEntityDTO: aembit.AccessEntityDTO{
				EntityDTO: aembit.EntityDTO{
					ExternalID:  "cccccccc-cccc-cccc-cccc-cccccccccccc",
					Name:        "OAuth Auth Code",
					Description: "Auth code CP",
					IsActive:    true,
				},
			},
			Type:             "oauth-authorization-code",
			ClientID:         "client-3",
			Scope:            "openid",
			AuthorizationUrl: "https://aembit.io/authorize",
			TokenUrl:         "https://aembit.io/token",
			RevocationUrl:    "https://ignored.example.com/revoke",
			CustomParameters: []aembit.CustomClaimsDTO{},
			CredentialOAuthAuthorizationCodeV2DTO: aembit.CredentialOAuthAuthorizationCodeV2DTO{
				OAuthUrl:             "https://aembit.io/.well-known/openid-configuration",
				UserAuthorizationUrl: "https://aembit.io/user-authorize",
				FinalCallbackUrl:     "https://aembit.io/final-callback",
				IsPkceRequired:       true,
				CallBackUrl:          "https://aembit.io/callback3",
			},
		},
	}

	mockClient, _ := newMockCredentialProviderCloudClient(t, "", listed)
	ds, ok := NewCredentialProvidersDataSource().(*credentialProvidersDataSource)
	require.True(t, ok)
	ds.client = mockClient

	var schemaResp datasource.SchemaResponse
	ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)
	require.False(t, schemaResp.Diagnostics.HasError())

	cfg := tfsdk.Config{
		Schema: schemaResp.Schema,
		Raw: tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
			"resource_set_id":      tftypes.NewValue(tftypes.String, nil),
			"credential_providers": tftypes.NewValue(schemaResp.Schema.Attributes["credential_providers"].GetType().TerraformType(ctx), nil),
		}),
	}
	readResp := datasource.ReadResponse{
		State: tfsdk.State{
			Schema: schemaResp.Schema,
			Raw: tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
				"resource_set_id":      tftypes.NewValue(tftypes.String, nil),
				"credential_providers": tftypes.NewValue(schemaResp.Schema.Attributes["credential_providers"].GetType().TerraformType(ctx), nil),
			}),
		},
	}

	ds.Read(ctx, datasource.ReadRequest{Config: cfg}, &readResp)
	require.False(t, readResp.Diagnostics.HasError(), "data source Read diagnostics: %v", readResp.Diagnostics)

	var state models.CredentialProvidersDataSourceModel
	require.False(t, readResp.State.Get(ctx, &state).HasError())
	require.Len(t, state.CredentialProviders, 3)

	require.NotNil(t, state.CredentialProviders[0].McpUserBasedAccessToken)
	assert.Equal(t, "https://aembit.io/revoke", state.CredentialProviders[0].McpUserBasedAccessToken.OAuthRevocationUrl.ValueString())

	require.NotNil(t, state.CredentialProviders[1].McpUserBasedAccessToken)
	assert.Equal(t, "", state.CredentialProviders[1].McpUserBasedAccessToken.OAuthRevocationUrl.ValueString())

	assert.Nil(t, state.CredentialProviders[2].McpUserBasedAccessToken)
	require.NotNil(t, state.CredentialProviders[2].OAuthAuthorizationCode)
}
