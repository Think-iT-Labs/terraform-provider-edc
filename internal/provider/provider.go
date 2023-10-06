package provider

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/Think-iT-Labs/edc-connector-client-go/config"
	"github.com/Think-iT-Labs/edc-connector-client-go/edc"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &EDCProvider{}

// EDCProvider defines the provider implementation.
type EDCProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// EDCProviderModel describes the provider data model.
type EDCProviderModel struct {
	Token     types.String `tfsdk:"token"`
	Addresses Addresses    `tfsdk:"addresses"`
}

type Addresses struct {
	ControlEndpoint    types.String `tfsdk:"control"`
	ManagementEndpoint types.String `tfsdk:"management"`
	ProtocolEndpoint   types.String `tfsdk:"protocol"`
	PublicEndpoint     types.String `tfsdk:"public"`
	DefaultEndpoint    types.String `tfsdk:"default"`
}

func (p *EDCProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "edc"
	resp.Version = p.version
}

func (p *EDCProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Optional: true,
			},
			"addresses": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"control": schema.StringAttribute{
						Optional: true,
					},
					"management": schema.StringAttribute{
						Optional: true,
					},
					"protocol": schema.StringAttribute{
						Optional: true,
					},
					"public": schema.StringAttribute{
						Optional: true,
					},
					"default": schema.StringAttribute{
						Optional: true,
					},
				},
			},
		},
	}
}

func (p *EDCProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data EDCProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	token, edcAddresses := validateProviderOptions(data, resp)

	if resp.Diagnostics.HasError() {
		return
	}

	cfg, err := config.LoadConfig(
		token,
		edcAddresses,
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create EDC API Client",
			"An unexpected error occurred when creating the EDC API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"EDC Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = cfg
	resp.ResourceData = cfg
}

func validateProviderOptions(data EDCProviderModel, resp *provider.ConfigureResponse) (string, edc.Addresses) {
	if data.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("Token"),
			"Unknown EDC API Token",
			"The provider cannot create the EDC API client as there is an unknown configuration value for the EDC API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the EDC_TOKEN environment variable.",
		)
	}

	if data.Addresses.ControlEndpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("Control"),
			"Unknown EDC Control Endpoint",
			"The provider cannot create the EDC Control Endpoint client as there is an unknown configuration value for the EDC Control Endpoint."+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the EDC_TOKEN environment variable.",
		)
	}

	if data.Addresses.DefaultEndpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("Default"),
			"Unknown EDC Default Endpoint",
			"The provider cannot create the EDC Default Endpoint client as there is an unknown configuration value for the EDC Default Endpoint."+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the EDC_ENDPOINT environment.",
		)
	}

	if data.Addresses.ManagementEndpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("Management"),
			"Unknown EDC Management Endpoint",
			"The provider cannot create the EDC Management Endpoint client as there is an unknown configuration value for the EDC Management Endpoint."+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the EDC_MANAGEMENT environment.",
		)
	}

	if data.Addresses.PublicEndpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("Public"),
			"Unknown EDC Public Endpoint",
			"The provider cannot create the EDC Public Endpoint client as there is an unknown configuration value for the EDC Public Endpoint."+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the EDC_PUBLIC environment.",
		)
	}

	if data.Addresses.ProtocolEndpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("Protocol"),
			"Unknown EDC Protocol Endpoint",
			"The provider cannot create the EDC Protocol Endpoint client as there is an unknown configuration value for the EDC Protocol Endpoint."+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the EDC_PROTOCOL environment.",
		)
	}

	token := os.Getenv("EDC_TOKEN")
	controlAddress := os.Getenv("EDC_CONTROL")
	publicAddress := os.Getenv("EDC_PUBLIC")
	protocolAddress := os.Getenv("EDC_PROTOCOL")
	managementAddress := os.Getenv("EDC_MANAGEMENT")
	defaultAddress := os.Getenv("EDC_DEFAULT")

	if !data.Token.IsNull() {
		token = data.Token.ValueString()
	}

	if !data.Addresses.ControlEndpoint.IsNull() {
		controlAddress = data.Addresses.ControlEndpoint.ValueString()
	}

	if !data.Addresses.PublicEndpoint.IsNull() {
		publicAddress = data.Addresses.ProtocolEndpoint.ValueString()
	}

	if !data.Addresses.ManagementEndpoint.IsNull() {
		managementAddress = data.Addresses.ManagementEndpoint.ValueString()
	}

	if !data.Addresses.ProtocolEndpoint.IsNull() {
		protocolAddress = data.Addresses.ProtocolEndpoint.ValueString()
	}

	if !data.Addresses.DefaultEndpoint.IsNull() {
		defaultAddress = data.Addresses.DefaultEndpoint.ValueString()
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing EDC Token",
			"The provider cannot create the EDC API client as there is a missing or empty value for the EDC API token. "+
				"Set the token value in the configuration or use the EDC_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	edcAddresses := edc.Addresses{
		Default:    &defaultAddress,
		Management: &managementAddress,
		Protocol:   &protocolAddress,
		Public:     &publicAddress,
		Control:    &controlAddress,
	}

	v := reflect.ValueOf(edcAddresses)
	typeOfAddresses := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Elem().String() == "" {
			fieldName := typeOfAddresses.Field(i).Name
			resp.Diagnostics.AddAttributeError(
				path.Root(fieldName),
				fmt.Sprintf("Missing EDC %s Address", fieldName),
				fmt.Sprintf("The provider cannot create the EDC API client as there is a missing or empty value for the EDC API %s. "+
					"Set the %s value in the configuration or use the EDC_%s environment variable. "+
					"If either is already set, ensure the value is not empty.", fieldName, fieldName, strings.ToUpper(fieldName)),
			)
		}
	}
	return token, edcAddresses
}

func (p *EDCProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAssetsResource,
		NewPoliciesResource,
		NewContractDefinitionResource,
	}
}

func (p *EDCProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAssetDataSource,
		NewPolicyDataSource,
		NewContractDefinitionDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &EDCProvider{
			version: version,
		}
	}
}
