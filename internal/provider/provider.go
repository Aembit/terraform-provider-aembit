// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"context"
	"crypto/tls"
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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	ClientID    types.String `tfsdk:"client_id"`
	ResourceSet types.String `tfsdk:"resource_set_id"`
}

// AembitProvider defines the provider implementation.
type aembitProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Configure adds the provider configured client to the resource.
func resourceConfigure(
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) *aembit.CloudClient {
	if req.ProviderData == nil {
		return nil
	}

	var client *aembit.CloudClient
	var ok bool
	if client, ok = req.ProviderData.(*aembit.CloudClient); !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configuration Type",
			fmt.Sprintf(
				"Expected *aembit.CloudClient, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return nil
	}

	return client
}

// Configure adds the provider configured client to the data source.
func datasourceConfigure(
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) *aembit.CloudClient {
	if req.ProviderData == nil {
		return nil
	}

	var client *aembit.CloudClient
	var ok bool
	if client, ok = req.ProviderData.(*aembit.CloudClient); !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configuration Type",
			fmt.Sprintf(
				"Expected *aembit.CloudClient, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return nil
	}

	return client
}

// Metadata returns the provider type name.
func (p *aembitProvider) Metadata(
	_ context.Context,
	_ provider.MetadataRequest,
	resp *provider.MetadataResponse,
) {
	resp.TypeName = "aembit"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *aembitProvider) Schema(
	_ context.Context,
	_ provider.SchemaRequest,
	resp *provider.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"tenant": schema.StringAttribute{
				Description: "Tenant ID of the specific Aembit Cloud instance.",
				Optional:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "The Aembit Trust Provider Client ID to use for authentication to the Aembit Cloud Tenant instance (recommended).",
				Optional:    true,
			},
			"token": schema.StringAttribute{
				Description: "Access Token to use for authentication to the Aembit Cloud Tenant instance.",
				Optional:    true,
				Sensitive:   true,
			},
			"resource_set_id": schema.StringAttribute{
				Description: "The Aembit Resource Set to use for resources associated with this Terraform Provider.",
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

func (p *aembitProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {
	var err error
	tflog.Info(ctx, "Configuring Aembit client...")

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

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	tenant := os.Getenv("AEMBIT_TENANT_ID")
	token := os.Getenv("AEMBIT_TOKEN")
	stackDomain := os.Getenv("AEMBIT_STACK_DOMAIN")
	resourceSetId := os.Getenv("AEMBIT_RESOURCE_SET_ID")
	if len(stackDomain) == 0 {
		stackDomain = "useast2.aembit.io"
	}

	if !config.Tenant.IsNull() && len(config.Tenant.ValueString()) > 0 {
		tenant = config.Tenant.ValueString()
	}
	if !config.Token.IsNull() && len(config.Token.ValueString()) > 0 {
		token = config.Token.ValueString()
	}
	if !config.ResourceSet.IsNull() && len(config.ResourceSet.ValueString()) > 0 {
		resourceSetId = config.ResourceSet.ValueString()
	}

	// Check for the Aembit Client ID - if provided, then we need to try TrustProvider Attestation Authentication
	aembitClientID := os.Getenv("AEMBIT_CLIENT_ID")
	if !config.ClientID.IsNull() && len(config.ClientID.ValueString()) > 0 {
		// If there is a provider block ClientID, prefer that
		aembitClientID = config.ClientID.ValueString()
	}
	if len(aembitClientID) > 0 {
		tenant = getAembitTenantId(aembitClientID)
		tflog.Debug(
			ctx,
			"Using Aembit Native Authentication",
			map[string]interface{}{"tenantId": tenant},
		)

		// Try with the resourceSetId first
		if token, err = getToken(ctx, aembitClientID, stackDomain, resourceSetId, p.version); err != nil {
			tflog.Warn(ctx, "Failed to get Aembit Auth Token", map[string]interface{}{"error": err})
		} else {
			tflog.Debug(ctx, "Retrieved Aembit Auth Token for ResourceSet", map[string]interface{}{"resourceSetId": resourceSetId})
		}

		// If there was an error, try again without the resourceSetId
		// LEGACY: This is included to authenticate to the default resource set if resource_set_id is specified in the provider
		if token == "" {
			if token, err = getToken(ctx, aembitClientID, stackDomain, "", p.version); err != nil {
				tflog.Warn(
					ctx,
					"Failed to get Aembit Auth Token",
					map[string]interface{}{"error": err},
				)
			}
		} else {
			tflog.Debug(ctx, "Retrieved Aembit Auth Token for Default ResourceSet")
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

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "aembit_tenant", tenant)
	ctx = tflog.SetField(ctx, "aembit_token", token)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "aembit_token")

	tflog.Debug(ctx, "Creating Aembit client")

	// Create a new Aembit client using the configuration values
	client, err := aembit.NewClient(aembit.URLBuilder{}, &token, resourceSetId, p.version)
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

	tflog.Info(
		ctx,
		fmt.Sprintf("Configured Aembit client (%s)", p.version),
		map[string]any{"success": true},
	)
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
		NewRoleResource,
		NewSignInPolicyResource,
		NewStandaloneCertificateAuthorityResource,
		NewCredentialProviderIntegrationResource,
		NewDiscoveryIntegrationResource,
		NewLogStreamResource,
		// NewResourceSetResource,	// Preventing Resource Set Resources via Terraform until we add support for deleting Resource Sets
		NewGlobalPolicyComplianceResource,
		NewIdentityProviderResource,
	}
}

func (p *aembitProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewServerWorkloadsDataSource,
		NewCredentialProvidersDataSource,
		NewTrustProvidersDataSource,
		NewClientWorkloadsDataSource,
		NewIntegrationsDataSource,
		NewAccessConditionsDataSource,
		NewAccessPoliciesDataSource,
		NewAgentControllersDataSource,
		NewAgentControllerDeviceCodeDataSource,
		NewRolesDataSource,
		NewResourceSetDataSource,
		NewResourceSetsDataSource,
		NewStandaloneCertificateAuthoritiesDataSource,
		NewCountriesDataSource,
		NewTimeZonesDataSource,
		NewCredentialProviderIntegrationsDataSource,
		NewDiscoveryIntegrationsDataSource,
		NewGlobalPolicyComplianceDataSource,
		NewLogStreamsDataSource,
		NewCallerIdentityDataSource,
		NewIdentityProviderDataSource,
	}
}

// // Temporary until Aembit SDK is published.
var (
	GCP_ID_TOKEN       string
	GITHUB_ID_TOKEN    string
	TERRAFORM_ID_TOKEN string
	AEMBIT_TOKEN       string
)

type ClientRequestNetwork struct {
	TargetHost        string `json:"targetHost"`
	TargetPort        int16  `json:"targetPort"`
	TransportProtocol string `json:"transportProtocol"`
}

type ClientRequest struct {
	Version string               `json:"version"`
	Network ClientRequestNetwork `json:"network"`
}

type WorkloadAssessmentIdToken struct {
	IdentityToken string `json:"identityToken"`
}

type WorkloadAssessment struct {
	Version   string                    `json:"version"`
	GCP       WorkloadAssessmentIdToken `json:"gcp,omitempty"`
	GitHub    WorkloadAssessmentIdToken `json:"github,omitempty"`
	Terraform WorkloadAssessmentIdToken `json:"terraform,omitempty"`
	OS        OperatingSystemData       `json:"os,omitempty"`
}

type OperatingSystemData struct {
	Environment EnvironmentVariablesData `json:"environment,omitempty"`
}

type EnvironmentVariablesData struct {
	ResourceSet string `json:"AEMBIT_RESOURCE_SET_ID,omitempty"`
}

type tokenAuth struct {
	token string
}

func (t tokenAuth) GetRequestMetadata(
	ctx context.Context,
	in ...string,
) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + t.token,
	}, nil
}

