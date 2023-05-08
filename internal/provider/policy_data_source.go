package provider

import (
	"context"
	"fmt"

	"github.com/Think-iT-Labs/edc-connector-client/go/edc"
	"github.com/Think-iT-Labs/edc-connector-client/go/service/policies"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	policy, apiError, err := d.client.GetPolicy(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Policy, got error: %s", err))
		return
	}

	if apiError != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read Policy, got error: %v", *apiError[0].Message))
		return
	}

	tflog.Info(ctx, "read a data source")
	tflog.Info(ctx, "Policy", map[string]interface{}{
		"obj": policy,
	})
	// save into the Terraform state.
	data.Policy = toTFObject(policy.Policy)
	data.CreatedAt = types.Int64Value(policy.CreatedAt)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func toTFObject(policy policies.Policy) *Policy {

	return &Policy{
		// ExtensibleProperties: ,
		// UID:      basetypes.NewStringPointerValue(policy.UID),
		// Assignee: basetypes.NewStringValue(*policy.Assignee),
		// Assigner: basetypes.NewStringValue(*policy.Assigner),
		// Target: basetypes.NewStringValue(*policy.Target),
	}
}
