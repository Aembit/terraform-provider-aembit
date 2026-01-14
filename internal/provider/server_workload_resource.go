package provider

import (
	"context"
	"fmt"
	"terraform-provider-aembit/internal/provider/models"
	"terraform-provider-aembit/internal/provider/validators"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &serverWorkloadResource{}
	_ resource.ResourceWithConfigure   = &serverWorkloadResource{}
	_ resource.ResourceWithImportState = &serverWorkloadResource{}
	_ resource.ResourceWithModifyPlan  = &serverWorkloadResource{}
)

const mcpAuthServerTemplate string = "https://%s.mcp.%s/auth"

// NewServerWorkloadResource is a helper function to simplify the provider implementation.
func NewServerWorkloadResource() resource.Resource {
	return &serverWorkloadResource{}
}

// serverWorkloadResource is the resource implementation.
type serverWorkloadResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *serverWorkloadResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_server_workload"
}

// Configure adds the provider configured client to the resource.
func (r *serverWorkloadResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = resourceConfigure(req, resp)
}

// Schema defines the schema for the resource.
func (r *serverWorkloadResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the Server Workload.",
				Computed:    true,
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name for the Server Workload.",
				Required:    true,
				Validators: []validator.String{
					validators.NameLengthValidation(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description for the Server Workload.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active status of the Server Workload.",
				Optional:    true,
				Computed:    true,
			},
			"tags":     TagsMapAttribute(),
			"tags_all": TagsAllMapAttribute(),
			"service_endpoint": schema.SingleNestedAttribute{
				Description: "Service endpoint details.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"external_id": schema.StringAttribute{
						Description: "Unique identifier of the service endpoint.",
						Computed:    true,
					},
					"id": schema.Int64Attribute{
						Description: "Number identifier of the service endpoint.",
						Computed:    true,
					},
					"host": schema.StringAttribute{
						Description: "Hostname of the Server Workload service endpoint.\n" +
							"Wildcard hostnames are supported, for example `*.amazonaws.com`, `*.azure.com`, or `*.googleapis.com`.\n" +
							"Note: Wildcards are not supported in the top or second-level domain, such as `*.com`.",
						Required: true,
						Validators: []validator.String{
							validators.SafeWildcardHostNameValidation(),
						},
					},
					"aembit_mcp_authorization_server_url": schema.StringAttribute{
						Description: "Aembit MCP Authorization Server URL",
						Computed:    true,
					},
					"port": schema.Int64Attribute{
						Description: "Port of the Server Workload service endpoint.",
						Required:    true,
						Validators: []validator.Int64{
							int64validator.Between(1, 65535),
						},
					},
					"app_protocol": schema.StringAttribute{
						Description: "Application Protocol of the Server Workload service endpoint. Possible values are: \n" +
							"\t* `Amazon Redshift`\n" +
							"\t* `HTTP`\n" +
							"\t* `MCP`\n" +
							"\t* `MySQL`\n" +
							"\t* `PostgreSQL`\n" +
							"\t* `Redis`\n" +
							"\t* `Snowflake`\n",
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{
								"Amazon Redshift",
								"HTTP",
								"MCP",
								"MySQL",
								"OAuth",
								"PostgreSQL",
								"Redis",
								"Snowflake",
							}...),
						},
					},
					"requested_port": schema.Int64Attribute{
						Description: "Requested port of the Server Workload service endpoint.",
						Optional:    true,
						Computed:    true,
						Validators: []validator.Int64{
							int64validator.Between(1, 65535),
						},
					},
					"tls_verification": schema.StringAttribute{
						Description: "TLS verification configuration of the Server Workload service endpoint. Possible values are `full` (default) or `none`.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{
								"full",
								"none",
							}...),
						},
					},
					"transport_protocol": schema.StringAttribute{
						Description: "Transport protocol of the Server Workload service endpoint. This value must be set to the default `TCP`.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{
								"TCP",
							}...),
						},
					},
					"requested_tls": schema.BoolAttribute{
						Description: "TLS requested on the Server Workload service endpoint.",
						Optional:    true,
						Computed:    true,
					},
					"tls": schema.BoolAttribute{
						Description: "TLS indicated on the Server Workload service endpoint.",
						Optional:    true,
						Computed:    true,
					},
					"url_path": schema.StringAttribute{
						Description: "URL path of the Server Workload service endpoint. <br>This value is only applicable when the `app_protocol` is set to `OAuth`.",
						Required:    false,
						Optional:    true,
						Computed:    true,
					},
					"http_headers": schema.MapAttribute{
						Description: "HTTP Headers are key-value pairs.",
						ElementType: types.StringType,
						Optional:    true,
					},
					"authentication_config": schema.SingleNestedAttribute{
						Description: "Service authentication details.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"method": schema.StringAttribute{
								Description: "Server Workload Service authentication method. Possible values are: \n" +
									"\t* `API Key`\n" +
									"\t* `HTTP Authentication`\n" +
									"\t* `JWT Token Authentication`\n" +
									"\t* `OAuth Client Authentication`\n" +
									"\t* `Password Authentication`\n",
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf([]string{
										"API Key",
										"HTTP Authentication",
										"JWT Token Authentication",
										"OAuth Client Authentication",
										"Password Authentication",
									}...),
								},
							},
							"scheme": schema.StringAttribute{
								Description: "Server Workload Service authentication scheme. Possible values are: \n" +
									"\t* For Authentation Method `API Key`:\n" +
									"\t\t* `Header`\n" +
									"\t\t* `Query Parameter`\n" +
									"\t* For Authentation Method `HTTP Authentication`:\n" +
									"\t\t* `Basic`\n" +
									"\t\t* `Bearer`\n" +
									"\t\t* `Header`\n" +
									"\t\t* `AWS Signature v4`\n" +
									"\t* For Authentation Method `JWT Token Authentication`:\n" +
									"\t\t* `Snowflake JWT`\n" +
									"\t* For Authentation Method `OAuth Client Authentication`:\n" +
									"\t\t* `POST Body`\n" +
									"\t* For Authentation Method `Password Authentication`:\n" +
									"\t\t* `Password`\n",
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf([]string{
										"Header",
										"Query Parameter",
										"Basic",
										"Bearer",
										"Header",
										"AWS Signature v4",
										"Snowflake JWT",
										"Password",
										"POST Body",
									}...),
								},
							},
							"config": schema.StringAttribute{
								Description: "Server Workload Service authentication config. <br>This value is used to identify the HTTP Header or Query Parameter used for the associated authentication scheme. <br>**Note:** This value is required in cases where an HTTP Header or Query Parameter is required, for example with `HTTP Authentication` and scheme `Header`.",
								Optional:    true,
								Computed:    true,
							},
						},
					},
				},
			},
		},
	}
}

