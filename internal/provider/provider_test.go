// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"context"
	"os"
	"testing"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-log/tflogtest"
	"github.com/stretchr/testify/assert"
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
