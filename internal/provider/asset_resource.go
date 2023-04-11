package provider

import (
	"context"
	"fmt"

	"github.com/Think-iT-Labs/edc-connector-client/go/edc"
	"github.com/Think-iT-Labs/edc-connector-client/go/service/assets"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AssetsResource{}
var _ resource.ResourceWithImportState = &AssetsResource{}

func NewAssetsResource() resource.Resource {
	return &AssetsResource{}
}

// AssetsResource defines the resource implementation.
type AssetsResource struct {
	client *assets.Client
}

// AssetsResourceModel describes the resource data model.
type AssetsResourceModel struct {
	AssetProperties `tfsdk:"asset"`
	DataAddress     `tfsdk:"data"`
	Id              types.String `tfsdk:"id"`
}

type AssetProperties map[string]string

type DataAddress struct {
	HttpDataAddress         *HttpDataAddress         `tfsdk:"http"`
	S3StorageDataAddress    *S3StorageDataAddress    `tfsdk:"s3"`
	AzureStorageDataAddress *AzureStorageDataAddress `tfsdk:"azure"`
}

type HttpDataAddress struct {
	Type             types.String `tfsdk:"type"`
	Name             types.String `tfsdk:"name"`
	Path             types.String `tfsdk:"path"`
	Method           types.String `tfsdk:"method"`
	BaseUrl          types.String `tfsdk:"base_url"`
	AuthKey          types.String `tfsdk:"auth_key"`
	AuthCode         types.String `tfsdk:"auth_code"`
	SecretName       types.String `tfsdk:"secret_name"`
	ProxyBody        types.String `tfsdk:"proxy_body"`
	ProxyPath        types.String `tfsdk:"proxy_path"`
	ProxyQueryParams types.String `tfsdk:"proxy_query_params"`
	ProxyMethod      types.String `tfsdk:"proxy_method"`
	ContentType      types.String `tfsdk:"content_type"`
}

type S3StorageDataAddress struct {
	Type            types.String `tfsdk:"type"`
	Name            types.String `tfsdk:"name"`
	BucketName      types.String `tfsdk:"bucket_name"`
	AccessKeyId     types.String `tfsdk:"access_key_id"`
	SecretAccessKey types.String `tfsdk:"secret_access_key"`
}

type AzureStorageDataAddress struct {
	Type      types.String `tfsdk:"type"`
	Container types.String `tfsdk:"container"`
	Account   types.String `tfsdk:"account"`
	BlobName  types.String `tfsdk:"blob_name"`
}

func (r *AssetsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asset"
}

func (r *AssetsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Assets resource",

		Attributes: map[string]schema.Attribute{
			"asset": AssetsSchema(),
			"data":  DataAssetsSchema(),
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Assets identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// AssetsSchema returns the schema to use for tags.
func AssetsSchema() *schema.MapAttribute {
	return &schema.MapAttribute{
		Required:    true,
		ElementType: types.StringType,
	}
}

// DataAssetsSchema returns the schema to use fo tags.
func DataAssetsSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"s3":    S3Schema(),
			"http":  HTTPSchema(),
			"azure": AzureSchema(),
		},
	}
}

func S3Schema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Optional: true,
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
			"bucket_name": schema.StringAttribute{
				Optional: true,
			},
			"access_key_id": schema.StringAttribute{
				Optional: true,
			},
			"secret_access_key": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func AzureSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Optional: true,
			},
			"container": schema.StringAttribute{
				Optional: true,
			},
			"account": schema.StringAttribute{
				Optional: true,
			},
			"blob_name": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func HTTPSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Optional: true,
			},
			"path": schema.StringAttribute{
				Optional: true,
			},
			"method": schema.StringAttribute{
				Optional: true,
			},
			"base_url": schema.StringAttribute{
				Optional: true,
			},
			"auth_key": schema.StringAttribute{
				Optional: true,
			},
			"auth_code": schema.StringAttribute{
				Optional: true,
			},
			"secret_name": schema.StringAttribute{
				Optional: true,
			},
			"proxy_body": schema.StringAttribute{
				Optional: true,
			},
			"proxy_path": schema.StringAttribute{
				Optional: true,
			},
			"proxy_query_params": schema.StringAttribute{
				Optional: true,
			},
			"proxy_method": schema.StringAttribute{
				Optional: true,
			},
			"content_type": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (r *AssetsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	client, err := assets.New(*cfg)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to initiate assets client",
			fmt.Sprintf("Client Error: %v", err),
		)
		return
	}

	r.client = client
}

func (r *AssetsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *AssetsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	output, err := r.client.CreateAsset(*data.toSDKObject())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Assets, got error: %s", err))
		return
	}

	// For the purposes of this Assets code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.StringValue(output.Id)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AssetsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *AssetsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	asset, err := r.client.GetAsset(data.Id.String())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Assets, got error: %s", err))
		return
	}

	data.AssetProperties = AssetProperties(asset.AssetProperties)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AssetsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *AssetsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update Assets, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AssetsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *AssetsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAsset(data.Id.String())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Assets, got error: %s", err))
		return
	}
}

