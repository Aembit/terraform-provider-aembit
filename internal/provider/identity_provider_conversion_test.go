// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"terraform-provider-aembit/internal/provider/models"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertIdentityProviderModelToDTO_OidcUserInfoSettings(t *testing.T) {
	t.Parallel()

	model := models.IdentityProviderResourceModel{
		Name:        types.StringValue("Test OIDC IdP"),
		Description: types.StringValue("Test Description"),
		IsActive:    types.BoolValue(true),
		Tags:        types.MapNull(types.StringType),
		Oidc: &models.IdentityProviderOidcModel{
			OidcBaseUrl:               types.StringValue("https://oidc.example.com"),
			ClientId:                  types.StringValue("client-123"),
			Scopes:                    types.StringValue("openid profile email"),
			ClientSecret:              types.StringValue("super-secret"),
			AuthType:                  types.StringValue("ClientSecret"),
			PkceRequired:              types.BoolValue(true),
			EnableUserinfoGroupLookup: types.BoolValue(true),
			GroupClaimAttribute:       types.StringValue("user_groups"),
			GroupCacheTtlMinutes:      types.Int64Value(45),
		},
	}

	dto := convertIdentityProviderModelToDTO(
		context.Background(),
		model,
		nil,
		nil,
	)

	assert.Equal(t, "OIDCv1", dto.Type)
	assert.Equal(t, "Test OIDC IdP", dto.Name)
	assert.True(t, dto.EnableUserinfoGroupLookup)
	assert.Equal(t, "user_groups", dto.GroupClaimAttribute)
	assert.Equal(t, 45, dto.GroupCacheTtlMinutes)
}

func TestConvertIdentityProviderDTOToModel_OidcUserInfoSettings(t *testing.T) {
	t.Parallel()

	dto := aembit.IdentityProviderDTO{
		EntityDTO: aembit.EntityDTO{
			ExternalID:  "00000000-0000-0000-0000-000000000001",
			Name:        "Test OIDC IdP",
			Description: "Test Description",
			IsActive:    true,
		},
		Type: "OIDCv1",
		IdentityProviderOidcDTO: aembit.IdentityProviderOidcDTO{
			OidcBaseUrl:               "https://oidc.example.com",
			ClientId:                  "client-123",
			Scopes:                    "openid profile email",
			AuthType:                  "ClientSecret",
			PkceRequired:              true,
			AembitRedirectUrl:         "https://aembit.example.com/oauth/callback",
			AembitJwksUrl:             "https://aembit.example.com/.well-known/jwks.json",
			EnableUserinfoGroupLookup: true,
			GroupClaimAttribute:       "custom_groups",
			GroupCacheTtlMinutes:      60,
		},
	}

	planModel := &models.IdentityProviderResourceModel{
		Tags: types.MapNull(types.StringType),
		Oidc: &models.IdentityProviderOidcModel{
			ClientSecret: types.StringValue("retained-secret"),
		},
	}

	model := convertIdentityProviderDTOToModel(
		context.Background(),
		planModel,
		dto,
	)

	assert.NotNil(t, model.Oidc)
	assert.True(t, model.Oidc.EnableUserinfoGroupLookup.ValueBool())
	assert.Equal(t, "custom_groups", model.Oidc.GroupClaimAttribute.ValueString())
	assert.Equal(t, int64(60), model.Oidc.GroupCacheTtlMinutes.ValueInt64())
	assert.Equal(t, "retained-secret", model.Oidc.ClientSecret.ValueString())
}

