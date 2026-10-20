// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"testing"

	"terraform-provider-aembit/internal/provider/models"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
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
