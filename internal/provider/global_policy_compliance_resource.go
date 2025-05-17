package provider

import (
	"context"
	"fmt"
	"terraform-provider-aembit/internal/provider/models"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	AccessPolicyTrustProviderComplianceSettingName    = "ap_trust_prov_governance"
	AccessPolicyAccessConditionComplianceSettingName  = "ap_access_cond_governance"
	AgentControllerTrustProviderComplianceSettingName = "edge_trust_prov_governance"
	AgentControllerTlsHostNameComplianceSettingName   = "edge_tls_hostname_governance"
)

var (
	_ resource.Resource                = &GlobalPolicyComplianceResource{}
	_ resource.ResourceWithConfigure   = &GlobalPolicyComplianceResource{}
	_ resource.ResourceWithImportState = &GlobalPolicyComplianceResource{}
)

type GlobalPolicyComplianceResource struct {
	client *aembit.CloudClient
}

func NewGlobalPolicyComplianceResource() resource.Resource {
	return &GlobalPolicyComplianceResource{}
}

// ImportState implements resource.ResourceWithImportState.
func (gpcResource *GlobalPolicyComplianceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	gpcSettings, err := gpcResource.client.GetGlobalPolicyComplianceSettings(nil)
	if err != nil {
		resp.Diagnostics.AddError("Unable to retrieve Global Policy Compliance settings", err.Error())
		return
	}
	state := convertGlobalPolicyComplianceDTOToModel(gpcSettings, gpcResource.client.Tenant)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Configure implements resource.ResourceWithConfigure.
func (g *GlobalPolicyComplianceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	g.client = resourceConfigure(req, resp)
}

// Create implements resource.Resource.
func (g *GlobalPolicyComplianceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var gpcModel models.GlobalPolicyComplianceModel
	diags := req.Plan.Get(ctx, &gpcModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := updateComplianceSettings(g.client, &gpcModel, nil)

	if err != nil {
		resp.Diagnostics.AddError("Error updating Global Policy Compliance settings during resource creation", "Error: "+err.Error())
		return
	}
	gpcModel.Id = types.StringValue(g.client.Tenant + "-gpc")
	diags = resp.State.Set(ctx, gpcModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete implements resource.Resource.
func (g *GlobalPolicyComplianceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning("Resetting Global Policy Compliance settings default values",
		"Deleting Global Policy Compliance settings will result in resetting their values to the default 'Recommended' value.")

	var state models.GlobalPolicyComplianceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var defaultModel = models.GlobalPolicyComplianceModel{
		Id:                             types.StringValue(g.client.Tenant + "-gpc"),
		APTrustProviderCompliance:      types.StringValue("Recommended"),
		APAccessConditionCompliance:    types.StringValue("Recommended"),
		ACTrustProviderCompliance:      types.StringValue("Recommended"),
		ACAllowedTLSHostanmeCompliance: types.StringValue("Recommended"),
	}
	err := updateComplianceSettings(g.client, &defaultModel, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Global Policy Compliance settings to their default values", "Error: "+err.Error())
		return
	}
}

// Metadata implements resource.Resource.
func (g *GlobalPolicyComplianceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_global_policy_compliance"
}

// Read implements resource.Resource.
func (g *GlobalPolicyComplianceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.GlobalPolicyComplianceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	gpcSettingsDto, err := g.client.GetGlobalPolicyComplianceSettings(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"API Error while reading Global Policy Compliance data",
			fmt.Sprintf("Details: %s", err),
		)
		return
	}
	state = convertGlobalPolicyComplianceDTOToModel(gpcSettingsDto, g.client.Tenant)
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Schema implements resource.Resource.
func (g *GlobalPolicyComplianceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Allows configuration of the Global Policy Compliance settings",
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "This identifier is generated by the Terraform Provider automatically.",
				Computed:    true,
			},
			"access_policy_trust_provider_compliance": schema.StringAttribute{
				Description: "Defines a compliance requirement for an Access Policy to have a Trust Provider. Possible values are: Required, Recommended, Optional.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Required", "Recommended", "Optional"),
				},
			},
			"access_policy_access_condition_compliance": schema.StringAttribute{
				Description: "Defines a compliance requirement for an Access Policy to have an Access Condition. Possible values are: Required, Recommended, Optional.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Required", "Recommended", "Optional"),
				},
			},
			"agent_controller_trust_provider_compliance": schema.StringAttribute{
				Description: "Defines a compliance requirement for an Agent Controller to have a Trust Provider. Possible values are: Required, Recommended, Optional.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Required", "Recommended", "Optional"),
				},
			},
			"agent_controller_allowed_tls_hostname_compliance": schema.StringAttribute{
				Description: "Defines a compliance requirement for an Agent Controller to specify Allowed TLS Hostname parameter. Possible values are: Required, Recommended, Optional.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Required", "Recommended", "Optional"),
				},
			},
		},
	}
}