func (tokenAuth) RequireTransportSecurity() bool {
	return true
}

func getToken(
	ctx context.Context,
	aembitClientID, stackDomain, resourceSetId, version string,
) (string, error) {
	idToken, err := getIdentityToken(aembitClientID, stackDomain)
	if err == nil {
		aembitToken, err := getAembitToken(
			aembitClientID,
			stackDomain,
			idToken,
			resourceSetId,
			version,
		)
		if err == nil {
			roleToken, err := getAembitCredential(
				fmt.Sprintf("%s.api.%s", getAembitTenantId(aembitClientID), stackDomain),
				443,
				aembitClientID,
				stackDomain,
				idToken,
				aembitToken,
				resourceSetId,
			)
			if err == nil {
				return roleToken, nil
			} else {
				tflog.Warn(ctx, "Failed to get Aembit API Role Token: %v", map[string]interface{}{
					"error": err,
				})
				return "", err
			}
		} else {
			tflog.Warn(ctx, "Failed to get Aembit Token: %v", map[string]interface{}{
				"error": err,
			})
			return "", err
		}
	} else {
		tflog.Warn(ctx, "Failed to get ID Token: %v", map[string]interface{}{
			"error": err,
		})
		return "", err
	}
}

func getAembitCredential(
	targetHost string,
	targetPort int16,
	clientId, stackDomain, idToken, aembitToken, resourceSetId string,
) (string, error) {
	var err error
	var clientRequest, workloadAssessment string
	var conn *grpc.ClientConn
	var aembitClient EdgeCommanderClient
	var credResponse *CredentialResponse

	tlsCreds := credentials.NewTLS(
		&tls.Config{InsecureSkipVerify: false, MinVersion: tls.VersionTLS12},
	)
	if conn, err = grpc.NewClient(fmt.Sprintf("%s.ec.%s:443", getAembitTenantId(clientId), stackDomain), grpc.WithTransportCredentials(tlsCreds), grpc.WithPerRPCCredentials(tokenAuth{token: aembitToken})); err != nil {
		return "", err
	}
	defer conn.Close()

	if clientRequest, err = getClientRequest(targetHost, targetPort); err != nil {
		return "", err
	}
	if workloadAssessment, err = getWorkloadAssessment(clientId, idToken, resourceSetId); err != nil {
		return "", err
	}

	aembitClient = NewEdgeCommanderClient(conn)
	if credResponse, err = aembitClient.GetCredential(context.Background(), &CredentialRequest{
		ClientRequest:      clientRequest,
		AgentAssessment:    workloadAssessment,
		WorkloadAssessment: workloadAssessment,
	}); err != nil {
		return "", err
	}

	return credResponse.Credential, nil
}

