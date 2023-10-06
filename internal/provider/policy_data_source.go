package provider

import (
	"context"
	"fmt"

	"github.com/Think-iT-Labs/edc-connector-client-go/edc"
	"github.com/Think-iT-Labs/edc-connector-client-go/service/policies"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AssetDataSource{}

func NewPolicyDataSource() datasource.DataSource {
	return &PolicyDataSource{}
}

// PolicyDataSource defines the data source implementation.
type PolicyDataSource struct {
	client *policies.Client
}

// AssetDataSourceModel describes the data source data model.
type PolicyDataSourceModel struct {
	Id        types.String `tfsdk:"id"`
	*Policy   `tfsdk:"policy"`
	CreatedAt types.Int64 `tfsdk:"created_at"`
}

func (d *PolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy"
}

func (d *PolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Policy data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Policy identifier",
				Required:            true,
			},
			"policy": PolicySchema(),
			"created_at": schema.Int64Attribute{
				Computed: true,
			},
		},
	}
}

func (d *PolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(*edc.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *edc.Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	client, err := policies.New(*cfg)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to initiate assets client",
			fmt.Sprintf("Client Error: %v", err),
		)
		return
	}

	d.client = client
}

func (d *PolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PolicyDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := d.client.GetPolicy(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Policy, got error: %s", err))
		return
	}

	tflog.Info(ctx, "read a data source")

	// save into the Terraform state.
	data.Policy = toTFPolicy(policy.Policy)
	data.CreatedAt = types.Int64Value(policy.CreatedAt)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a policy data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func toTFConstraint(constraint policies.Constraint) *Constraint {
	tfConstraint := &Constraint{}

	tfConstraint.EdcType = constraint.EdcType

	return tfConstraint
}

func toTFAction(action policies.Action) *Action {
	tfAction := &Action{}

	if action.IncludedIn != nil {
		tfAction.IncludedIn = basetypes.NewStringPointerValue(action.IncludedIn)
	}

	if action.ActionType != nil {
		tfAction.ActionType = basetypes.NewStringPointerValue(action.ActionType)
	}

	if action.Constraint != nil {
		tfAction.Constraint = toTFConstraint(*action.Constraint)
	}

	return tfAction

}

func toTFDuty(duty policies.Duty) *Duty {
	tfDuty := &Duty{}

	if duty.Assignee != nil {
		tfDuty.Assignee = basetypes.NewStringPointerValue(duty.Assignee)
	}

	if duty.Assigner != nil {
		tfDuty.Assigner = basetypes.NewStringPointerValue(duty.Assigner)
	}

	if duty.Target != nil {
		tfDuty.Target = basetypes.NewStringPointerValue(duty.Target)
	}

	if duty.UID != nil {
		tfDuty.UID = basetypes.NewStringPointerValue(duty.UID)
	}

	if duty.ParentPermission != nil {
		tfDuty.ParentPermission = toTFPermission(*duty.ParentPermission)
	}

	if duty.Action != nil {
		tfDuty.Action = toTFAction(*duty.Action)
	}

	if duty.Consequence != nil {
		tfDuty.Consequence = toTFDuty(*duty.Consequence)
	}

	if duty.Constraints != nil {
		var constraints []Constraint
		for _, constraint := range *duty.Constraints {
			constraints = append(constraints, *toTFConstraint(constraint))
		}
		tfDuty.Constraints = &constraints
	}

	return tfDuty

}

func toTFPermission(permission policies.Permission) *Permission {
	tfPermission := &Permission{}

	if permission.Assignee != nil {
		tfPermission.Assignee = basetypes.NewStringPointerValue(permission.Assignee)
	}

	if permission.Assigner != nil {
		tfPermission.Assigner = basetypes.NewStringPointerValue(permission.Assigner)
	}

	if permission.Target != nil {
		tfPermission.Target = basetypes.NewStringPointerValue(permission.Target)
	}

	if permission.UID != nil {
		tfPermission.UID = basetypes.NewStringPointerValue(permission.UID)
	}

	if permission.EdcType != nil {
		tfPermission.EdcType = basetypes.NewStringPointerValue(permission.EdcType)
	}

	if permission.Action != nil {
		tfPermission.Action = toTFAction(*permission.Action)
	}

	if permission.Constraints != nil {
		var constraints []Constraint
		for _, constraint := range *permission.Constraints {
			constraints = append(constraints, *toTFConstraint(constraint))
		}
		tfPermission.Constraints = &constraints
	}

	if permission.Duties != nil {
		var duties []Duty
		for _, duty := range *permission.Duties {
			duties = append(duties, *toTFDuty(duty))
		}
		tfPermission.Duties = &duties
	}

	return tfPermission

}

func toTFProhibition(prohibition policies.Prohibition) *Prohibition {
	tfProhibition := &Prohibition{}

	if prohibition.Assignee != nil {
		tfProhibition.Assignee = basetypes.NewStringPointerValue(prohibition.Assignee)
	}

	if prohibition.Assigner != nil {
		tfProhibition.Assigner = basetypes.NewStringPointerValue(prohibition.Assigner)
	}

	if prohibition.Target != nil {
		tfProhibition.Target = basetypes.NewStringPointerValue(prohibition.Target)
	}

	if prohibition.UID != nil {
		tfProhibition.UID = basetypes.NewStringPointerValue(prohibition.UID)
	}

	if prohibition.Action != nil {
		tfProhibition.Action = toTFAction(*prohibition.Action)
	}

	if prohibition.Constraints != nil {
		var constraints []Constraint
		for _, constraint := range *prohibition.Constraints {
			constraints = append(constraints, *toTFConstraint(constraint))
		}
		tfProhibition.Constraints = &constraints
	}

	return tfProhibition

}

func toTFPolicyType(policyType map[string]policies.PolicyType) basetypes.MapValue {
	var tfPolicyType = make(map[string]attr.Value, len(policyType))

	for k, v := range policyType {
		tfPolicyType[k] = types.StringValue(string(v))
	}

	return basetypes.NewMapValueMust(types.StringType, tfPolicyType)
}

func toTFPolicy(policy policies.Policy) *Policy {

	tfPolicy := &Policy{}

	if policy.Assignee != nil {
		tfPolicy.Assignee = basetypes.NewStringPointerValue(policy.Assignee)
	}

	if policy.Assigner != nil {
		tfPolicy.Assigner = basetypes.NewStringPointerValue(policy.Assigner)
	}

	if policy.ExtensibleProperties != nil {
		tfPolicy.ExtensibleProperties = (*ExtensibleProperties)(policy.ExtensibleProperties)
	}

	if policy.InheritsFrom != nil {
		tfPolicy.InheritsFrom = basetypes.NewStringPointerValue(policy.InheritsFrom)
	}

	if policy.Target != nil {
		tfPolicy.Target = basetypes.NewStringPointerValue(policy.Target)
	}

	if policy.UID != nil {
		tfPolicy.UID = basetypes.NewStringPointerValue(policy.UID)
	}

	if policy.Type != nil {
		tfPolicy.Type = toTFPolicyType(policy.Type)
	}

	if policy.Permissions != nil {
		var permissions []Permission
		for _, permission := range *policy.Permissions {
			permissions = append(permissions, *toTFPermission(permission))
		}
		tfPolicy.Permissions = &permissions
	}

	if policy.Obligations != nil {
		var obligations []Duty
		for _, obligation := range *policy.Obligations {
			obligations = append(obligations, *toTFDuty(obligation))
		}
		tfPolicy.Obligations = &obligations
	}

	if policy.Prohibitions != nil {
		var prohibitions []Prohibition
		for _, prohibition := range *policy.Prohibitions {
			prohibitions = append(prohibitions, *toTFProhibition(prohibition))
		}
		tfPolicy.Prohibitions = &prohibitions
	}

	return tfPolicy

}