func TestConvertIdentityProviderDTOToModel_OidcUserInfoDefaultsWhenOmitted(t *testing.T) {
	t.Parallel()

	// When enable_userinfo_group_lookup is false, backend omits groupClaimAttribute and groupCacheTtlMinutes (or returns zero-values "")
	dto := aembit.IdentityProviderDTO{
		EntityDTO: aembit.EntityDTO{
			ExternalID:  "00000000-0000-0000-0000-000000000002",
			Name:        "Test OIDC IdP Defaults",
			Description: "Test Description",
			IsActive:    true,
		},
		Type: "OIDCv1",
		IdentityProviderOidcDTO: aembit.IdentityProviderOidcDTO{
			OidcBaseUrl:               "https://oidc.example.com",
			ClientId:                  "client-456",
			Scopes:                    "openid profile email",
			AuthType:                  "ClientSecret",
			PkceRequired:              true,
			AembitRedirectUrl:         "https://aembit.example.com/oauth/callback",
			AembitJwksUrl:             "https://aembit.example.com/.well-known/jwks.json",
			EnableUserinfoGroupLookup: false,
			GroupClaimAttribute:       "",
			GroupCacheTtlMinutes:      0,
		},
	}

	// 1. Without plan values, falls back to schema defaults ("groups" and 15)
	planModelWithoutCustom := &models.IdentityProviderResourceModel{
		Tags: types.MapNull(types.StringType),
		Oidc: &models.IdentityProviderOidcModel{
			ClientSecret: types.StringValue("retained-secret"),
		},
	}

	modelDefault := convertIdentityProviderDTOToModel(
		context.Background(),
		planModelWithoutCustom,
		dto,
	)

	assert.NotNil(t, modelDefault.Oidc)
	assert.False(t, modelDefault.Oidc.EnableUserinfoGroupLookup.ValueBool())
	assert.Equal(t, "groups", modelDefault.Oidc.GroupClaimAttribute.ValueString())
	assert.Equal(t, int64(15), modelDefault.Oidc.GroupCacheTtlMinutes.ValueInt64())
	assert.Equal(t, "retained-secret", modelDefault.Oidc.ClientSecret.ValueString())

	// 2. With plan values, retains configured values to prevent persistent plan diff
	planModelWithCustom := &models.IdentityProviderResourceModel{
		Tags: types.MapNull(types.StringType),
		Oidc: &models.IdentityProviderOidcModel{
			ClientSecret:         types.StringValue("retained-secret"),
			GroupClaimAttribute:  types.StringValue("custom_roles"),
			GroupCacheTtlMinutes: types.Int64Value(30),
		},
	}

	modelCustom := convertIdentityProviderDTOToModel(
		context.Background(),
		planModelWithCustom,
		dto,
	)

	assert.NotNil(t, modelCustom.Oidc)
	assert.False(t, modelCustom.Oidc.EnableUserinfoGroupLookup.ValueBool())
	assert.Equal(t, "custom_roles", modelCustom.Oidc.GroupClaimAttribute.ValueString())
	assert.Equal(t, int64(30), modelCustom.Oidc.GroupCacheTtlMinutes.ValueInt64())
	assert.Equal(t, "retained-secret", modelCustom.Oidc.ClientSecret.ValueString())

	// 3. Import scenario: planModel is nil or planModel.Oidc is nil
	modelImport := convertIdentityProviderDTOToModel(
		context.Background(),
		&models.IdentityProviderResourceModel{},
		dto,
	)

	assert.NotNil(t, modelImport.Oidc)
	assert.Equal(t, "groups", modelImport.Oidc.GroupClaimAttribute.ValueString())
	assert.Equal(t, int64(15), modelImport.Oidc.GroupCacheTtlMinutes.ValueInt64())
	assert.True(t, modelImport.Oidc.ClientSecret.IsNull())
}

// TestConvertIdentityProviderModelToDTO_Prompt verifies conversion of prompt from Terraform model to DTO.
func TestConvertIdentityProviderModelToDTO_Prompt(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		modelPrompt    types.String
		expectedPrompt string
	}{
		{
			name:           "null prompt results in empty DTO prompt",
			modelPrompt:    types.StringNull(),
			expectedPrompt: "",
		},
		{
			name:           "unknown prompt results in empty DTO prompt",
			modelPrompt:    types.StringUnknown(),
			expectedPrompt: "",
		},
		{
			name:           "single value prompt is converted",
			modelPrompt:    types.StringValue("login"),
			expectedPrompt: "login",
		},
		{
			name:           "multi-value prompt is converted",
			modelPrompt:    types.StringValue("login consent"),
			expectedPrompt: "login consent",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			model := models.IdentityProviderResourceModel{
				Name:        types.StringValue("Test IdP"),
				Description: types.StringValue("Description"),
				IsActive:    types.BoolValue(true),
				Tags:        types.MapNull(types.StringType),
				Oidc: &models.IdentityProviderOidcModel{
					OidcBaseUrl:  types.StringValue("https://oidc.example.com"),
					ClientId:     types.StringValue("client-123"),
					Scopes:       types.StringValue("openid profile email"),
					AuthType:     types.StringValue("ClientSecret"),
					PkceRequired: types.BoolValue(true),
					Prompt:       tc.modelPrompt,
				},
			}

			dto := convertIdentityProviderModelToDTO(
				context.Background(),
				model,
				nil,
				nil,
			)

			assert.Equal(t, "OIDCv1", dto.Type)
			assert.Equal(t, tc.expectedPrompt, dto.Prompt)
		})
	}
}

