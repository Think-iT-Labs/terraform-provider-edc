package provider

import (
	"context"
	"fmt"

	"github.com/Think-iT-Labs/edc-connector-client/go/service/policies"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type PolicyType string

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &PoliciesResource{}
var _ resource.ResourceWithImportState = &PoliciesResource{}

func NewPoliciesResource() resource.Resource {
	return &PoliciesResource{}
}

const (
	SetPolicyType      PolicyType = "set"
	OfferPolicyType    PolicyType = "offer"
	ContractPolicyType PolicyType = "contract"
	MaxRecursionLevel  int        = 3
)

type ExtensibleProperties map[string]string

type Constraint struct {
	EdcType types.String `tfsdk:"edctype"`
}

type Action struct {
	Constraint *Constraint  `tfsdk:"constraint"`
	IncludedIn types.String `tfsdk:"included_in"`
	ActionType types.String `tfsdk:"type"`
}

type Permission struct {
	Assignee    types.String  `tfsdk:"assignee"`
	Assigner    types.String  `tfsdk:"assigner"`
	Duties      *[]Duty       `tfsdk:"duties"`
	Target      types.String  `tfsdk:"target"`
	UID         types.String  `tfsdk:"uid"`
	Constraints *[]Constraint `tfsdk:"constraints"`
	Action      *Action       `tfsdk:"action"`
	EdcType     types.String  `tfsdk:"edctype"`
}

type Duty struct {
	Assignee         types.String  `tfsdk:"assignee"`
	Assigner         types.String  `tfsdk:"assigner"`
	Consequence      *Duty         `tfsdk:"consequence"`
	Target           types.String  `tfsdk:"target"`
	UID              types.String  `tfsdk:"uid"`
	Constraints      *[]Constraint `tfsdk:"constraints"`
	ParentPermission *Permission   `tfsdk:"parent_permission"`
	Action           *Action       `tfsdk:"action"`
}

type Prohibition struct {
	Assignee    types.String  `tfsdk:"assignee"`
	Assigner    types.String  `tfsdk:"assigner"`
	Target      types.String  `tfsdk:"target"`
	UID         types.String  `tfsdk:"uid"`
	Constraints *[]Constraint `tfsdk:"constraints"`
	Action      *Action       `tfsdk:"action"`
}

type Policy struct {
	UID                  types.String          `tfsdk:"uid"`
	Type                 map[string]PolicyType `tfsdk:"type"`
	Assignee             types.String          `tfsdk:"assignee"`
	Assigner             types.String          `tfsdk:"assigner"`
	ExtensibleProperties *ExtensibleProperties `tfsdk:"extensible_properties"`
	InheritsFrom         types.String          `tfsdk:"inherits_from"`
	Obligations          *[]Duty               `tfsdk:"obligations"`
	Permissions          *[]Permission         `tfsdk:"permissions"`
	Prohibitions         *[]Prohibition        `tfsdk:"prohibitions"`
	Target               types.String          `tfsdk:"target"`
}

type PolicyDefinition struct {
	Id        string `tfsdk:"id"`
	CreatedAt int64  `tfsdk:"created_at"`
	Policy    Policy `tfsdk:"policy"`
}

// PolicyResourceModel describes the resource data model.
type PolicyResourceModel struct {
	Policy `tfsdk:"policy"`
	Id     types.String `tfsdk:"id"`
}

func (p *PoliciesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy"
}

func (p *PoliciesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Policy resource",

		Attributes: map[string]schema.Attribute{
			"policy": PolicySchema(),
			"id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Policy identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func ConstraintSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"edctype": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func ActionSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"constraint": ConstraintSchema(),
			"included_in": schema.StringAttribute{
				Optional: true,
			},
			"type": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func PermissionSchema(level int) map[string]schema.Attribute {
	attributes := map[string]schema.Attribute{
		"assignee": schema.StringAttribute{
			Optional: true,
		},
		"assigner": schema.StringAttribute{
			Optional: true,
		},
		"target": schema.StringAttribute{
			Optional: true,
		},
		"uid": schema.StringAttribute{
			Optional: true,
		},
		"constraints": schema.ListAttribute{
			Optional:    true,
			ElementType: ConstraintSchema().GetType(),
		},
		"action": ActionSchema(),
		"edctype": schema.StringAttribute{
			Optional: true,
		},
	}
	if level != 0 {
		attributes["duties"] = schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: DutySchema(level - 1),
			},
		}
	}

	return attributes
}

func DutySchema(level int) map[string]schema.Attribute {
	attributes := map[string]schema.Attribute{
		"assignee": schema.StringAttribute{
			Optional: true,
		},
		"assigner": schema.StringAttribute{
			Optional: true,
		},
		"target": schema.StringAttribute{
			Optional: true,
		},
		"uid": schema.StringAttribute{
			Optional: true,
		},
		"constraints": schema.ListAttribute{
			Optional:    true,
			ElementType: ConstraintSchema().GetType(),
		},
		"parent_permission": schema.SingleNestedAttribute{
			Optional:   true,
			Attributes: PermissionSchema(level),
		},
		"action": ActionSchema(),
	}
	if level != 0 {
		attributes["consequence"] = schema.SingleNestedAttribute{
			Optional:   true,
			Attributes: DutySchema(level - 1),
		}
	}
	return attributes
}

func ProhibitionSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"assignee": schema.StringAttribute{
				Optional: true,
			},
			"assigner": schema.StringAttribute{
				Optional: true,
			},
			"target": schema.StringAttribute{
				Optional: true,
			},
			"uid": schema.StringAttribute{
				Optional: true,
			},
			"constraints": schema.ListAttribute{
				Optional:    true,
				ElementType: ConstraintSchema().GetType(),
			},
			"action": ActionSchema(),
		},
	}
}