func getClientRequest(targetHost string, targetPort int16) (string, error) {
	var request []byte
	var err error
	clientRequest := ClientRequest{
		Version: "1.0.0",
		Network: ClientRequestNetwork{
			TargetHost:        targetHost,
			TargetPort:        targetPort,
			TransportProtocol: "TCP",
		},
	}

	if request, err = json.Marshal(clientRequest); err != nil {
		return "", err
	}
	return string(request), nil
}

func getWorkloadAssessment(clientId, idToken, resourceSetId string) (string, error) {
	var assessment []byte
	var err error
	var workload WorkloadAssessment

	switch getAembitIdentityType(clientId) {
	case "gcp_idtoken":
		workload = WorkloadAssessment{
			Version: "1.0.0",
			GCP:     WorkloadAssessmentIdToken{IdentityToken: idToken},
		}
	case "github_idtoken":
		workload = WorkloadAssessment{
			Version: "1.0.0",
			GitHub:  WorkloadAssessmentIdToken{IdentityToken: idToken},
		}
	case "terraform_idtoken":
		workload = WorkloadAssessment{
			Version:   "1.0.0",
			Terraform: WorkloadAssessmentIdToken{IdentityToken: idToken},
		}
	default:
		return "", fmt.Errorf("invalid aembit client id")
	}

	// Add the Aembit Resource Set if it's been configured in the environment
	if len(resourceSetId) > 0 {
		workload.OS = OperatingSystemData{
			Environment: EnvironmentVariablesData{
				ResourceSet: resourceSetId,
			},
		}
	}

	if assessment, err = json.Marshal(workload); err != nil {
		return "", err
	}
	return string(assessment), nil
}