// TestConvertIdentityProviderDTOToModel_Prompt verifies conversion of prompt from DTO to Terraform model.
func TestConvertIdentityProviderDTOToModel_Prompt(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		dtoPrompt      string
		expectedPrompt types.String
	}{
		{
			name:           "empty DTO prompt results in types.StringNull()",
			dtoPrompt:      "",
			expectedPrompt: types.StringNull(),
		},
		{
			name:           "single value prompt results in types.StringValue",
			dtoPrompt:      "login",
			expectedPrompt: types.StringValue("login"),
		},
		{
			name:           "multi-value prompt results in types.StringValue",
			dtoPrompt:      "login consent",
			expectedPrompt: types.StringValue("login consent"),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dto := aembit.IdentityProviderDTO{
				EntityDTO: aembit.EntityDTO{
					ExternalID:  "00000000-0000-0000-0000-000000000003",
					Name:        "Test IdP",
					Description: "Description",
					IsActive:    true,
				},
				Type: "OIDCv1",
				IdentityProviderOidcDTO: aembit.IdentityProviderOidcDTO{
					OidcBaseUrl:  "https://oidc.example.com",
					ClientId:     "client-123",
					Scopes:       "openid profile email",
					AuthType:     "ClientSecret",
					PkceRequired: true,
					Prompt:       tc.dtoPrompt,
				},
			}

			planModel := &models.IdentityProviderResourceModel{
				Tags: types.MapNull(types.StringType),
				Oidc: &models.IdentityProviderOidcModel{},
			}

			model := convertIdentityProviderDTOToModel(
				context.Background(),
				planModel,
				dto,
			)

			assert.NotNil(t, model.Oidc)
			assert.Equal(t, tc.expectedPrompt, model.Oidc.Prompt)
		})
	}
}

// TestIdentityProviderResourceSchema_OidcTokenUrlIsOptionalOnly verifies that oidc.token_url is Optional-only (not Computed).
func TestIdentityProviderResourceSchema_OidcTokenUrlIsOptionalOnly(t *testing.T) {
	t.Parallel()

	r := NewIdentityProviderResource()
	schemaResp := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), fwresource.SchemaRequest{}, schemaResp)

	oidcAttr, ok := schemaResp.Schema.Attributes["oidc"].(schema.SingleNestedAttribute)
	require.True(t, ok)

	tokenUrlAttr, ok := oidcAttr.Attributes["token_url"]
	require.True(t, ok)
	assert.True(t, tokenUrlAttr.IsOptional(), "token_url should be Optional")
	assert.False(t, tokenUrlAttr.IsComputed(), "token_url must not be Computed so removing it from config plans null")
}

