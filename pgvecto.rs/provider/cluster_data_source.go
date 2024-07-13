package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/tensorchord/terraform-provider-pgvecto-rs-cloud/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ClusterDataSource{}

func NewClusterDataSource() datasource.DataSource {
	return &ClusterDataSource{}
}

// ClusterDataSource defines the data source implementation.
type ClusterDataSource struct {
	client *client.Client
}

// ClusterDataSourceModel describes the cluster data model.
type ClusterDataSourceModel struct {
	ClusterId       types.String `tfsdk:"id"`
	AccountId       types.String `tfsdk:"account_id"`
	ClusterName     types.String `tfsdk:"cluster_name"`
	Plan            types.String `tfsdk:"plan"`
	Region          types.String `tfsdk:"region"`
	ServerResource  types.String `tfsdk:"server_resource"`
	ClusterProvider types.String `tfsdk:"cluster_provider"`
	Status          types.String `tfsdk:"status"`
	ConnectEndpoint types.String `tfsdk:"connect_endpoint"`
	PGDataDiskSize  types.String `tfsdk:"pg_data_disk_size"`
	DatabaseName    types.String `tfsdk:"database_name"`
}

func (d *ClusterDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

func (r *ClusterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Cluster Data Source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Cluster identifier",
				Required:            true,
			},
			"account_id": schema.StringAttribute{
				MarkdownDescription: "Default Account Identifier for the PGVecto.rs cloud",
				Required:            true,
			},
			"cluster_name": schema.StringAttribute{
				MarkdownDescription: "The name of the cluster to be created. It is a string of no more than 32 characters.",
				Computed:            true,
			},
			"plan": schema.StringAttribute{
				MarkdownDescription: "The plan tier of the PGVecto.rs Cloud service. Available options are Starter and Enterprise.",
				Optional:            true,
			},
			"server_resource": schema.StringAttribute{
				MarkdownDescription: "The server resource of the cluster instance. Avaliable aws-t3-xlarge-4c-16g, aws-m7i-large-2c-8g, aws-r7i-large-2c-16g,aws-r7i-xlarge-4c-32g",
				Computed:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The region of the cluster instance.Avaliable options are us-east-1,eu-west-1",
				Computed:            true,
			},
			"cluster_provider": schema.StringAttribute{
				MarkdownDescription: "The cloud provider of the cluster instance. At present, only aws is supported.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The current status of the cluster. Possible values are Initializing, Ready, NotReady, Deleted, Upgrading, Suspended, Resuming.",
				Computed:            true,
			},
			"connect_endpoint": schema.StringAttribute{
				MarkdownDescription: "The psql connection endpoint of the cluster.",
				Computed:            true,
			},
			"pg_data_disk_size": schema.StringAttribute{
				MarkdownDescription: "The size of the PGData disk in GB, please insert between 1 and 16384.",
				Computed:            true,
			},
			"database_name": schema.StringAttribute{
				MarkdownDescription: "The name of the database.",
				Computed:            true,
			},
		},
	}
}

func (d *ClusterDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *ClusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ClusterDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "sending describe project request...")
	if state.ClusterId.IsNull() {
		resp.Diagnostics.AddError("Invalid Cluster ID", "Cluster ID is required")
		return
	}

	if state.AccountId.IsNull() {
		resp.Diagnostics.AddError("Invalid Account ID", "Account ID is required")
		return
	}
	c, err := d.client.GetCluster(state.AccountId.ValueString(), state.ClusterId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to GetCluster %s, got error: %v", state.ClusterId.ValueString(), err))
		return
	}

	// Save data into Terraform state
	state.ClusterId = types.StringValue(c.Spec.ID)
	state.ClusterName = types.StringValue(c.Spec.Name)
	state.Plan = types.StringValue(string(c.Spec.Plan))
	state.ServerResource = types.StringValue(string(c.Spec.ServerResource))
	state.Region = types.StringValue(string(c.Spec.ClusterProvider.Region))
	state.ClusterProvider = types.StringValue(string(c.Spec.ClusterProvider.Type))
	state.Status = types.StringValue(string(c.Status.Status))
	state.ConnectEndpoint = types.StringValue(c.Status.VectorUserEndpoint)
	state.PGDataDiskSize = types.StringValue(c.Spec.PostgreSQLConfig.PGDataDiskSize)
	state.DatabaseName = types.StringValue(c.Spec.PostgreSQLConfig.VectorConfig.DatabaseName)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
