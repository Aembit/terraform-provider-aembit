package provider

import (
	"context"
	"fmt"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &serverWorkloadsDataSource{}
	_ datasource.DataSourceWithConfigure = &serverWorkloadsDataSource{}
)

// NewServerWorkloadsDataSource is a helper function to simplify the provider implementation.
func NewServerWorkloadsDataSource() datasource.DataSource {
	return &serverWorkloadsDataSource{}
}

// serverWorkloadsDataSource is the data source implementation.
type serverWorkloadsDataSource struct {
	client *aembit.Client
}

// Configure adds the provider configured client to the data source.
func (d *serverWorkloadsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*aembit.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *aembit.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Metadata returns the data source type name.
func (d *serverWorkloadsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_workloads"
}

// Schema defines the schema for the resource.
func (r *serverWorkloadsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an server workload.",
		Attributes: map[string]schema.Attribute{
			"server_workloads": schema.ListNestedAttribute{
				Description: "List of server workloads.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// ID field is required for Terraform Framework acceptance testing.
						"id": schema.StringAttribute{
							Description: "Alphanumeric identifier of the server workload.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "User-provided name of the server workload.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "User-provided description of the server workload.",
							Computed:    true,
						},
						"is_active": schema.BoolAttribute{
							Description: "Active/Inactive status of the server workload.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Type of server workload.",
							Computed:    true,
						},
						"tags": schema.ListNestedAttribute{
							Description: "List of Tags.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										Description: "Tag key.",
										Computed:    true,
									},
									"value": schema.StringAttribute{
										Description: "Tag value.",
										Computed:    true,
									},
								},
							},
						},
						"service_endpoint": schema.SingleNestedAttribute{
							Description: "Service endpoint details.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"external_id": schema.StringAttribute{
									Description: "Alphanumeric identifier of the service endpoint.",
									Computed:    true,
								},
								"id": schema.Int64Attribute{
									Description: "Number identifier of the service endpoint.",
									Computed:    true,
								},
								"host": schema.StringAttribute{
									Description: "hostname of the service endpoint.",
									Computed:    true,
								},
								"port": schema.Int64Attribute{
									Description: "hostname of the service endpoint.",
									Computed:    true,
								},
								"app_protocol": schema.StringAttribute{
									Description: "hostname of the service endpoint.",
									Computed:    true,
								},
								"requested_port": schema.Int64Attribute{
									Description: "hostname of the service endpoint.",
									Computed:    true,
								},
								"requested_tls": schema.BoolAttribute{
									Description: "hostname of the service endpoint.",
									Computed:    true,
								},
								"tls_verification": schema.StringAttribute{
									Description: "hostname of the service endpoint.",
									Computed:    true,
								},
								"transport_protocol": schema.StringAttribute{
									Description: "hostname of the service endpoint.",
									Computed:    true,
								},
								"tls": schema.BoolAttribute{
									Description: "hostname of the service endpoint.",
									Computed:    true,
								},
								"http_headers": schema.ListNestedAttribute{
									Description: "List of HTTP headers.",
									Computed:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"key": schema.StringAttribute{
												Description: "HTTP Header Key.",
												Computed:    true,
											},
											"value": schema.StringAttribute{
												Description: "HTTP Header value.",
												Computed:    true,
											},
										},
									},
								},
								"workload_service_authentication": schema.SingleNestedAttribute{
									Description: "Service authentication details.",
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										"method": schema.StringAttribute{
											Description: "Service authentication method.",
											Computed:    true,
										},
										"scheme": schema.StringAttribute{
											Description: "Service authentication method.",
											Computed:    true,
										},
										"config": schema.StringAttribute{
											Description: "Service authentication method.",
											Computed:    true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *serverWorkloadsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serverWorkloadsDataSourceModel

	server_workloads, err := d.client.GetServerWorkloads(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Aembit Server Workloads",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, server_workload := range server_workloads {
		serverWorkloadState := serverWorkloadResourceModel{
			ID:          types.StringValue(server_workload.EntityDTO.ExternalId),
			Name:        types.StringValue(server_workload.EntityDTO.Name),
			Description: types.StringValue(server_workload.EntityDTO.Description),
			IsActive:    types.BoolValue(server_workload.EntityDTO.IsActive),
			Type:        types.StringValue(server_workload.Type),
		}

		for _, tag := range server_workload.EntityDTO.Tags {
			tagState := tagModel{
				Key:   types.StringValue(tag.Key),
				Value: types.StringValue(tag.Value),
			}

			serverWorkloadState.Tags = append(serverWorkloadState.Tags, tagState)
		}

		serverWorkloadState.ServiceEndpoint = &serviceEndpointModel{
			ExternalId:        types.StringValue(server_workload.ServiceEndpoint.ExternalId),
			Id:                types.Int64Value(int64(server_workload.ServiceEndpoint.Id)),
			Host:              types.StringValue(server_workload.ServiceEndpoint.Host),
			AppProtocol:       types.StringValue(server_workload.ServiceEndpoint.AppProtocol),
			TransportProtocol: types.StringValue(server_workload.ServiceEndpoint.TransportProtocol),
			RequestedPort:     types.Int64Value(int64(server_workload.ServiceEndpoint.RequestedPort)),
			RequestedTls:      types.BoolValue(server_workload.ServiceEndpoint.RequestedTls),
			Port:              types.Int64Value(int64(server_workload.ServiceEndpoint.Port)),
			Tls:               types.BoolValue(server_workload.ServiceEndpoint.Tls),
			TlsVerification:   types.StringValue(server_workload.ServiceEndpoint.TlsVerification),
		}

		for _, header := range server_workload.ServiceEndpoint.HttpHeaders {
			headerState := keyValueModel{
				Key:   types.StringValue(header.Key),
				Value: types.StringValue(header.Value),
			}

			serverWorkloadState.ServiceEndpoint.HttpHeaders = append(serverWorkloadState.ServiceEndpoint.HttpHeaders, headerState)
		}

		if server_workload.ServiceEndpoint.WorkloadServiceAuthentication != nil {
			serverWorkloadState.ServiceEndpoint.WorkloadServiceAuthentication = &workloadServiceAuthenticationModel{
				Method: types.StringValue(server_workload.ServiceEndpoint.WorkloadServiceAuthentication.Method),
				Scheme: types.StringValue(server_workload.ServiceEndpoint.WorkloadServiceAuthentication.Scheme),
				Config: types.StringValue(server_workload.ServiceEndpoint.WorkloadServiceAuthentication.Config),
			}
		}

		state.ServerWorkloads = append(state.ServerWorkloads, serverWorkloadState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
