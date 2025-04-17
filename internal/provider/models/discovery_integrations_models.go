package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// models.DiscoveryIntegrationResourceModel maps the resource schema.
type DiscoveryIntegrationResourceModel struct {
	// ID is required for Framework acceptance testing
	ID                   types.String                  `tfsdk:"id"`
	Name                 types.String                  `tfsdk:"name"`
	Description          types.String                  `tfsdk:"description"`
	IsActive             types.Bool                    `tfsdk:"is_active"`
	Tags                 types.Map                     `tfsdk:"tags"`
	Type                 types.String                  `tfsdk:"type"`
	SyncFrequencySeconds types.Int64                   `tfsdk:"sync_frequency_seconds"`
	LastSync             types.String                  `tfsdk:"last_sync"`
	LastSyncStatus       types.String                  `tfsdk:"last_sync_status"`
	Endpoint             types.String                  `tfsdk:"endpoint"`
	Wiz                  *DiscoveryIntegrationWizModel `tfsdk:"wiz_integration"`
}

type DiscoveryIntegrationWizModel struct {
	TokenUrl     types.String `tfsdk:"token_url"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	Audience     types.String `tfsdk:"audience"`
}

type DiscoveryIntegrationsDataSourceModel struct {
	DiscoveryIntegrations []DiscoveryIntegrationResourceModel `tfsdk:"discovery_integrations"`
}