// TestConvertIdentityProvider_OidcTokenUrlSetAndRemoved verifies model-to-DTO and DTO-to-model conversion when token_url is set and removed.
func TestConvertIdentityProvider_OidcTokenUrlSetAndRemoved(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// 1. Model with token_url set -> DTO has TokenUrl populated
	modelWithTokenUrl := models.IdentityProviderResourceModel{
		Name:        types.StringValue("Test OIDC IdP"),
		Description: types.StringValue("Test Description"),
		IsActive:    types.BoolValue(true),
		Tags:        types.MapNull(types.StringType),
		Oidc: &models.IdentityProviderOidcModel{
			OidcBaseUrl:               types.StringValue("https://oidc.example.com"),
			ClientId:                  types.StringValue("client-123"),
			Scopes:                    types.StringValue("openid profile email"),
			AuthType:                  types.StringValue("KeyPair"),
			PkceRequired:              types.BoolValue(true),
			EnableUserinfoGroupLookup: types.BoolValue(false),
			GroupClaimAttribute:       types.StringValue("groups"),
			GroupCacheTtlMinutes:      types.Int64Value(15),
			TokenUrl:                  types.StringValue("https://oidc.example.com/oauth2/token"),
		},
	}

	dtoWithTokenUrl := convertIdentityProviderModelToDTO(ctx, modelWithTokenUrl, nil, nil)
	assert.Equal(t, "https://oidc.example.com/oauth2/token", dtoWithTokenUrl.TokenUrl)

	// 2. Model with token_url removed (null) -> DTO has empty TokenUrl ("")
	modelWithoutTokenUrl := modelWithTokenUrl
	oidcCopy := *modelWithTokenUrl.Oidc
	oidcCopy.TokenUrl = types.StringNull()
	modelWithoutTokenUrl.Oidc = &oidcCopy

	dtoWithoutTokenUrl := convertIdentityProviderModelToDTO(ctx, modelWithoutTokenUrl, nil, nil)
	assert.Empty(t, dtoWithoutTokenUrl.TokenUrl)

	// 3. Backend response with empty TokenUrl -> Model has TokenUrl = StringNull()
	dtoResponse := aembit.IdentityProviderDTO{
		EntityDTO: aembit.EntityDTO{
			ExternalID:  "00000000-0000-0000-0000-000000000003",
			Name:        "Test OIDC IdP",
			Description: "Test Description",
			IsActive:    true,
		},
		Type: "OIDCv1",
		IdentityProviderOidcDTO: aembit.IdentityProviderOidcDTO{
			OidcBaseUrl:  "https://oidc.example.com",
			ClientId:     "client-123",
			Scopes:       "openid profile email",
			AuthType:     "KeyPair",
			PkceRequired: true,
			TokenUrl:     "",
		},
	}

	convertedModel := convertIdentityProviderDTOToModel(ctx, &modelWithoutTokenUrl, dtoResponse)
	require.NotNil(t, convertedModel.Oidc)
	assert.True(t, convertedModel.Oidc.TokenUrl.IsNull())
}

func TestConvertIdentityProviderDTOToModel_AembitClientId(t *testing.T) {
	t.Parallel()

	t.Run("populated aembit client id mapped to string value", func(t *testing.T) {
		dto := aembit.IdentityProviderDTO{
			EntityDTO: aembit.EntityDTO{
				ExternalID: "00000000-0000-0000-0000-000000000003",
				Name:       "OIDC with ClientId",
				IsActive:   true,
			},
			Type: "OIDCv1",
			IdentityProviderOidcDTO: aembit.IdentityProviderOidcDTO{
				OidcBaseUrl:    "https://oidc.example.com",
				ClientId:       "client-789",
				Scopes:         "openid",
				AuthType:       "ClientSecret",
				AembitClientId: "arn:aembit:iam:tenant:client/test-client-id",
			},
		}

		model := convertIdentityProviderDTOToModel(
			context.Background(),
			&models.IdentityProviderResourceModel{},
			dto,
		)

		assert.NotNil(t, model.Oidc)
		assert.False(t, model.Oidc.AembitClientId.IsNull())
		assert.Equal(t, "arn:aembit:iam:tenant:client/test-client-id", model.Oidc.AembitClientId.ValueString())
	})

	t.Run("empty aembit client id mapped to string null", func(t *testing.T) {
		dto := aembit.IdentityProviderDTO{
			EntityDTO: aembit.EntityDTO{
				ExternalID: "00000000-0000-0000-0000-000000000004",
				Name:       "OIDC without ClientId",
				IsActive:   true,
			},
			Type: "OIDCv1",
			IdentityProviderOidcDTO: aembit.IdentityProviderOidcDTO{
				OidcBaseUrl:    "https://oidc.example.com",
				ClientId:       "client-789",
				Scopes:         "openid",
				AuthType:       "ClientSecret",
				AembitClientId: "",
			},
		}

		model := convertIdentityProviderDTOToModel(
			context.Background(),
			&models.IdentityProviderResourceModel{},
			dto,
		)

		assert.NotNil(t, model.Oidc)
		assert.True(t, model.Oidc.AembitClientId.IsNull())
	})

	t.Run("saml provider leaves oidc model nil", func(t *testing.T) {
		dto := aembit.IdentityProviderDTO{
			EntityDTO: aembit.EntityDTO{
				ExternalID: "00000000-0000-0000-0000-000000000005",
				Name:       "SAML IdP",
				IsActive:   true,
			},
			Type: "SAMLv2",
			IdentityProviderSamlDTO: aembit.IdentityProviderSamlDTO{
				MetadataUrl: "https://saml.example.com/metadata",
			},
		}

		model := convertIdentityProviderDTOToModel(
			context.Background(),
			&models.IdentityProviderResourceModel{},
			dto,
		)

		assert.Nil(t, model.Oidc)
		assert.NotNil(t, model.Saml)
	})
}

