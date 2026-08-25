package provider

import (
	"context"

	"terraform-provider-aembit/internal/provider/models"
	"terraform-provider-aembit/internal/provider/validators"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &mcpToolAccessRuleResource{}
	_ resource.ResourceWithConfigure   = &mcpToolAccessRuleResource{}
	_ resource.ResourceWithImportState = &mcpToolAccessRuleResource{}
)

// NewMcpToolAccessRuleResource is a helper function to simplify the provider implementation.
func NewMcpToolAccessRuleResource() resource.Resource {
	return &mcpToolAccessRuleResource{}
}

// mcpToolAccessRuleResource is the resource implementation.
type mcpToolAccessRuleResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *mcpToolAccessRuleResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_mcp_tool_access_rule"
}

// Configure adds the provider configured client to the resource.
func (r *mcpToolAccessRuleResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = resourceConfigure(req, resp)
}

// Schema defines the schema for the resource.
func (r *mcpToolAccessRuleResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages an MCP Tool Access Rule resource associated with an MCP Tool Access Control Content Security.",
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the MCP Tool Access Rule.",
				Computed:    true,
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
			},
			"content_security_id": schema.StringAttribute{
				Description: "Unique identifier of the Content Security.",
				Required:    true,
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
			},
			"tool_name": schema.StringAttribute{
				Description: "The name or pattern matching tool name.",
				Required:    true,
			},
			"is_visible": schema.BoolAttribute{
				Description: "Whether the tool is visible to the agent.",
				Required:    true,
			},
			"is_invocable": schema.BoolAttribute{
				Description: "Whether the tool is invocable by the agent.",
				Required:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *mcpToolAccessRuleResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Retrieve values from plan
	var plan models.McpToolAccessRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch Content Security to get its Resource Set
	cs, err, notFound := r.client.GetContentSecurity(plan.ContentSecurityID.ValueString(), nil, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching Content Security",
			"Could not read Aembit Content Security resource "+plan.ContentSecurityID.ValueString()+": "+err.Error(),
		)
		return
	}
	if notFound {
		resp.Diagnostics.AddError(
			"Content Security not found",
			"The specified Content Security "+plan.ContentSecurityID.ValueString()+" does not exist.",
		)
		return
	}

	// Generate API request body from plan
	dto := convertMcpToolAccessRuleResourceModelToDTO(plan, nil)

	// Create new MCP Tool Access Rule
	newRule, err := r.client.UpsertMcpToolRule(
		plan.ContentSecurityID.ValueString(),
		dto,
		nil,
		&cs.ResourceSet,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating MCP Tool Access Rule",
			"Could not create MCP Tool Access Rule, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertMcpToolAccessRuleDTOToModel(*newRule, &plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *mcpToolAccessRuleResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Get current state
	var state models.McpToolAccessRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch Content Security to get its Resource Set
	cs, err, notFound := r.client.GetContentSecurity(state.ContentSecurityID.ValueString(), nil, nil)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading Aembit Content Security",
			"Could not read Aembit Content Security resource "+state.ContentSecurityID.ValueString()+": "+err.Error(),
		)
		if notFound {
			resp.State.RemoveResource(ctx)
		}
		return
	}

	// Get refreshed rule value from Aembit
	rule, err, notFound := r.client.GetMcpToolRule(
		state.ContentSecurityID.ValueString(),
		state.ID.ValueString(),
		nil,
		&cs.ResourceSet,
	)

	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading Aembit MCP Tool Access Rule",
			"Could not read MCP Tool Access Rule "+state.ID.ValueString()+" from Aembit: "+err.Error(),
		)

		// If the resource is not found on Aembit Cloud, delete it locally
		if notFound {
			resp.State.RemoveResource(ctx)
		}
		return
	}

	state = convertMcpToolAccessRuleDTOToModel(rule, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *mcpToolAccessRuleResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Get current state
	var state models.McpToolAccessRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from plan
	var plan models.McpToolAccessRuleResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch Content Security to get its Resource Set
	cs, err, notFound := r.client.GetContentSecurity(plan.ContentSecurityID.ValueString(), nil, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching Content Security",
			"Could not read Aembit Content Security resource "+plan.ContentSecurityID.ValueString()+": "+err.Error(),
		)
		return
	}
	if notFound {
		resp.Diagnostics.AddError(
			"Content Security not found",
			"The specified Content Security "+plan.ContentSecurityID.ValueString()+" does not exist.",
		)
		return
	}

	ruleID := state.ID.ValueString()
	dto := convertMcpToolAccessRuleResourceModelToDTO(plan, &ruleID)

	// Update MCP Tool Access Rule
	updatedRule, err := r.client.UpsertMcpToolRule(
		plan.ContentSecurityID.ValueString(),
		dto,
		nil,
		&cs.ResourceSet,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating MCP Tool Access Rule",
			"Could not update MCP Tool Access Rule "+ruleID+", unexpected error: "+err.Error(),
		)
		return
	}

	state = convertMcpToolAccessRuleDTOToModel(*updatedRule, &plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *mcpToolAccessRuleResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve values from state
	var state models.McpToolAccessRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch Content Security to get its Resource Set
	cs, err, notFound := r.client.GetContentSecurity(state.ContentSecurityID.ValueString(), nil, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching Content Security",
			"Could not read Aembit Content Security resource "+state.ContentSecurityID.ValueString()+": "+err.Error(),
		)
		return
	}
	if notFound {
		// Parent resource is gone, rule is deleted by cascade
		return
	}

	// Delete existing rule
	_, err = r.client.DeleteMcpToolRule(
		ctx,
		state.ContentSecurityID.ValueString(),
		state.ID.ValueString(),
		nil,
		&cs.ResourceSet,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting MCP Tool Access Rule",
			"Could not delete MCP Tool Access Rule, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing content_security_id,rule_id.
func (r *mcpToolAccessRuleResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Custom import logic or passthrough
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertMcpToolAccessRuleResourceModelToDTO(
	model models.McpToolAccessRuleResourceModel,
	externalID *string,
) aembit.McpToolAccessControlRuleDTO {
	var dto aembit.McpToolAccessControlRuleDTO
	dto.ToolName = model.ToolName.ValueString()
	dto.IsVisible = model.IsVisible.ValueBool()
	dto.IsInvocable = model.IsInvocable.ValueBool()

	if externalID != nil {
		dto.ExternalID = *externalID
	}

	return dto
}

func convertMcpToolAccessRuleDTOToModel(
	dto aembit.McpToolAccessControlRuleDTO,
	planModel *models.McpToolAccessRuleResourceModel,
) models.McpToolAccessRuleResourceModel {
	var model models.McpToolAccessRuleResourceModel
	model.ID = types.StringValue(dto.ExternalID)
	model.ContentSecurityID = planModel.ContentSecurityID
	model.ToolName = types.StringValue(dto.ToolName)
	model.IsVisible = types.BoolValue(dto.IsVisible)
	model.IsInvocable = types.BoolValue(dto.IsInvocable)

	return model
}
