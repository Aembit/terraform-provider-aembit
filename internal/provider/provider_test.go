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
	"testing"
	"time"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-log/tflogtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testClient *aembit.CloudClient

func init() {
	tenant := os.Getenv("AEMBIT_TENANT_ID")
	stackDomain := os.Getenv("AEMBIT_STACK_DOMAIN")

	token := os.Getenv("AEMBIT_TOKEN")
	if token == "" {
		aembitClientID := os.Getenv("AEMBIT_CLIENT_ID")
		tenant = getAembitTenantId(aembitClientID)
		token, _ = getToken(context.Background(), aembitClientID, stackDomain, "", "test")
	}
	testClient, _ = aembit.NewClient(aembit.URLBuilder{}, &token, "", "test")
	testClient.Tenant = tenant
	testClient.StackDomain = stackDomain
}

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"aembit": providerserver.NewProtocol6WithError(New("test", "unittest")()),
}

func TestUnitResourceConfigure(t *testing.T) {
	t.Parallel()
	configResponse := resource.ConfigureResponse{}
	resourceConfigure(resource.ConfigureRequest{ProviderData: nil}, &configResponse)
	assert.Empty(t, configResponse.Diagnostics)

	resourceConfigure(resource.ConfigureRequest{ProviderData: "invalidData"}, &configResponse)
	assert.NotEmpty(t, configResponse.Diagnostics)
}

func TestUnitDataSourceConfigure(t *testing.T) {
	t.Parallel()
	configResponse := datasource.ConfigureResponse{}
	datasourceConfigure(datasource.ConfigureRequest{ProviderData: nil}, &configResponse)
	assert.Empty(t, configResponse.Diagnostics)

	datasourceConfigure(datasource.ConfigureRequest{ProviderData: "invalidData"}, &configResponse)
	assert.NotEmpty(t, configResponse.Diagnostics)
}

func TestUnitConfigureLogging(t *testing.T) {
	t.Parallel()
	// 1. Create a buffer to capture logs
	var buf bytes.Buffer

	// 2. Initialize a context with the test logger attached to the buffer
	// This context will intercept all calls to tflog.Debug, tflog.Info, etc.
	ctx := tflogtest.RootLogger(context.Background(), &buf)

	// 3. Define your provider/struct for testing
	p := New("1.2.3", "unittest")()

	// 4. Execute the code that calls tflog.Debug
	p.Configure(ctx, provider.ConfigureRequest{}, nil)

	// 5. Verify the output in the buffer
	loggedOutput := buf.String()

	assert.Contains(t, loggedOutput, "Aembit Provider version: 1.2.3")
	assert.Contains(t, loggedOutput, "Aembit Provider release time: unittest")
}

func TestUnitConfigureOldReleaseWarning(t *testing.T) {
	t.Parallel()
	// 1. Create a buffer to capture logs
	var buf bytes.Buffer

	// 2. Initialize a context with the test logger attached to the buffer
	ctx := tflogtest.RootLogger(context.Background(), &buf)

	// 3. Define your provider/struct for testing with an old release time
	// One year and one day ago
	oldDate := time.Now().AddDate(-1, 0, -1).Format(time.RFC3339)
	p := New("1.2.3", oldDate)()

	// 4. Execute the code that calls tflog.Warn and adds a diagnostic warning
	resp := &provider.ConfigureResponse{}
	ap, ok := p.(*aembitProvider)
	if !ok {
		t.Fatalf("expected *aembitProvider, got %T", p)
	}
	ap.checkVersionWarning(ctx, resp)

	// 5. Verify the output in the buffer
	loggedOutput := buf.String()
	assert.Contains(t, loggedOutput, "This Aembit Provider version (1.2.3) is more than 1 year old")

	// 6. Verify the diagnostic warning
	assert.NotEmpty(t, resp.Diagnostics)
	found := false
	for _, diag := range resp.Diagnostics {
		if diag.Summary() == "Aembit Provider Version Warning" &&
			assert.Contains(t, diag.Detail(), "This Aembit Provider version (1.2.3) is more than 1 year old") {
			found = true
			break
		}
	}
	assert.True(t, found, "Expected diagnostic warning not found")
}

type mockRoundTripper struct {
	fn func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.fn(req)
}

