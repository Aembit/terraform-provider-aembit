// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure AembitProvider satisfies various provider interfaces.
var _ provider.Provider = &aembitProvider{}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &aembitProvider{
			version: version,
		}
	}
}

// aembitProviderModel maps provider schema data to a Go type.
type aembitProviderModel struct {
	Tenant      types.String `tfsdk:"tenant"`
	Token       types.String `tfsdk:"token"`
	StackDomain types.String `tfsdk:"stack_domain"`
	ClientID    types.String `tfsdk:"client_id"`
}

// AembitProvider defines the provider implementation.
type aembitProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *aembitProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "aembit"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *aembitProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"tenant": schema.StringAttribute{
				Description: "Tenant ID of the specific Aembit Cloud instance",
				Optional:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "The Aembit Trust Provider Client ID to use for authentication to the Aembit Cloud Tenant instance (recommended).",
				Optional:    true,
			},
			"token": schema.StringAttribute{
				Description: "Access Token to use for authentication to the Aembit Cloud Tenant instance",
				Optional:    true,
				Sensitive:   true,
			},
			"stack_domain": schema.StringAttribute{
				Description: "For development purposes only",
				Optional:    true,
			},
		},
	}
}

// Configure validators to ensure that only one credential provider type is specified.
func (p *aembitProvider) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("client_id"),
			path.MatchRoot("token"),
		),
	}
}

func (p *aembitProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Aembit client")

	// Retrieve provider data from configuration
	var config aembitProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.Tenant.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("tenant"),
			"Unknown Aembit API Tenant",
			"The provider cannot create the Aembit API client as there is an unknown configuration value for the Aembit API Tenant. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the AEMBIT_TENANT_ID environment variable.",
		)
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Aembit API Access Token",
			"The provider cannot create the Aembit API client as there is an unknown configuration value for the Aembit API Access Token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the AEMBIT_TOKEN environment variable.",
		)
	}

	if config.StackDomain.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("stack_domain"),
			"Unknown Aembit API Stack Domain",
			"The provider cannot create the Aembit API client as there is an unknown configuration value for the Aembit API Stack Domain. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the AEMBIT_STACK_DOMAIN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	tenant := os.Getenv("AEMBIT_TENANT_ID")
	token := os.Getenv("AEMBIT_TOKEN")
	stackDomain := os.Getenv("AEMBIT_STACK_DOMAIN")

	if !config.Tenant.IsNull() {
		tenant = config.Tenant.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if !config.StackDomain.IsNull() {
		stackDomain = config.StackDomain.ValueString()
	}

	// Check for the Aembit Client ID - if provided, then we need to try TrustProvider Attestation Authentication
	aembitClientID := os.Getenv("AEMBIT_CLIENT_ID")
	if len(aembitClientID) > 0 {
		idToken, err := getIdentityToken(aembitClientID, stackDomain)
		if err == nil {
			fmt.Printf("DEBUG: Got ID Token: %s\n", idToken)
			aembitToken, err := getAembitToken(aembitClientID, stackDomain, idToken)
			if err == nil {
				fmt.Printf("DEBUG: Got Aembit Token: %s\n", aembitToken)
				token = aembitToken
			} else {
				fmt.Printf("WARNING: Failed to get Aembit Token: %v\n", err)
			}
		} else {
			fmt.Printf("WARNING: Failed to get ID Token: %v\n", err)
		}
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if tenant == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("tenant"),
			"Missing Aembit API Tenant",
			"The provider cannot create the Aembit API client as there is a missing or empty value for the Aembit API Tenant. "+
				"Set the host value in the configuration or use the AEMBIT_TENANT_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Aembit API Access Token",
			"The provider cannot create the Aembit API client as there is a missing or empty value for the Aembit API Access Token. "+
				"Set the password value in the configuration or use the AEMBIT_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if stackDomain == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("stack_domain"),
			"Missing Aembit API Stack Domain",
			"The provider cannot create the Aembit API client as there is a missing or empty value for the Aembit API Stack Domain. "+
				"Set the password value in the configuration or use the AEMBIT_STACK_DOMAIN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "aembit_tenant", tenant)
	ctx = tflog.SetField(ctx, "aembit_token", token)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "aembit_token")

	tflog.Debug(ctx, "Creating Aembit client")

	// Create a new Aembit client using the configuration values
	client, err := aembit.NewClient(aembit.URLBuilder{}, &token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Aembit API Client",
			"An unexpected error occurred when creating the Aembit API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Aembit Client Error: "+err.Error(),
		)
		return
	}
	client.Tenant = tenant
	client.StackDomain = stackDomain

	// Make the Aembit client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Aembit client", map[string]any{"success": true})
}

func (p *aembitProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewServerWorkloadResource,
		NewCredentialProviderResource,
		NewTrustProviderResource,
		NewClientWorkloadResource,
		NewIntegrationResource,
		NewAccessConditionResource,
		NewAccessPolicyResource,
		NewAgentControllerResource,
	}
}

