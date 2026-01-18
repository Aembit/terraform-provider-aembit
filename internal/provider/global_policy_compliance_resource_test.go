package provider

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"terraform-provider-aembit/internal/provider/models"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_GPC_CreateImportUpdate(t *testing.T) {
	const gpcResourceName string = "aembit_global_policy_compliance.test"
	const gpcResourceDef string = `provider "aembit" {}

		resource "aembit_global_policy_compliance" "test" {
			access_policy_trust_provider_compliance = "Recommended"
			access_policy_access_condition_compliance = "Recommended"
			agent_controller_trust_provider_compliance = "Recommended"
			agent_controller_allowed_tls_hostname_compliance = "Recommended"
		}`
	const gpcResourceUpdate string = `provider "aembit" {}

		resource "aembit_global_policy_compliance" "test" {
			access_policy_trust_provider_compliance = "Optional"
			access_policy_access_condition_compliance = "Optional"
			agent_controller_trust_provider_compliance = "Optional"
			agent_controller_allowed_tls_hostname_compliance = "Optional"
		}`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				ResourceName: gpcResourceName,
				Config:       gpcResourceDef,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						gpcResourceName,
						"access_policy_trust_provider_compliance",
						"Recommended",
					),
					resource.TestCheckResourceAttr(
						gpcResourceName,
						"access_policy_access_condition_compliance",
						"Recommended",
					),
					resource.TestCheckResourceAttr(
						gpcResourceName,
						"agent_controller_trust_provider_compliance",
						"Recommended",
					),
					resource.TestCheckResourceAttr(
						gpcResourceName,
						"agent_controller_allowed_tls_hostname_compliance",
						"Recommended",
					),
				),
			},
			// ImportState testing
			{ResourceName: gpcResourceName, ImportState: true, ImportStateVerify: true},
			// Update
			{
				ResourceName: gpcResourceName,
				Config:       gpcResourceUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						gpcResourceName,
						"access_policy_trust_provider_compliance",
						"Optional",
					),
					resource.TestCheckResourceAttr(
						gpcResourceName,
						"access_policy_access_condition_compliance",
						"Optional",
					),
					resource.TestCheckResourceAttr(
						gpcResourceName,
						"agent_controller_trust_provider_compliance",
						"Optional",
					),
					resource.TestCheckResourceAttr(
						gpcResourceName,
						"agent_controller_allowed_tls_hostname_compliance",
						"Optional",
					),
				),
			},
		},
	})
}

func Test_GPC_convertGlobalPolicyComplianceDTOToModel(t *testing.T) {
	type args struct {
		gpcSettings *aembit.GlobalPolicyComplianceSettingsDTO
	}
	tests := []struct {
		name string
		args args
		want models.GlobalPolicyComplianceModel
	}{
		{
			name: "success",
			args: args{gpcSettings: &aembit.GlobalPolicyComplianceSettingsDTO{
				aembit.TenantSettingDTO{
					Name:  AccessPolicyTrustProviderComplianceSettingName,
					Value: "Required",
				},
				aembit.TenantSettingDTO{
					Name:  AccessPolicyAccessConditionComplianceSettingName,
					Value: "Required",
				},
				aembit.TenantSettingDTO{
					Name:  AgentControllerTrustProviderComplianceSettingName,
					Value: "Required",
				},
				aembit.TenantSettingDTO{
					Name:  AgentControllerTlsHostNameComplianceSettingName,
					Value: "Required",
				},
			}},
			want: models.GlobalPolicyComplianceModel{
				Id: types.StringValue(
					"testId-gpc",
				), APTrustProviderCompliance: types.StringValue("Required"), APAccessConditionCompliance: types.StringValue("Required"),
				ACTrustProviderCompliance: types.StringValue(
					"Required",
				), ACAllowedTLSHostanmeCompliance: types.StringValue("Required"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := convertGlobalPolicyComplianceDTOToModel(tt.args.gpcSettings, "testId"); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("convertGlobalPolicyComplianceDTOToModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GPC_convertGlobalPolicyComplianceModelToDTO(t *testing.T) {
	tests := []struct {
		name  string
		model models.GlobalPolicyComplianceModel
		want  aembit.GlobalPolicyComplianceSettingsDTO
	}{
		{
			name: "success",
			model: models.GlobalPolicyComplianceModel{
				APTrustProviderCompliance:      types.StringValue("Required"),
				APAccessConditionCompliance:    types.StringValue("Required"),
				ACTrustProviderCompliance:      types.StringValue("Required"),
				ACAllowedTLSHostanmeCompliance: types.StringValue("Required"),
			},
			want: aembit.GlobalPolicyComplianceSettingsDTO{
				aembit.TenantSettingDTO{
					Name:  AccessPolicyTrustProviderComplianceSettingName,
					Value: "Required",
				},
				aembit.TenantSettingDTO{
					Name:  AccessPolicyAccessConditionComplianceSettingName,
					Value: "Required",
				},
				aembit.TenantSettingDTO{
					Name:  AgentControllerTrustProviderComplianceSettingName,
					Value: "Required",
				},
				aembit.TenantSettingDTO{
					Name:  AgentControllerTlsHostNameComplianceSettingName,
					Value: "Required",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertGlobalPolicyComplianceModelToDTO(tt.model); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("convertGlobalPolicyComplianceModelToDTO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GPC_updateComplianceSettings_API_Errors(t *testing.T) {

	type args struct {
		currentModel  *models.GlobalPolicyComplianceModel
		previousModel *models.GlobalPolicyComplianceModel
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "api_error_no_previous_state",
			args: args{
				currentModel: &models.GlobalPolicyComplianceModel{},
				// this indicates that there is no previous state
				previousModel: nil,
			},
			wantErr: true,
		},
		{
			name: "api_error_when_updating_ap_tp",
			args: args{
				currentModel: &models.GlobalPolicyComplianceModel{
					APTrustProviderCompliance: types.StringValue("Value"),
				},
				previousModel: &models.GlobalPolicyComplianceModel{
					APTrustProviderCompliance: types.StringValue("ADifferentValue"),
				},
			},
			wantErr: true,
		},
		{
			name: "api_error_when_updating_ap_ac",
			args: args{
				currentModel: &models.GlobalPolicyComplianceModel{
					APAccessConditionCompliance: types.StringValue("Value"),
				},
				previousModel: &models.GlobalPolicyComplianceModel{
					APAccessConditionCompliance: types.StringValue("ADifferentValue"),
				},
			},
			wantErr: true,
		},
		{
			name: "api_error_when_updating_ac_tp",
			args: args{
				currentModel: &models.GlobalPolicyComplianceModel{
					ACTrustProviderCompliance: types.StringValue("Value"),
				},
				previousModel: &models.GlobalPolicyComplianceModel{
					ACTrustProviderCompliance: types.StringValue("ADifferentValue"),
				},
			},
			wantErr: true,
		},
		{
			name: "api_error_when_updating_ap_tls",
			args: args{
				currentModel: &models.GlobalPolicyComplianceModel{
					ACAllowedTLSHostanmeCompliance: types.StringValue("Value"),
				},
				previousModel: &models.GlobalPolicyComplianceModel{
					ACAllowedTLSHostanmeCompliance: types.StringValue("ADifferentValue"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
				}),
			)
			defer server.Close()

			c, _ := aembit.NewClient(&aembit.URLBuilder{}, nil, "", "test")
			if err := updateComplianceSettings(c, tt.args.currentModel, tt.args.previousModel); (err != nil) != tt.wantErr {
				t.Errorf("updateComplianceSettings() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
