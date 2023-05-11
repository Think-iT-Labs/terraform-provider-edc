package provider

import (
	"context"
	"fmt"

	"github.com/Think-iT-Labs/edc-connector-client/go/edc"
	"github.com/Think-iT-Labs/edc-connector-client/go/service/contractdefinition"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
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
	ContractDefinitionResourceModel
	CreatedAt int64 `tfsdk:"created_at"`
}

func (d *ContractDefinitionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contract_definition"
}

func (d *ContractDefinitionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Asset data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Asset identifier",
				Required:            true,
			},
			"access_policy_id": schema.StringAttribute{
				MarkdownDescription: "Access policy identifier",
				Required:            true,
			},
			"contract_policy_id": schema.StringAttribute{
				MarkdownDescription: "Contract policy identifier",
				Required:            true,
			},
			"validity": schema.Int64Attribute{
				MarkdownDescription: "Validity",
				Required:            true,
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "Created at",
				Required:            true,
			},
			"criteria": CriteriaSchema(),
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
			"Failed to initiate assets client",
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

	tflog.Info(ctx, "read contract definition")
	// For the purposes of this Asset code, hardcoding a response value to
	// save into the Terraform state.
	data.AccessPolicyId = types.StringValue(cd.AccessPolicyId)
	data.ContractPolicyId = types.StringValue(cd.ContractPolicyId)
	data.Validity = types.Int64Value(cd.Validity)
	data.Criteria = CriteriaModel(cd.Criteria)
	// TODO add created_at when it is fixed

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func CriteriaModel(c []contractdefinition.Criterion) []Criterion {
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