// mockIdentityProviderTestProvider embeds the real provider while injecting an isolated CloudClient for unit testing.
type mockIdentityProviderTestProvider struct {
	// Provider delegates Metadata, Schema, Resources, and DataSources to the standard Aembit provider.
	provider.Provider
	// client is the isolated mock-backed CloudClient supplied to resources and data sources.
	client *aembit.CloudClient
}

// Configure injects the isolated mock CloudClient without reading environment variables or global transport state.
func (p *mockIdentityProviderTestProvider) Configure(
	_ context.Context,
	_ provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {
	resp.DataSourceData = p.client
	resp.ResourceData = p.client
}

// TestIdentityProviderResource_OidcPromptSchema verifies that oidc.prompt is Optional-only (not Computed)
// on the managed resource while remaining Computed on the data source, ensuring omitted HCL config plans null.
func TestIdentityProviderResource_OidcPromptSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	res := NewIdentityProviderResource()
	var resp fwresource.SchemaResponse
	res.Schema(ctx, fwresource.SchemaRequest{}, &resp)
	require.False(t, resp.Diagnostics.HasError(), "unexpected schema diagnostics: %v", resp.Diagnostics)

	oidcAttr, ok := resp.Schema.Attributes["oidc"]
	require.True(t, ok, "missing oidc attribute in identity_provider schema")

	oidcNested, ok := oidcAttr.(schema.SingleNestedAttribute)
	require.True(t, ok, "expected oidc to be schema.SingleNestedAttribute")

	promptAttr, ok := oidcNested.Attributes["prompt"]
	require.True(t, ok, "missing prompt attribute in oidc schema")

	promptStrAttr, ok := promptAttr.(schema.StringAttribute)
	require.True(t, ok, "expected oidc.prompt to be schema.StringAttribute")

	assert.True(t, promptStrAttr.Optional, "oidc.prompt must be Optional")
	assert.False(t, promptStrAttr.Computed, "oidc.prompt must not be Computed so omitted config clears prior state")
	assert.Nil(t, promptStrAttr.Default, "oidc.prompt must not define a static default")
	assert.Len(t, promptStrAttr.Validators, 1, "oidc.prompt must retain OidcPromptValidation")

	// Also verify Protocol v6 schema advertised to Terraform Core and data source Computed contract.
	server, err := providerserver.NewProtocol6WithError(New("test", "unittest")())()
	require.NoError(t, err)

	protoSchema, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	require.NoError(t, err)
	require.Empty(t, protoSchema.Diagnostics)

	resProtoSchema, ok := protoSchema.ResourceSchemas["aembit_identity_provider"]
	require.True(t, ok)

	var foundPrompt *tfprotov6.SchemaAttribute
	for _, attr := range resProtoSchema.Block.Attributes {
		if attr.Name == "oidc" && attr.NestedType != nil {
			for _, nestedAttr := range attr.NestedType.Attributes {
				if nestedAttr.Name == "prompt" {
					foundPrompt = nestedAttr
					break
				}
			}
		}
	}
	require.NotNil(t, foundPrompt, "missing oidc.prompt in protocol v6 resource schema")
	assert.True(t, foundPrompt.Optional)
	assert.False(t, foundPrompt.Computed)

	dsProtoSchema, ok := protoSchema.DataSourceSchemas["aembit_identity_providers"]
	require.True(t, ok)

	var dsFoundPrompt *tfprotov6.SchemaAttribute
	for _, attr := range dsProtoSchema.Block.Attributes {
		if attr.Name == "identity_providers" && attr.NestedType != nil {
			for _, idpAttr := range attr.NestedType.Attributes {
				if idpAttr.Name == "oidc" && idpAttr.NestedType != nil {
					for _, oidcAttr := range idpAttr.NestedType.Attributes {
						if oidcAttr.Name == "prompt" {
							dsFoundPrompt = oidcAttr
							break
						}
					}
				}
			}
		}
	}
	require.NotNil(t, dsFoundPrompt, "missing oidc.prompt in protocol v6 data source schema")
	assert.True(t, dsFoundPrompt.Computed)
	assert.False(t, dsFoundPrompt.Optional)
	assert.False(t, dsFoundPrompt.Required)
}