func createMockJWT(exp time.Time) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"exp":%d}`, exp.Unix())))
	sig := base64.RawURLEncoding.EncodeToString([]byte("signature"))
	return fmt.Sprintf("%s.%s.%s", header, payload, sig)
}

func TestUnitOidcIdentityToken(t *testing.T) {
	origOidc := OIDC_ID_TOKEN
	origEnv := os.Getenv("AEMBIT_OIDC_ID_TOKEN")
	defer func() {
		OIDC_ID_TOKEN = origOidc
		if origEnv == "" {
			_ = os.Unsetenv("AEMBIT_OIDC_ID_TOKEN")
		} else {
			_ = os.Setenv("AEMBIT_OIDC_ID_TOKEN", origEnv)
		}
	}()

	t.Run("cached valid token", func(t *testing.T) {
		validToken := createMockJWT(time.Now().Add(10 * time.Minute))
		OIDC_ID_TOKEN = validToken
		_ = os.Setenv("AEMBIT_OIDC_ID_TOKEN", "different-env-token")

		token, err := getOidcIdentityToken()
		require.NoError(t, err)
		assert.Equal(t, validToken, token)
	})

	t.Run("fallback to env var when cache invalid or empty", func(t *testing.T) {
		OIDC_ID_TOKEN = ""
		envToken := "test-oidc-env-token"
		_ = os.Setenv("AEMBIT_OIDC_ID_TOKEN", envToken)

		token, err := getOidcIdentityToken()
		require.NoError(t, err)
		assert.Equal(t, envToken, token)
		assert.Equal(t, envToken, OIDC_ID_TOKEN)
	})

	t.Run("returns empty string when no token configured", func(t *testing.T) {
		OIDC_ID_TOKEN = ""
		_ = os.Unsetenv("AEMBIT_OIDC_ID_TOKEN")

		token, err := getOidcIdentityToken()
		require.NoError(t, err)
		assert.Empty(t, token)
	})
}

func TestUnitGetIdentityTokenOidc(t *testing.T) {
	origOidc := OIDC_ID_TOKEN
	origEnv := os.Getenv("AEMBIT_OIDC_ID_TOKEN")
	defer func() {
		OIDC_ID_TOKEN = origOidc
		if origEnv == "" {
			_ = os.Unsetenv("AEMBIT_OIDC_ID_TOKEN")
		} else {
			_ = os.Setenv("AEMBIT_OIDC_ID_TOKEN", origEnv)
		}
	}()

	t.Run("valid oidc_id_token client id", func(t *testing.T) {
		OIDC_ID_TOKEN = ""
		_ = os.Setenv("AEMBIT_OIDC_ID_TOKEN", "my-oidc-token")

		clientID := "aembit:dev:tenant123:identity:oidc_id_token:ext456"
		token, err := getIdentityToken(clientID, "aembit.io")
		require.NoError(t, err)
		assert.Equal(t, "my-oidc-token", token)
	})

	t.Run("invalid client id identity type", func(t *testing.T) {
		clientID := "aembit:dev:tenant123:identity:unsupported_type:ext456"
		token, err := getIdentityToken(clientID, "aembit.io")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no matching id token configuration")
		assert.Empty(t, token)
	})
}

func TestUnitGetWorkloadAssessmentOidc(t *testing.T) {
	t.Run("oidc_id_token without resource set", func(t *testing.T) {
		clientID := "aembit:dev:tenant123:identity:oidc_id_token:ext456"
		idToken := "test-oidc-id-token"

		assessmentStr, err := getWorkloadAssessment(clientID, idToken, "")
		require.NoError(t, err)

		var workload WorkloadAssessment
		err = json.Unmarshal([]byte(assessmentStr), &workload)
		require.NoError(t, err)

		assert.Equal(t, "1.0.0", workload.Version)
		assert.Equal(t, idToken, workload.Oidc.IdentityToken)
		assert.Empty(t, workload.GCP.IdentityToken)
		assert.Empty(t, workload.GitHub.IdentityToken)
		assert.Empty(t, workload.Terraform.IdentityToken)
		assert.Empty(t, workload.OS.Environment.ResourceSet)
	})

	t.Run("oidc_id_token with resource set", func(t *testing.T) {
		clientID := "aembit:dev:tenant123:identity:oidc_id_token:ext456"
		idToken := "test-oidc-id-token"
		resourceSetID := "rs-789"

		assessmentStr, err := getWorkloadAssessment(clientID, idToken, resourceSetID)
		require.NoError(t, err)

		var workload WorkloadAssessment
		err = json.Unmarshal([]byte(assessmentStr), &workload)
		require.NoError(t, err)

		assert.Equal(t, "1.0.0", workload.Version)
		assert.Equal(t, idToken, workload.Oidc.IdentityToken)
		assert.Equal(t, resourceSetID, workload.OS.Environment.ResourceSet)
	})

	t.Run("invalid client id", func(t *testing.T) {
		clientID := "invalid-client-id"
		_, err := getWorkloadAssessment(clientID, "some-token", "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid aembit client id")
	})
}

func TestUnitGetAembitTokenOidc(t *testing.T) {
	origAembit := AEMBIT_TOKEN
	origTransport := http.DefaultTransport
	defer func() {
		AEMBIT_TOKEN = origAembit
		http.DefaultTransport = origTransport
	}()

	clientID := "aembit:dev:tenant123:identity:oidc_id_token:ext456"
	stackDomain := "aembit.io"

	t.Run("return cached AEMBIT_TOKEN if valid", func(t *testing.T) {
		validToken := createMockJWT(time.Now().Add(10 * time.Minute))
		AEMBIT_TOKEN = validToken

		token, err := getAembitToken(clientID, stackDomain, "id-token", "", "1.0.0")
		require.NoError(t, err)
		assert.Equal(t, validToken, token)
	})

	t.Run("fetch token via HTTP for oidc_id_token", func(t *testing.T) {
		AEMBIT_TOKEN = ""
		idToken := "test-oidc-token-value"
		resourceSetID := "rs-test-123"

		expectedAccessToken := createMockJWT(time.Now().Add(10 * time.Minute))

		http.DefaultTransport = &mockRoundTripper{
			fn: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, "POST", req.Method)
				assert.Equal(t, "https://tenant123.id.aembit.io/connect/token", req.URL.String())
				assert.Equal(t, "application/x-www-form-urlencoded;charset=UTF-8", req.Header.Get("Content-Type"))
				assert.Equal(t, "AembitTerraformProvider/1.0.0", req.Header.Get("User-Agent"))
				assert.Equal(t, resourceSetID, req.Header.Get("X-Aembit-ResourceSet"))

				bodyBytes, err := io.ReadAll(req.Body)
				require.NoError(t, err)

				values, err := url.ParseQuery(string(bodyBytes))
				require.NoError(t, err)

				assert.Equal(t, "client_credentials", values.Get("grant_type"))
				assert.Equal(t, clientID, values.Get("client_id"))

				attestationJSON := values.Get("attestation")
				var attestation map[string]interface{}
				err = json.Unmarshal([]byte(attestationJSON), &attestation)
				require.NoError(t, err)

				assert.Equal(t, "1.0.0", attestation["version"])
				oidcMap, ok := attestation["oidc"].(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, idToken, oidcMap["identityToken"])

				respBody := fmt.Sprintf(`{"access_token":"%s"}`, expectedAccessToken)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(respBody)),
					Header:     make(http.Header),
				}, nil
			},
		}

		token, err := getAembitToken(clientID, stackDomain, idToken, resourceSetID, "1.0.0")
		require.NoError(t, err)
		assert.Equal(t, expectedAccessToken, token)
		assert.Equal(t, expectedAccessToken, AEMBIT_TOKEN)
	})

	t.Run("invalid client id identity type", func(t *testing.T) {
		AEMBIT_TOKEN = ""
		invalidClientID := "aembit:dev:tenant123:identity:unknown:ext456"

		_, err := getAembitToken(invalidClientID, stackDomain, "id-token", "", "1.0.0")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid aembit client id")
	})

	t.Run("handles HTTP response unmarshal error", func(t *testing.T) {
		AEMBIT_TOKEN = ""

		http.DefaultTransport = &mockRoundTripper{
			fn: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("invalid-json")),
					Header:     make(http.Header),
				}, nil
			},
		}

		_, err := getAembitToken(clientID, stackDomain, "id-token", "", "1.0.0")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal response")
	})
}

func TestUnitGetTokenOidcError(t *testing.T) {
	origOidc := OIDC_ID_TOKEN
	origAembit := AEMBIT_TOKEN
	origEnv := os.Getenv("AEMBIT_OIDC_ID_TOKEN")
	origTransport := http.DefaultTransport
	defer func() {
		OIDC_ID_TOKEN = origOidc
		AEMBIT_TOKEN = origAembit
		http.DefaultTransport = origTransport
		if origEnv == "" {
			_ = os.Unsetenv("AEMBIT_OIDC_ID_TOKEN")
		} else {
			_ = os.Setenv("AEMBIT_OIDC_ID_TOKEN", origEnv)
		}
	}()

	OIDC_ID_TOKEN = ""
	AEMBIT_TOKEN = ""
	_ = os.Setenv("AEMBIT_OIDC_ID_TOKEN", "valid-id-token")

	clientID := "aembit:dev:tenant123:identity:oidc_id_token:ext456"

	// Mock http transport to return an error
	http.DefaultTransport = &mockRoundTripper{
		fn: func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("network connection error")
		},
	}

	_, err := getToken(context.Background(), clientID, "aembit.io", "", "1.0.0")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "network connection error")
}
