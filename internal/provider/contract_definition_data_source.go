package provider

import (
	"context"
	"fmt"

	"github.com/Think-iT-Labs/edc-connector-client-go/edc"
	"github.com/Think-iT-Labs/edc-connector-client-go/service/contractdefinition"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure ContractDefinitionDataSource fully satisfies interfaces defined by the terraform provider framework.
var _ datasource.DataSource = &ContractDefinitionDataSource{}

func NewContractDefinitionDataSource() datasource.DataSource {
	return &ContractDefinitionDataSource{}
}

// ContractDefinitionDataSource defines the data source implementation.
type ContractDefinitionDataSource struct {
	client *contractdefinition.Client
}

// ContractDefinitionDataSourceModel describes the data source data model.
type ContractDefinitionDataSourceModel struct {
	Id               types.String `tfsdk:"id"`
	AccessPolicyId   types.String `tfsdk:"access_policy_id"`
	ContractPolicyId types.String `tfsdk:"contract_policy_id"`
	Validity         types.Int64  `tfsdk:"validity"`
	Criteria         []Criterion  `tfsdk:"criteria"`
	CreatedAt        types.Int64  `tfsdk:"created_at"`
}

func (d *ContractDefinitionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contract_definition"
}

func (d *ContractDefinitionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Contract Definition Data Source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Contract definition identifier",
				Required:            true,
			},
			"access_policy_id": schema.StringAttribute{
				MarkdownDescription: "Access policy identifier",
				Computed:            true,
			},
			"contract_policy_id": schema.StringAttribute{
				MarkdownDescription: "Contract policy identifier",
				Computed:            true,
			},
			"validity": schema.Int64Attribute{
				MarkdownDescription: "Validity",
				Computed:            true,
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "Created at timestamp in seconds",
				Computed:            true,
			},
			"criteria": &schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"operand_left":  schema.StringAttribute{Computed: true},
						"operator":      schema.StringAttribute{Computed: true},
						"operand_right": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *ContractDefinitionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

func (d *ContractDefinitionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ContractDefinitionDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	cd, err := d.client.GetContractDefinition(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read contract definition, got error: %s", err))
		return
	}

	tflog.Info(ctx, fmt.Sprintf("read contract definition %s", cd.Id))
	data.AccessPolicyId = types.StringValue(cd.AccessPolicyId)
	data.ContractPolicyId = types.StringValue(cd.ContractPolicyId)
	data.Validity = types.Int64Value(cd.Validity)
	data.Criteria = criteriaModel(cd.Criteria)
	data.CreatedAt = types.Int64Value(cd.CreatedAt)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func criteriaModel(c []contractdefinition.Criterion) []Criterion {
	criteria := make([]Criterion, len(c))
	for i, criterion := range c {
		criteria[i] = Criterion{
			OperandLeft:  types.StringValue(criterion.OperandLeft),
			Operator:     types.StringValue(criterion.Operator),
			OperandRight: types.StringValue(criterion.OperandRight),
		}
	}
	return criteria
}
