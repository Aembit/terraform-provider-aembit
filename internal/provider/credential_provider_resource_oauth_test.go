package provider

import (
	"context"
	"testing"

	"terraform-provider-aembit/internal/provider/models"

	"aembit.io/aembit"

	"github.com/hashicorp/terraform-plugin-framework/types"
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
		CredentialOAuthAuthorizationCodeV2DTO: aembit.CredentialOAuthAuthorizationCodeV2DTO{
			AuthorizationUrl: "https://aembit.io/authorize",
			TokenUrl:         "https://aembit.io/token",
			IntrospectionUrl: "https://aembit.io/introspect",
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