// newMockIdpCloudClient creates an isolated mock-backed CloudClient and captures PUT payloads for assertions.
func newMockIdpCloudClient(t *testing.T, externalID string) (*aembit.CloudClient, func() []aembit.IdentityProviderDTO) {
	t.Helper()

	var (
		mu          sync.Mutex
		storedIdp   *aembit.IdentityProviderDTO
		putPayloads []aembit.IdentityProviderDTO
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
					var incoming aembit.IdentityProviderDTO
					if unmarshalErr := json.Unmarshal(bodyBytes, &incoming); unmarshalErr != nil {
						return nil, unmarshalErr
					}
					incoming.ExternalID = externalID
					incoming.AembitRedirectUrl = "https://unit-tenant.aembit.io/callback"
					incoming.AembitJwksUrl = "https://unit-tenant.aembit.io/jwks"
					incoming.AembitClientId = "aembit-client-1"
					incoming.ClientSecret = ""
					storedCopy := incoming
					storedIdp = &storedCopy

					respBytes, marshalErr := json.Marshal(storedIdp)
					if marshalErr != nil {
						return nil, marshalErr
					}
					return &http.Response{
						StatusCode: http.StatusCreated,
						Body:       io.NopCloser(bytes.NewReader(respBytes)),
						Header:     make(http.Header),
					}, nil

				case http.MethodGet:
					if storedIdp == nil {
						return &http.Response{
							StatusCode: http.StatusNoContent,
							Body:       io.NopCloser(bytes.NewReader(nil)),
							Header:     make(http.Header),
						}, nil
					}
					respBytes, marshalErr := json.Marshal(storedIdp)
					if marshalErr != nil {
						return nil, marshalErr
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(respBytes)),
						Header:     make(http.Header),
					}, nil

				case http.MethodPut:
					bodyBytes, readErr := io.ReadAll(req.Body)
					if readErr != nil {
						return nil, readErr
					}
					var incoming aembit.IdentityProviderDTO
					if unmarshalErr := json.Unmarshal(bodyBytes, &incoming); unmarshalErr != nil {
						return nil, unmarshalErr
					}
					putPayloads = append(putPayloads, incoming)

					incoming.ExternalID = externalID
					incoming.AembitRedirectUrl = "https://unit-tenant.aembit.io/callback"
					incoming.AembitJwksUrl = "https://unit-tenant.aembit.io/jwks"
					incoming.AembitClientId = "aembit-client-1"
					incoming.ClientSecret = ""
					storedCopy := incoming
					storedIdp = &storedCopy

					respBytes, marshalErr := json.Marshal(storedIdp)
					if marshalErr != nil {
						return nil, marshalErr
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(respBytes)),
						Header:     make(http.Header),
					}, nil

				case http.MethodPatch:
					if storedIdp != nil {
						storedIdp.IsActive = false
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
						Header:     make(http.Header),
					}, nil

				case http.MethodDelete:
					storedIdp = nil
					return &http.Response{
						StatusCode: http.StatusNoContent,
						Body:       io.NopCloser(bytes.NewReader(nil)),
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

	getPuts := func() []aembit.IdentityProviderDTO {
		mu.Lock()
		defer mu.Unlock()
		out := make([]aembit.IdentityProviderDTO, len(putPayloads))
		copy(out, putPayloads)
		return out
	}

	return mockClient, getPuts
}