func (r *serverWorkloadResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	modifyPlanForTagsAll(ctx, req, resp, r.client.DefaultTags)
}

// Create creates the resource and sets the initial Terraform state.
func (r *serverWorkloadResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Retrieve values from plan
	var plan models.ServerWorkloadResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validation: requested_port must equal port when app_protocol is MCP
	if plan.ServiceEndpoint.AppProtocol.ValueString() == "MCP" &&
		plan.ServiceEndpoint.RequestedPort.ValueInt64() != plan.ServiceEndpoint.Port.ValueInt64() {
		resp.Diagnostics.AddError(
			"Invalid requested_port for MCP protocol",
			"When app_protocol is 'MCP', requested_port must be equal to port.",
		)
		return
	}

	// Generate API request body from plan
	workload := convertServerWorkloadModelToDTO(ctx, plan, nil, r.client.DefaultTags)

	// Create new Server Workload
	serverWorkload, err := r.client.CreateServerWorkload(workload, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating server workload",
			"Could not create server workload, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertServerWorkloadDTOToModel(ctx, *serverWorkload, &plan,
		r.client.Tenant,
		r.client.StackDomain)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *serverWorkloadResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Get current state
	var state models.ServerWorkloadResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed workload value from Aembit
	serverWorkload, err, notFound := r.client.GetServerWorkload(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading Aembit Server Workload",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)

		// If the resource is not found on Aembit Cloud, delete it locally
		if notFound {
			resp.State.RemoveResource(ctx)
		}
		return
	}

	// Overwrite items with refreshed state
	state = convertServerWorkloadDTOToModel(ctx, serverWorkload, &state,
		r.client.Tenant,
		r.client.StackDomain)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *serverWorkloadResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Get current state
	var state models.ServerWorkloadResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan models.ServerWorkloadResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validation: requested_port must equal port when app_protocol is MCP
	if plan.ServiceEndpoint.AppProtocol.ValueString() == "MCP" &&
		plan.ServiceEndpoint.RequestedPort.ValueInt64() != plan.ServiceEndpoint.Port.ValueInt64() {
		resp.Diagnostics.AddError(
			"Invalid requested_port for MCP protocol",
			"When app_protocol is 'MCP', requested_port must be equal to port.",
		)
		return	
	}

	// Generate API request body from plan
	workload := convertServerWorkloadModelToDTO(ctx, plan, &externalID, r.client.DefaultTags)

	// Update Server Workload
	serverWorkload, err := r.client.UpdateServerWorkload(workload, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating server workload",
			"Could not update server workload, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertServerWorkloadDTOToModel(ctx, *serverWorkload, &plan,
		r.client.Tenant,
		r.client.StackDomain)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *serverWorkloadResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve values from state
	var state models.ServerWorkloadResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Server Workload is Active - if it is, disable it first
	if state.IsActive == types.BoolValue(true) {
		_, err := r.client.DisableServerWorkload(state.ID.ValueString(), nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error disabling Server Workload",
				"Could not disable Server Workload, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Delete existing Server Workload
	_, err := r.client.DeleteServerWorkload(ctx, state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Server Workload",
			"Could not delete server workload, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalID.
func (r *serverWorkloadResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Retrieve import externalID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertServerWorkloadModelToDTO(
	ctx context.Context,
	model models.ServerWorkloadResourceModel,
	externalID *string,
	defaultTags map[string]string,
) aembit.ServerWorkloadExternalDTO {
	var workload aembit.ServerWorkloadExternalDTO
	workload.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}

	workload.ServiceEndpoint = aembit.WorkloadServiceEndpointDTO{
		Host:              model.ServiceEndpoint.Host.ValueString(),
		ID:                int(model.ServiceEndpoint.ID.ValueInt64()),
		Port:              int(model.ServiceEndpoint.Port.ValueInt64()),
		AppProtocol:       model.ServiceEndpoint.AppProtocol.ValueString(),
		TransportProtocol: model.ServiceEndpoint.TransportProtocol.ValueString(),
		RequestedPort: func() int {
			if model.ServiceEndpoint.AppProtocol.ValueString() == "MCP" {
				return int(model.ServiceEndpoint.Port.ValueInt64())
			}
			return int(model.ServiceEndpoint.RequestedPort.ValueInt64())
		}(),
		RequestedTLS: func() bool {
			if model.ServiceEndpoint.AppProtocol.ValueString() == "MCP" {
				return model.ServiceEndpoint.TLS.ValueBool()
			}
			return model.ServiceEndpoint.RequestedTLS.ValueBool()
		}(),
		TLS:             model.ServiceEndpoint.TLS.ValueBool(),
		TLSVerification: model.ServiceEndpoint.TLSVerification.ValueString(),
		URLPath:         model.ServiceEndpoint.URLPath.ValueString(),
	}

	if model.ServiceEndpoint.WorkloadServiceAuthentication != nil {
		workload.ServiceEndpoint.WorkloadServiceAuthentication = &aembit.WorkloadServiceAuthenticationDTO{
			Method: model.ServiceEndpoint.WorkloadServiceAuthentication.Method.ValueString(),
			Scheme: model.ServiceEndpoint.WorkloadServiceAuthentication.Scheme.ValueString(),
			Config: model.ServiceEndpoint.WorkloadServiceAuthentication.Config.ValueString(),
		}
	}

	if len(model.ServiceEndpoint.HTTPHeaders.Elements()) > 0 {
		headersMap := make(map[string]string)
		_ = model.ServiceEndpoint.HTTPHeaders.ElementsAs(ctx, &headersMap, true)

		for key, value := range headersMap {
			workload.ServiceEndpoint.HTTPHeaders = append(
				workload.ServiceEndpoint.HTTPHeaders,
				aembit.KeyValuePair{
					Key:   key,
					Value: value,
				},
			)
		}
	}

	if externalID != nil {
		workload.ExternalID = *externalID
	}

	workload.Tags = collectAllTagsDto(ctx, defaultTags, model.Tags)
	return workload
}

func convertServerWorkloadDTOToModel(
	ctx context.Context,
	dto aembit.ServerWorkloadExternalDTO,
	planModel *models.ServerWorkloadResourceModel,
	tenant string,
	stackDomain string,
) models.ServerWorkloadResourceModel {
	var model models.ServerWorkloadResourceModel
	model.ID = types.StringValue(dto.ExternalID)
	model.Name = types.StringValue(dto.Name)
	model.Description = types.StringValue(dto.Description)
	model.IsActive = types.BoolValue(dto.IsActive)

	// handle tags
	model.Tags = newTagsModelFromPlan(ctx, planModel.Tags)
	model.TagsAll = newTagsModel(ctx, dto.Tags)

	model.ServiceEndpoint = &models.ServiceEndpointModel{
		ExternalID:        types.StringValue(dto.ServiceEndpoint.ExternalID),
		Host:              types.StringValue(dto.ServiceEndpoint.Host),
		Port:              types.Int64Value(int64(dto.ServiceEndpoint.Port)),
		AppProtocol:       types.StringValue(dto.ServiceEndpoint.AppProtocol),
		TransportProtocol: types.StringValue(dto.ServiceEndpoint.TransportProtocol),
		RequestedPort:     types.Int64Value(int64(dto.ServiceEndpoint.RequestedPort)),
		RequestedTLS:      types.BoolValue(dto.ServiceEndpoint.RequestedTLS),
		TLS:               types.BoolValue(dto.ServiceEndpoint.TLS),
		TLSVerification:   types.StringValue(dto.ServiceEndpoint.TLSVerification),
		URLPath:           types.StringValue(dto.ServiceEndpoint.URLPath),
		MCPAuthorizationServerURL: types.StringValue(
			fmt.Sprintf(mcpAuthServerTemplate, tenant, stackDomain),
		),
	}
	model.ServiceEndpoint.HTTPHeaders = newHTTPHeadersModel(ctx, dto.ServiceEndpoint.HTTPHeaders)

	if dto.ServiceEndpoint.WorkloadServiceAuthentication != nil {
		model.ServiceEndpoint.WorkloadServiceAuthentication = &models.WorkloadServiceAuthenticationModel{
			Scheme: types.StringValue(dto.ServiceEndpoint.WorkloadServiceAuthentication.Scheme),
			Method: types.StringValue(dto.ServiceEndpoint.WorkloadServiceAuthentication.Method),
			Config: types.StringValue(dto.ServiceEndpoint.WorkloadServiceAuthentication.Config),
		}
	}

	return model
}