func getAembitToken(clientId, stackDomain, idToken, resourceSetId, version string) (string, error) {
	if isTokenValid(AEMBIT_TOKEN) {
		return AEMBIT_TOKEN, nil
	}

	idTokenType := ""
	switch getAembitIdentityType((clientId)) {
	case "gcp_idtoken":
		idTokenType = "gcp"
	case "github_idtoken":
		idTokenType = "github"
	case "terraform_idtoken":
		idTokenType = "terraform"
	default:
		return "", fmt.Errorf("invalid aembit client id")
	}

	details := url.Values{}
	details.Set("grant_type", "client_credentials")
	details.Set("client_id", clientId)

	attestationData := map[string]interface{}{
		"version": "1.0.0",
		idTokenType: map[string]interface{}{
			"identityToken": idToken,
		},
	}
	attestationJSON, err := json.Marshal(attestationData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal attestation data: %w", err)
	}
	details.Set("attestation", string(attestationJSON))

	tokenEndpoint := fmt.Sprintf(
		"https://%s.id.%s/connect/token",
		getAembitTenantId(clientId),
		stackDomain,
	)
	req, err := http.NewRequest("POST", tokenEndpoint, bytes.NewBufferString(details.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("User-Agent", fmt.Sprintf("AembitTerraformProvider/%s", version))

	// Add the Aembit Resource Set if it's been configured in the environment
	if len(resourceSetId) > 0 {
		req.Header.Set("X-Aembit-ResourceSet", resourceSetId)
	}

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
	switch getAembitIdentityType((clientId)) {
	case "gcp_idtoken":
		return getGcpIdentityToken(clientId, stackDomain)
	case "github_idtoken":
		return getGitHubIdentityToken(clientId, stackDomain)
	case "terraform_idtoken":
		return getTerraformIdentityToken()
	}
	return "", fmt.Errorf("no matching id token configuration")
}

func getGcpIdentityToken(clientId, stackDomain string) (string, error) {
	if isTokenValid(GCP_ID_TOKEN) {
		return GCP_ID_TOKEN, nil
	}

	audience := fmt.Sprintf("https://%s.id.%s", getAembitTenantId(clientId), stackDomain)
	metadataIdentityTokenUrl := fmt.Sprintf(
		"http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/identity?format=full&audience=%s",
		url.QueryEscape(audience),
	)

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

	audience := fmt.Sprintf("https://%s.id.%s", getAembitTenantId(clientId), stackDomain)
	identityTokenURL := fmt.Sprintf("%s&audience=%s", tokenRequestURL, url.QueryEscape(audience))

	req, err := http.NewRequest("GET", identityTokenURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenRequestToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch github id token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	jsonBody := make(map[string]interface{})
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		return "", fmt.Errorf("failed to parse response body: %w", err)
	}

	GITHUB_ID_TOKEN, ok := jsonBody["value"].(string)
	if !ok {
		return "", fmt.Errorf("failed to parse response value: %w", err)
	}

	return GITHUB_ID_TOKEN, nil
}

func getTerraformIdentityToken() (string, error) {
	if isTokenValid(TERRAFORM_ID_TOKEN) {
		return TERRAFORM_ID_TOKEN, nil
	}

	TERRAFORM_ID_TOKEN := os.Getenv("TFC_WORKLOAD_IDENTITY_TOKEN")
	return TERRAFORM_ID_TOKEN, nil
}

func getAembitTenantId(clientId string) string {
	clientIdSplit := strings.Split(clientId, ":")
	if len(clientIdSplit) >= 3 {
		return clientIdSplit[2]
	}

	return ""
}

func getAembitIdentityType(clientId string) string {
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