// TestIdentityProviderResource_UpdateClearsOidcPrompt verifies directly against identityProviderResource
// Create, Update, and Read that transitioning oidc.prompt from a value to null clears it in both API and state.
func TestIdentityProviderResource_UpdateClearsOidcPrompt(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	const externalID = "11111111-2222-3333-4444-555555555555"

	mockClient, getPuts := newMockIdpCloudClient(t, externalID)

	res, ok := NewIdentityProviderResource().(*identityProviderResource)
	require.True(t, ok)
	res.client = mockClient

	var schemaResp fwresource.SchemaResponse
	res.Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)
	require.False(t, schemaResp.Diagnostics.HasError())
	resSchema := schemaResp.Schema

	buildModel := func(id types.String, prompt types.String) models.IdentityProviderResourceModel {
		return models.IdentityProviderResourceModel{
			ID:          id,
			Name:        types.StringValue("OIDC Prompt Test"),
			Description: types.StringValue("Testing prompt lifecycle"),
			IsActive:    types.BoolValue(true),
			Tags:        types.MapNull(types.StringType),
			TagsAll:     types.MapNull(types.StringType),
			Oidc: &models.IdentityProviderOidcModel{
				OidcBaseUrl:               types.StringValue("https://auth.example.com"),
				ClientId:                  types.StringValue("client-1"),
				Scopes:                    types.StringValue("openid profile email"),
				ClientSecret:              types.StringValue("secret-1"),
				AuthType:                  types.StringValue("ClientSecret"),
				PkceRequired:              types.BoolValue(true),
				AembitRedirectUrl:         types.StringUnknown(),
				AembitJwksUrl:             types.StringUnknown(),
				EnableUserinfoGroupLookup: types.BoolValue(false),
				GroupClaimAttribute:       types.StringValue("groups"),
				GroupCacheTtlMinutes:      types.Int64Value(15),
				TokenUrl:                  types.StringNull(),
				Prompt:                    prompt,
				AembitClientId:            types.StringUnknown(),
			},
		}
	}

	// 1. Create with prompt = "login".
	createPlan := tfsdk.Plan{Schema: resSchema}
	require.False(t, createPlan.Set(ctx, buildModel(types.StringUnknown(), types.StringValue("login"))).HasError())

	createResp := fwresource.CreateResponse{
		State: tfsdk.State{Schema: resSchema, Raw: tftypes.NewValue(resSchema.Type().TerraformType(ctx), nil)},
	}
	res.Create(ctx, fwresource.CreateRequest{Plan: createPlan}, &createResp)
	require.False(t, createResp.Diagnostics.HasError(), "Create diagnostics: %v", createResp.Diagnostics)

	var stateAfterCreate models.IdentityProviderResourceModel
	require.False(t, createResp.State.Get(ctx, &stateAfterCreate).HasError())
	require.NotNil(t, stateAfterCreate.Oidc)
	assert.Equal(t, types.StringValue("login"), stateAfterCreate.Oidc.Prompt)

	// 2. Update with prompt = null (removed from HCL config).
	updatePlan := tfsdk.Plan{Schema: resSchema}
	require.False(t, updatePlan.Set(ctx, buildModel(types.StringValue(externalID), types.StringNull())).HasError())

	updateResp := fwresource.UpdateResponse{
		State: createResp.State,
	}
	res.Update(ctx, fwresource.UpdateRequest{
		State: createResp.State,
		Plan:  updatePlan,
	}, &updateResp)
	require.False(t, updateResp.Diagnostics.HasError(), "Update diagnostics: %v", updateResp.Diagnostics)

	var stateAfterClear models.IdentityProviderResourceModel
	require.False(t, updateResp.State.Get(ctx, &stateAfterClear).HasError())
	require.NotNil(t, stateAfterClear.Oidc)
	assert.True(t, stateAfterClear.Oidc.Prompt.IsNull(), "expected oidc.prompt to be null after clearing")

	puts := getPuts()
	require.Len(t, puts, 1)
	assert.Empty(t, puts[0].Prompt, "PUT payload must omit/clear prompt when planned prompt is null")

	// 3. Read refreshes state and preserves null prompt.
	readResp := fwresource.ReadResponse{State: updateResp.State}
	res.Read(ctx, fwresource.ReadRequest{State: updateResp.State}, &readResp)
	require.False(t, readResp.Diagnostics.HasError(), "Read diagnostics: %v", readResp.Diagnostics)

	var stateAfterRead models.IdentityProviderResourceModel
	require.False(t, readResp.State.Get(ctx, &stateAfterRead).HasError())
	require.NotNil(t, stateAfterRead.Oidc)
	assert.True(t, stateAfterRead.Oidc.Prompt.IsNull())
}

