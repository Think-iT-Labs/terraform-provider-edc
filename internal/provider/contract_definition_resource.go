package provider

import (
	"context"
	"fmt"

	"github.com/Think-iT-Labs/edc-connector-client-go/edc"
	"github.com/Think-iT-Labs/edc-connector-client-go/service/contractdefinition"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ContractDefinitionResource{}
var _ resource.ResourceWithImportState = &ContractDefinitionResource{}

func NewContractDefinitionResource() resource.Resource {
	return &ContractDefinitionResource{}
}

// ContractDefinitionResource defines the resource implementation.
type ContractDefinitionResource struct {
	client *contractdefinition.Client
}

// ContractDefinitionResourceModel describes the resource data model.
type ContractDefinitionResourceModel struct {
	Id               types.String `tfsdk:"id"`
	AccessPolicyId   types.String `tfsdk:"access_policy_id"`
	ContractPolicyId types.String `tfsdk:"contract_policy_id"`
	Validity         types.Int64  `tfsdk:"validity"`
	Criteria         []Criterion  `tfsdk:"criteria"`
}

type Criterion struct {
	OperandLeft  types.String `tfsdk:"operand_left"`
	Operator     types.String `tfsdk:"operator"`
	OperandRight types.String `tfsdk:"operand_right"`
}

func (r *ContractDefinitionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contract_definition"
}

func (r *ContractDefinitionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Contract definition resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Contract definition identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"access_policy_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Access policy identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"contract_policy_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Contract policy identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"validity": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Contract definition validity in seconds",
				PlanModifiers:       []planmodifier.Int64{},
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"criteria": CriteriaSchema(),
		},
	}
}

// CriteriaSchema returns the schema to use for tags.
func CriteriaSchema() *schema.ListNestedAttribute {
	return &schema.ListNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"operand_left":  schema.StringAttribute{Required: true},
				"operator":      schema.StringAttribute{Required: true},
				"operand_right": schema.StringAttribute{Optional: true},
			},
		},
	}
}

func (r *ContractDefinitionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	client, err := contractdefinition.New(*cfg)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to initiate contract definition client",
			fmt.Sprintf("Client Error: %v", err),
		)
		return
	}

	r.client = client
}

func (r *ContractDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ContractDefinitionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sdkObject := data.toSDKObject(ctx)
	output, err := r.client.CreateContractDefinition(*sdkObject)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create ContractDefinition, got error: %s", err))
		return
	}

	data.Id = types.StringValue(output.Id)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a contract definition")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ContractDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ContractDefinitionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	ctx = tflog.SetField(ctx, "contract_definition_id", data.Id.ValueString)

	if resp.Diagnostics.HasError() {
		return
	}

	cd, err := r.client.GetContractDefinition(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read contract definition, got error: %s", err))
		return
	}

	data.AccessPolicyId = types.StringValue(cd.AccessPolicyId)
	data.ContractPolicyId = types.StringValue(cd.ContractPolicyId)
	data.Validity = types.Int64Value(cd.Validity)
	data.Criteria = criteriaModel(cd.Criteria)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ContractDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ContractDefinitionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ContractDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ContractDefinitionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteContractDefinition(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete contract definition, got error: %s", err))
		return
	}
}

func (r *ContractDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ContractDefinitionResourceModel) toSDKObject(ctx context.Context) *contractdefinition.ContractDefinition {
	tflog.Debug(ctx, "transform tf object to sdk object", map[string]interface{}{
		"tf object": r,
	})

	return &contractdefinition.ContractDefinition{
		AccessPolicyId:   r.AccessPolicyId.ValueString(),
		ContractPolicyId: r.ContractPolicyId.ValueString(),
		Validity:         r.Validity.ValueInt64(),
		Criteria:         criteriaToSDKObject(r.Criteria),
	}
}

func criteriaToSDKObject(c []Criterion) []contractdefinition.Criterion {
	var sdkObject []contractdefinition.Criterion
	for _, criterion := range c {
		sdkObject = append(sdkObject, contractdefinition.Criterion{
			OperandLeft:  criterion.OperandLeft.ValueString(),
			Operator:     criterion.Operator.ValueString(),
			OperandRight: criterion.OperandRight.ValueString(),
		},
		)
	}
	return sdkObject
}