func (p *aembitProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewServerWorkloadsDataSource,
		NewCredentialProvidersDataSource,
		NewTrustProvidersDataSource,
		NewClientWorkloadsDataSource,
		NewAccessPoliciesDataSource,
	}
}

// // Temporary until Aembit SDK is published.
var GCP_ID_TOKEN string
var GITHUB_ID_TOKEN string
var TERRAFORM_ID_TOKEN string
var AEMBIT_TOKEN string

func getAembitToken(clientId, stackDomain, idToken string) (string, error) {
	if isTokenValid(AEMBIT_TOKEN) {
		return AEMBIT_TOKEN, nil
	}

	details := url.Values{}
	details.Set("grant_type", "client_credentials")
	details.Set("client_id", clientId)
	attestationData := map[string]interface{}{
		"version": "1.0.0",
		"github": map[string]interface{}{
			"identityToken": idToken,
		},
	}
	attestationJSON, err := json.Marshal(attestationData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal attestation data: %w", err)
	}
	details.Set("attestation", string(attestationJSON))

	tokenEndpoint := fmt.Sprintf("https://%s.id.%s/connect/token", getTenantId(clientId), stackDomain)
	req, err := http.NewRequest("POST", tokenEndpoint, bytes.NewBufferString(details.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch aembit token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	AEMBIT_TOKEN = tokenResponse.AccessToken
	return AEMBIT_TOKEN, nil
}

func getIdentityToken(clientId, stackDomain string) (string, error) {
	// First, determine which token type we need to get based on the identity type
	fmt.Printf("DEBUG: ID Token for ClientID: %s\n", clientId)
	switch getIdentityType((clientId)) {
	case "gcp_idtoken":
		return getGcpIdentityToken(clientId, stackDomain)
	case "github_idtoken":
		fmt.Printf("DEBUG: Getting GitHub ID Token\n")
		return getGitHubIdentityToken(clientId, stackDomain)
	case "terraform_idtoken":
	}
	return "", fmt.Errorf("no matching id token configuration")
}

func getGcpIdentityToken(clientId, stackDomain string) (string, error) {
	if isTokenValid(GCP_ID_TOKEN) {
		return GCP_ID_TOKEN, nil
	}

	audience := fmt.Sprintf("https://%s.id.%s", getTenantId(clientId), stackDomain)
	metadataIdentityTokenUrl := fmt.Sprintf("http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/identity?format=full&audience=%s", url.QueryEscape(audience))

	req, err := http.NewRequest("GET", metadataIdentityTokenUrl, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Metadata-Flavor", "Google")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch GCP ID Token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	GCP_ID_TOKEN = string(body)
	return GCP_ID_TOKEN, nil
}

func getGitHubIdentityToken(clientId, stackDomain string) (string, error) {
	if isTokenValid(GITHUB_ID_TOKEN) {
		return GITHUB_ID_TOKEN, nil
	}

	tokenRequestURL := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL")
	tokenRequestToken := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")
	if len(tokenRequestURL) == 0 || len(tokenRequestToken) == 0 {
		return "", fmt.Errorf("github action not configured for id_token access")
	}

	audience := fmt.Sprintf("https://%s.id.%s", getTenantId(clientId), stackDomain)
	identityTokenURL := fmt.Sprintf("%s&audience=%s", tokenRequestURL, url.QueryEscape(audience))

	req, err := http.NewRequest("GET", identityTokenURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenRequestToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch GCP ID Token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	jsonBody := make(map[string]string)
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		return "", fmt.Errorf("failed to parse response body: %w", err)
	}

	GITHUB_ID_TOKEN = jsonBody["value"]
	return GITHUB_ID_TOKEN, nil
}

func getTenantId(clientId string) string {
	clientIdSplit := strings.Split(clientId, ":")
	if len(clientIdSplit) >= 3 {
		return clientIdSplit[2]
	}

	return ""
}

func getIdentityType(clientId string) string {
	clientIdSplit := strings.Split(clientId, ":")
	if len(clientIdSplit) >= 5 {
		return clientIdSplit[4]
	}

	return ""
}

func isTokenValid(jwtToken string) bool {
	var payload []byte
	var expClaim float64
	var err error
	var ok bool

	if jwtToken == "" || !strings.Contains(jwtToken, ".") || strings.Count(jwtToken, ".") != 2 {
		return false
	}

	parts := strings.Split(jwtToken, ".")
	if payload, err = base64.RawURLEncoding.DecodeString(parts[1]); err != nil {
		return false
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return false
	}

	if expClaim, ok = claims["exp"].(float64); !ok {
		return false
	}

	// Calculate expiration with a 60-second safety window
	expiration := time.Unix(int64(expClaim), 0).UTC().Add(-60 * time.Second)
	return time.Now().UTC().Before(expiration)
}