// PolicySchema returns the schema to use for policy.
func PolicySchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"uid": schema.StringAttribute{
				Optional: true,
			},
			"type": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"assignee": schema.StringAttribute{
				Optional: true,
			},
			"assigner": schema.StringAttribute{
				Optional: true,
			},
			"extensible_properties": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"inherits_from": schema.StringAttribute{
				Optional: true,
			},
			"obligations": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: DutySchema(MaxRecursionLevel),
				},
			},
			"permissions": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: PermissionSchema(MaxRecursionLevel),
				},
			},
			"prohibitions": schema.ListAttribute{
				Optional:    true,
				ElementType: ProhibitionSchema().GetType(),
			},
			"target": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

// PoliciesResource defines the policies resource implementation.
type PoliciesResource struct {
	client *policies.Client
}

func (p *PoliciesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *PolicyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sdkObject := data.toSDKObject(ctx)

	// tflog.Debug(ctx, "SDK OBJECT", map[string]interface{}{
	// 	"piwpiw": sdkObject,
	// })
	output, apiError, err := p.client.CreatePolicy(*sdkObject)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Policy, got error: %s", err))
		return
	}

	if apiError != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to create Policy, got error: %s", err))
		return
	}

	// For the purposes of this Policies code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.StringValue(output.Id)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a policy")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *PoliciesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *PolicyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	ctx = tflog.SetField(ctx, "policy_id", data)

	if resp.Diagnostics.HasError() {
		return
	}

	policy, apiError, err := p.client.GetPolicy(data.Id.ValueString())

	tflog.Info(ctx, "Policy", map[string]any{
		"created_at": policy.CreatedAt,
		"id":         policy.Id,
		"policy":     policy.Policy,
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Policy with id %s, got error: %s", data.Id.String(), err))
		return
	}

	if apiError != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read Policy with id %s, got error: %s", data.Id.String(), err))
		return
	}

	// TODO: double check this
	// data.Policy = policy.Policy.Assignee

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *PoliciesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *PolicyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *PoliciesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *PolicyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiError, err := p.client.DeletePolicy(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Policy with id %s, got error: %s", data.Id.String(), err))
		return
	}

	if apiError != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to delete Policy with id %s, got error: %s", data.Id.String(), err))
		return
	}
}

