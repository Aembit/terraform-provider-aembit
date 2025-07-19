package validators

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type oidcEndpointValidator struct{}

var _ validator.String = oidcEndpointValidator{}

func OidcEndpointValidation() validator.String {
	return oidcEndpointValidator{}
}

func (v oidcEndpointValidator) Description(ctx context.Context) string {
	return "OIDC Endpoint must not include well-known endpoint."
}

func (v oidcEndpointValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v oidcEndpointValidator) ValidateString(
	ctx context.Context,
	req validator.StringRequest,
	resp *validator.StringResponse,
) {
	if strings.Contains(req.ConfigValue.ValueString(), "/.well-known/openid-configuration") {
		resp.Diagnostics.AddError(
			"Invalid OIDC Endpoint",
			"OIDC Endpoint must not include well-known endpoint.",
		)
	}
}