// skipIfCannotCreateDirectorySymlink skips Terraform CLI harness tests on Windows hosts lacking symlink privileges.
func skipIfCannotCreateDirectorySymlink(t *testing.T) {
	t.Helper()
	tmpDir := t.TempDir()
	srcDir := filepath.Join(tmpDir, "src")
	require.NoError(t, os.Mkdir(srcDir, 0o750))
	if err := os.Symlink(srcDir, filepath.Join(tmpDir, "dst")); err != nil {
		t.Skipf("skipping Terraform CLI harness test because directory symlinks are unavailable: %v", err)
	}
}

// TestUnitIdentityProviderResource_OidcPromptClear verifies end-to-end with Terraform CLI that
// configuring oidc.prompt, updating it, and removing it from HCL sends an update clearing prompt to null.
func TestUnitIdentityProviderResource_OidcPromptClear(t *testing.T) {
	t.Parallel()
	skipIfCannotCreateDirectorySymlink(t)

	const externalID = "11111111-2222-3333-4444-555555555555"
	mockClient, getPuts := newMockIdpCloudClient(t, externalID)

	factories := map[string]func() (tfprotov6.ProviderServer, error){
		"aembit": providerserver.NewProtocol6WithError(&mockIdentityProviderTestProvider{
			Provider: New("test", "unittest")(),
			client:   mockClient,
		}),
	}

	const resourceName = "aembit_identity_provider.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "aembit" {
	tenant = "unit-tenant"
	token  = "unit-test-token"
}

resource "aembit_identity_provider" "test" {
	name        = "OIDC Prompt Test"
	description = "Testing prompt lifecycle"
	is_active   = true
	oidc = {
		oidc_base_url = "https://auth.example.com"
		client_id     = "client-1"
		scopes        = "openid profile email"
		auth_type     = "ClientSecret"
		client_secret = "secret-1"
		prompt        = "login"
	}
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", externalID),
					resource.TestCheckResourceAttr(resourceName, "oidc.prompt", "login"),
				),
			},
			{
				Config: `
provider "aembit" {
	tenant = "unit-tenant"
	token  = "unit-test-token"
}

resource "aembit_identity_provider" "test" {
	name        = "OIDC Prompt Test"
	description = "Testing prompt lifecycle"
	is_active   = true
	oidc = {
		oidc_base_url = "https://auth.example.com"
		client_id     = "client-1"
		scopes        = "openid profile email"
		auth_type     = "ClientSecret"
		client_secret = "secret-1"
		prompt        = "none"
	}
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", externalID),
					resource.TestCheckResourceAttr(resourceName, "oidc.prompt", "none"),
				),
			},
			{
				Config: `
provider "aembit" {
	tenant = "unit-tenant"
	token  = "unit-test-token"
}

resource "aembit_identity_provider" "test" {
	name        = "OIDC Prompt Test"
	description = "Testing prompt lifecycle"
	is_active   = true
	oidc = {
		oidc_base_url = "https://auth.example.com"
		client_id     = "client-1"
		scopes        = "openid profile email"
		auth_type     = "ClientSecret"
		client_secret = "secret-1"
	}
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", externalID),
					resource.TestCheckNoResourceAttr(resourceName, "oidc.prompt"),
				),
			},
		},
	})

	puts := getPuts()
	require.Len(t, puts, 2, "expected 2 PUT updates (change to 'none', then clear to empty/null)")
	assert.Equal(t, "none", puts[0].Prompt)
	assert.Empty(t, puts[1].Prompt, "removing prompt from HCL must send empty/omitted prompt in PUT payload")
}