func (p *PoliciesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (c *Constraint) toSDKObject() *policies.Constraint {
	return &policies.Constraint{
		EdcType: c.EdcType.ValueString(),
	}

}

func (a *Action) toSDKObject() *policies.Action {
	return &policies.Action{
		Constraint: a.Constraint.toSDKObject(),
		IncludedIn: a.IncludedIn.ValueStringPointer(),
		ActionType: a.ActionType.ValueStringPointer(),
	}
}

func (p *Permission) toSDKObject() *policies.Permission {

	var duties []policies.Duty
	if p.Duties != nil {
		for _, v := range *p.Duties {
			duties = append(duties, *v.toSDKObject())
		}
	}

	var constraints []policies.Constraint
	if p.Constraints != nil {
		for _, v := range *p.Constraints {
			constraints = append(constraints, *v.toSDKObject())
		}
	}
	return &policies.Permission{
		Assignee:    p.Assignee.ValueStringPointer(),
		Assigner:    p.Assigner.ValueStringPointer(),
		Target:      p.Target.ValueStringPointer(),
		UID:         p.UID.ValueStringPointer(),
		Action:      p.Action.toSDKObject(),
		Duties:      &duties,
		Constraints: &constraints,
		EdcType:     p.EdcType.ValueStringPointer(),
	}
}

func (d *Duty) toSDKObject() *policies.Duty {

	var constraints []policies.Constraint
	if d.Constraints != nil {
		for _, v := range *d.Constraints {
			constraints = append(constraints, *v.toSDKObject())
		}
	}

	return &policies.Duty{
		Assignee:         d.Assignee.ValueStringPointer(),
		Assigner:         d.Assigner.ValueStringPointer(),
		Target:           d.Target.ValueStringPointer(),
		UID:              d.UID.ValueStringPointer(),
		Constraints:      &constraints,
		Consequence:      d.Consequence.toSDKObject(),
		Action:           d.Action.toSDKObject(),
		ParentPermission: d.ParentPermission.toSDKObject(),
	}
}

func (d *Prohibition) toSDKObject() *policies.Prohibition {
	var constraints []policies.Constraint
	if d.Constraints != nil {
		for _, v := range *d.Constraints {
			constraints = append(constraints, *v.toSDKObject())
		}
	}

	return &policies.Prohibition{
		Assignee:    d.Assignee.ValueStringPointer(),
		Assigner:    d.Assigner.ValueStringPointer(),
		Target:      d.Target.ValueStringPointer(),
		UID:         d.UID.ValueStringPointer(),
		Action:      d.Action.toSDKObject(),
		Constraints: &constraints,
	}
}

func (p *PolicyResourceModel) toSDKObject(ctx context.Context) *policies.CreatePolicyInput {
	tflog.Debug(ctx, "transform tf object to sdk object", map[string]interface{}{
		"tf object": p.Policy,
	})
	var extensibleProperties policies.ExtensibleProperties
	var obligations []policies.Duty
	var prohibitions []policies.Prohibition
	var permissions []policies.Permission

	if p.Policy.ExtensibleProperties != nil {
		for k, v := range *p.Policy.ExtensibleProperties {
			extensibleProperties[k] = v
		}
	}

	if p.Obligations != nil {
		for _, v := range *p.Obligations {
			obligations = append(obligations, *v.toSDKObject())
		}
	}

	if p.Prohibitions != nil {
		for _, v := range *p.Prohibitions {
			prohibitions = append(prohibitions, *v.toSDKObject())
		}
	}

	if p.Permissions != nil {
		for _, v := range *p.Permissions {
			permissions = append(permissions, *v.toSDKObject())
		}
	}

	policy := policies.Policy{
		UID:                  p.Policy.UID.ValueStringPointer(),
		Assignee:             p.Policy.Assignee.ValueStringPointer(),
		Assigner:             p.Policy.Assigner.ValueStringPointer(),
		InheritsFrom:         p.Policy.InheritsFrom.ValueStringPointer(),
		Target:               p.Policy.Target.ValueStringPointer(),
		ExtensibleProperties: &extensibleProperties,
		Obligations:          &obligations,
		Prohibitions:         &prohibitions,
		Permissions:          &permissions,
	}
	return &policies.CreatePolicyInput{
		Id:     p.Id.ValueStringPointer(),
		Policy: policy,
	}
}