// Update implements resource.Resource.
func (g *GlobalPolicyComplianceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var currentModel models.GlobalPolicyComplianceModel
	diags := req.State.Get(ctx, &currentModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Retrieve values from plan
	var updatedModel models.GlobalPolicyComplianceModel
	diags = req.Plan.Get(ctx, &updatedModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := updateComplianceSettings(g.client, &updatedModel, &currentModel)

	if err != nil {
		resp.Diagnostics.AddError("Error updating Global Policy Compliance setting", "Error: "+err.Error())
		return
	}
	updatedModel.Id = types.StringValue(g.client.Tenant + "-gpc")
	diags = resp.State.Set(ctx, updatedModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func convertGlobalPolicyComplianceDTOToModel(gpcSettings *aembit.GlobalPolicyComplianceSettingsDTO, tenantId string) models.GlobalPolicyComplianceModel {
	model := models.GlobalPolicyComplianceModel{}
	model.Id = types.StringValue(tenantId + "-gpc")
	for _, settingDto := range *gpcSettings {
		switch settingDto.Name {
		case AccessPolicyTrustProviderComplianceSettingName:
			model.APTrustProviderCompliance = types.StringValue(settingDto.Value)
		case AccessPolicyAccessConditionComplianceSettingName:
			model.APAccessConditionCompliance = types.StringValue(settingDto.Value)
		case AgentControllerTrustProviderComplianceSettingName:
			model.ACTrustProviderCompliance = types.StringValue(settingDto.Value)
		case AgentControllerTlsHostNameComplianceSettingName:
			model.ACAllowedTLSHostanmeCompliance = types.StringValue(settingDto.Value)
		}
	}
	return model
}

func convertGlobalPolicyComplianceModelToDTO(model models.GlobalPolicyComplianceModel) aembit.GlobalPolicyComplianceSettingsDTO {
	dto := aembit.GlobalPolicyComplianceSettingsDTO{}
	dto = append(dto, aembit.TenantSettingDTO{
		Name:  AccessPolicyTrustProviderComplianceSettingName,
		Value: model.APTrustProviderCompliance.ValueString()})
	dto = append(dto, aembit.TenantSettingDTO{
		Name:  AccessPolicyAccessConditionComplianceSettingName,
		Value: model.APAccessConditionCompliance.ValueString()})
	dto = append(dto, aembit.TenantSettingDTO{
		Name:  AgentControllerTrustProviderComplianceSettingName,
		Value: model.ACTrustProviderCompliance.ValueString()})
	dto = append(dto, aembit.TenantSettingDTO{
		Name:  AgentControllerTlsHostNameComplianceSettingName,
		Value: model.ACAllowedTLSHostanmeCompliance.ValueString()})
	return dto
}

func updateComplianceSettings(client *aembit.CloudClient, currentModel *models.GlobalPolicyComplianceModel, previousModel *models.GlobalPolicyComplianceModel) error {
	if previousModel == nil {
		dto := convertGlobalPolicyComplianceModelToDTO(*currentModel)
		for _, settingDto := range dto {
			_, err := client.UpdateGlobalPolicyComplianceSetting(settingDto, nil)
			if err != nil {
				return err
			}
		}
	} else {
		//update only what has changed
		if !currentModel.APTrustProviderCompliance.Equal(previousModel.APTrustProviderCompliance) {
			_, err := client.UpdateGlobalPolicyComplianceSetting(aembit.TenantSettingDTO{Name: AccessPolicyTrustProviderComplianceSettingName, Value: currentModel.APTrustProviderCompliance.ValueString()}, nil)
			if err != nil {
				return err
			}
		}
		if !currentModel.APAccessConditionCompliance.Equal(previousModel.APAccessConditionCompliance) {
			_, err := client.UpdateGlobalPolicyComplianceSetting(aembit.TenantSettingDTO{Name: AccessPolicyAccessConditionComplianceSettingName, Value: currentModel.APAccessConditionCompliance.ValueString()}, nil)
			if err != nil {
				return err
			}
		}
		if !currentModel.ACTrustProviderCompliance.Equal(previousModel.ACTrustProviderCompliance) {
			_, err := client.UpdateGlobalPolicyComplianceSetting(aembit.TenantSettingDTO{Name: AgentControllerTrustProviderComplianceSettingName, Value: currentModel.ACTrustProviderCompliance.ValueString()}, nil)
			if err != nil {
				return err
			}
		}
		if !currentModel.ACAllowedTLSHostanmeCompliance.Equal(previousModel.ACAllowedTLSHostanmeCompliance) {
			_, err := client.UpdateGlobalPolicyComplianceSetting(aembit.TenantSettingDTO{Name: AgentControllerTlsHostNameComplianceSettingName, Value: currentModel.ACAllowedTLSHostanmeCompliance.ValueString()}, nil)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
