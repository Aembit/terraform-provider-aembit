// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"testing"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/stretchr/testify/assert"
)

var (
	Client *aembit.CloudClient
)

func init() {
	tenant := os.Getenv("AEMBIT_TENANT_ID")
	stackDomain := os.Getenv("AEMBIT_STACK_DOMAIN")

	token := os.Getenv("AEMBIT_TOKEN")
	if token == "" {
		token, _ = getToken(context.Background(), os.Getenv("AEMBIT_CLIENT_ID"), stackDomain)
	}
	Client, _ = aembit.NewClient(aembit.URLBuilder{}, &token, "test")
	Client.Tenant = tenant
	Client.StackDomain = stackDomain
}

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"aembit": providerserver.NewProtocol6WithError(New("test")()),
}

func TestUnitResourceConfigure(t *testing.T) {
	var configResponse resource.ConfigureResponse = resource.ConfigureResponse{}
	resourceConfigure(resource.ConfigureRequest{ProviderData: nil}, &configResponse)
	assert.Empty(t, configResponse.Diagnostics)

	resourceConfigure(resource.ConfigureRequest{ProviderData: "invalidData"}, &configResponse)
	assert.NotEmpty(t, configResponse.Diagnostics)
}

func TestUnitDataSourceConfigure(t *testing.T) {
	var configResponse datasource.ConfigureResponse = datasource.ConfigureResponse{}
	datasourceConfigure(datasource.ConfigureRequest{ProviderData: nil}, &configResponse)
	assert.Empty(t, configResponse.Diagnostics)

	datasourceConfigure(datasource.ConfigureRequest{ProviderData: "invalidData"}, &configResponse)
	assert.NotEmpty(t, configResponse.Diagnostics)
}