func (r *AssetsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *AssetsResourceModel) toSDKObject() *assets.CreateAssetInput {
	assetProperties := assets.AssetProperties(r.AssetProperties)

	tflog.Info(context.TODO(), "address", map[string]interface{}{
		"data": r.DataAddress,
	})
	dataAddress := assets.DataAddress{}

	if r.DataAddress.HttpDataAddress != nil {
		dataAddress.HttpDataAddress = &assets.HttpData{
			Type:             r.DataAddress.HttpDataAddress.Type.ValueString(),
			Name:             r.DataAddress.HttpDataAddress.Name.ValueStringPointer(),
			Path:             r.DataAddress.HttpDataAddress.Path.ValueStringPointer(),
			Method:           r.DataAddress.HttpDataAddress.Method.ValueStringPointer(),
			BaseUrl:          r.DataAddress.HttpDataAddress.BaseUrl.ValueStringPointer(),
			AuthKey:          r.DataAddress.HttpDataAddress.AuthKey.ValueStringPointer(),
			SecretName:       r.DataAddress.HttpDataAddress.SecretName.ValueStringPointer(),
			AuthCode:         r.DataAddress.HttpDataAddress.AuthCode.ValueStringPointer(),
			ProxyBody:        r.DataAddress.HttpDataAddress.ProxyBody.ValueStringPointer(),
			ProxyPath:        r.DataAddress.HttpDataAddress.ProxyPath.ValueStringPointer(),
			ProxyQueryParams: r.DataAddress.HttpDataAddress.ProxyQueryParams.ValueStringPointer(),
			ProxyMethod:      r.DataAddress.HttpDataAddress.ProxyMethod.ValueStringPointer(),
			ContentType:      r.DataAddress.HttpDataAddress.ContentType.ValueStringPointer(),
		}
	}

	if r.DataAddress.S3StorageDataAddress != nil {
		dataAddress.S3StorageDataAddress = &assets.S3Data{
			Type:            r.DataAddress.S3StorageDataAddress.Type.ValueString(),
			Name:            r.DataAddress.S3StorageDataAddress.Name.ValueStringPointer(),
			BucketName:      r.DataAddress.S3StorageDataAddress.BucketName.ValueStringPointer(),
			AccessKeyId:     r.DataAddress.S3StorageDataAddress.AccessKeyId.ValueStringPointer(),
			SecretAccessKey: r.DataAddress.S3StorageDataAddress.SecretAccessKey.ValueStringPointer(),
		}
	}

	if r.DataAddress.AzureStorageDataAddress != nil {
		dataAddress.AzureStorageDataAddress = &assets.AzureData{
			Type:      r.DataAddress.AzureStorageDataAddress.Type.ValueString(),
			Container: r.DataAddress.AzureStorageDataAddress.Container.ValueStringPointer(),
			Account:   r.DataAddress.AzureStorageDataAddress.Account.ValueStringPointer(),
			BlobName:  r.DataAddress.AzureStorageDataAddress.BlobName.ValueStringPointer(),
		}
	}
	// dataAddress := assets.DataAddress{
	// 	HttpDataAddress: &assets.HttpData{
	// 		Type:             r.DataAddress.HttpDataAddress.Type.ValueString(),
	// 		Name:             r.DataAddress.HttpDataAddress.Name.ValueStringPointer(),
	// 		Path:             r.DataAddress.HttpDataAddress.Path.ValueStringPointer(),
	// 		Method:           r.DataAddress.HttpDataAddress.Method.ValueStringPointer(),
	// 		BaseUrl:          r.DataAddress.HttpDataAddress.BaseUrl.ValueStringPointer(),
	// 		AuthKey:          r.DataAddress.HttpDataAddress.AuthKey.ValueStringPointer(),
	// 		SecretName:       r.DataAddress.HttpDataAddress.SecretName.ValueStringPointer(),
	// 		AuthCode:         r.DataAddress.HttpDataAddress.AuthCode.ValueStringPointer(),
	// 		ProxyBody:        r.DataAddress.HttpDataAddress.ProxyBody.ValueStringPointer(),
	// 		ProxyPath:        r.DataAddress.HttpDataAddress.ProxyPath.ValueStringPointer(),
	// 		ProxyQueryParams: r.DataAddress.HttpDataAddress.ProxyQueryParams.ValueStringPointer(),
	// 		ProxyMethod:      r.DataAddress.HttpDataAddress.ProxyMethod.ValueStringPointer(),
	// 		ContentType:      r.DataAddress.HttpDataAddress.ContentType.ValueStringPointer(),
	// 	},
	// 	S3StorageDataAddress: &assets.S3Data{
	// 		Type:            r.DataAddress.S3StorageDataAddress.Type.ValueString(),
	// 		Name:            r.DataAddress.HttpDataAddress.Name.ValueStringPointer(),
	// 		BucketName:      r.DataAddress.S3StorageDataAddress.BucketName.ValueStringPointer(),
	// 		AccessKeyId:     r.DataAddress.S3StorageDataAddress.AccessKeyId.ValueStringPointer(),
	// 		SecretAccessKey: r.DataAddress.S3StorageDataAddress.SecretAccessKey.ValueStringPointer(),
	// 	},
	// 	AzureStorageDataAddress: &assets.AzureData{
	// 		Type:      r.DataAddress.AzureStorageDataAddress.Type.ValueString(),
	// 		Container: r.DataAddress.AzureStorageDataAddress.Container.ValueStringPointer(),
	// 		Account:   r.DataAddress.AzureStorageDataAddress.Account.ValueStringPointer(),
	// 		BlobName:  r.DataAddress.AzureStorageDataAddress.BlobName.ValueStringPointer(),
	// 	},
	// }

	return &assets.CreateAssetInput{
		Asset: assets.Asset{
			AssetProperties: assetProperties,
		},
		DataAddress: dataAddress,
	}

}
