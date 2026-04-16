package provider

import (
	"context"
	"testing"

	"aembit.io/aembit"
	"terraform-provider-aembit/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestConvertToOAuthAuthorizationCodeDTO_IncludesFinalCallbackUrl(t *testing.T) {
	model := models.CredentialProviderResourceModel{
		Tags: types.MapNull(types.StringType),
		OAuthAuthorizationCode: &models.CredentialProviderOAuthAuthorizationCodeModel{
			OAuthDiscoveryUrl:    types.StringValue("https://aembit.io/.well-known/openid-configuration"),
			UserAuthorizationUrl: types.StringValue("https://aembit.io/user-authorize"),
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
				FinalCallbackUrl:      types.StringValue("https://aembit.io/final-callback"),
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

func TestConvertOAuthCodeDTOToModel_IncludesFinalCallbackUrl(t *testing.T) {
	dto := aembit.CredentialProviderV2DTO{
		ClientID:         "client-id",
		Scope:            "scope",
		CustomParameters: []aembit.CustomClaimsDTO{},
		CredentialOAuthAuthorizationCodeV2DTO: aembit.CredentialOAuthAuthorizationCodeV2DTO{
			AuthorizationUrl: "https://aembit.io/authorize",
			TokenUrl:         "https://aembit.io/token",
			IntrospectionUrl: "https://aembit.io/introspect",
			IsPkceRequired:   true,
			CallBackUrl:      "https://aembit.io/callback",
			FinalCallbackUrl: "https://aembit.io/final-callback",
		},
	}

	model := convertOAuthCodeDTOToModel(dto)

	if model.FinalCallbackUrl.ValueString() != "https://aembit.io/final-callback" {
		t.Fatalf(
			"expected final callback url to be mapped, got %q",
			model.FinalCallbackUrl.ValueString(),
		)
	}
}
