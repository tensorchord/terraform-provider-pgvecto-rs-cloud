// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tensorchord/terraform-provider-pgvecto-rs-cloud/client"
)

// Ensure PineconeProvider satisfies various provider interfaces.
var _ provider.Provider = &PGVectorsProvider{}

// PGVectorsProvider defines the provider implementation.
type PGVectorsProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// PGVectorsProviderModel describes the provider data model.
type PGVectorsProviderModel struct {
	ApiKey types.String `tfsdk:"api_key"`
	ApiURL types.String `tfsdk:"api_url"`
}

func (p *PGVectorsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {

	resp.TypeName = "pgvecto-rs-cloud"
	resp.Version = p.version
}

func (p *PGVectorsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `You can use the this Terraform provider to manage resources supported 
by [PGVecto.rs Cloud](https://cloud.pgvecto.rs). The provider must be configured with the proper 
credentials before use. You can provide credentials via the PGVECTORS_CLOUD_API_KEY environment variable.`,

		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "PGVecto.rs Cloud API Key. Can be configured by setting PGVECTORS_CLOUD_API_KEY environment variable.",
				Required:            true,
				Sensitive:           true,
			},
			"api_url": schema.StringAttribute{
				MarkdownDescription: "The URL of the PGVecto.rs Cloud API.",
				Optional:            true,
			},
		},
	}
}

func (p *PGVectorsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data PGVectorsProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Default to environment variables, but override
	// with Terraform configuration value if set.
	apiKey := os.Getenv("PGVECTORS_CLOUD_API_KEY")
	if !data.ApiKey.IsNull() {
		apiKey = data.ApiKey.ValueString()
	}

	apiUrl := os.Getenv("PGVECTORS_CLOUD_API_URL")
	if !data.ApiURL.IsNull() {
		apiUrl = data.ApiURL.ValueString()
	}

	client, err := client.NewClient(
		client.WithApiKey(apiKey),
		client.OverrideApiUrl(apiUrl),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to create pgvecto.rs cloud client: %v", err.Error())
		return
	}

	// PGVecto.rs Cloud client for data sources and resources
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *PGVectorsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewClusterResource,
	}
}

func (p *PGVectorsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewClusterDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PGVectorsProvider{
			version: version,
		}
	}
}
