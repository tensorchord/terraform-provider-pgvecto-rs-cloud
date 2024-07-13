package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/tensorchord/terraform-provider-pgvecto-rs-cloud/client"
)

const (
	defaultClusterCreateTimeout time.Duration = 5 * time.Minute
	defaultClusterUpdateTimeout time.Duration = 5 * time.Minute
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ClusterResource{}
var _ resource.ResourceWithConfigure = &ClusterResource{}
var _ resource.ResourceWithImportState = &ClusterResource{}

func NewClusterResource() resource.Resource {
	return &ClusterResource{}
}

// ClusterResource defines the resource implementation.
type ClusterResource struct {
	client *client.Client
}

func (r *ClusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

func (r *ClusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Cluster resource. This resource allows you to create a new PGVecto.rs cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Cluster identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_id": schema.StringAttribute{
				MarkdownDescription: "Default Account Identifier for the PGVecto.rs cloud",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cluster_name": schema.StringAttribute{
				MarkdownDescription: "The name of the cluster to be created. It is a string of no more than 32 characters.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"plan": schema.StringAttribute{
				MarkdownDescription: "The plan tier of the PGVecto.rs Cloud service. Available options are Starter and Enterprise.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"server_resource": schema.StringAttribute{
				MarkdownDescription: "The server resource of the cluster instance. Available aws-t3-xlarge-4c-16g, aws-m7i-large-2c-8g, aws-r7i-large-2c-16g,aws-r7i-xlarge-4c-32g",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The region of the cluster instance.Available options are us-east-1,eu-west-1",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_provider": schema.StringAttribute{
				MarkdownDescription: "The cloud provider of the cluster instance. At present, only aws is supported.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"database_name": schema.StringAttribute{
				MarkdownDescription: "The name of the database.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx,
				timeouts.Opts{
					Create: true,
					CreateDescription: `Timeout defaults to 5 mins. Accepts a string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) ` +
						`consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are ` +
						`"s" (seconds), "m" (minutes), "h" (hours).`,
					Update: true,
					UpdateDescription: `Timeout defaults to 5 mins. Accepts a string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) ` +
						`consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are ` +
						`"s" (seconds), "m" (minutes), "h" (hours).`,
				},
			),
		},
	}
}

func (r *ClusterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

func (r *ClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Create Cluster...")
	var data ClusterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	checkPlan := func(data ClusterResourceModel) (bool, error) {

		if data.Plan.IsNull() {
			return false, fmt.Errorf("Plan is required")
		}

		switch client.CNPGClusterPlan(data.Plan.ValueString()) {
		case client.CNPGClusterPlanStarter, client.CNPGClusterPlanEnterprise:
			return true, nil
		default:
			return false, fmt.Errorf("Invalid plan: %s", data.Plan.ValueString())
		}
	}

	if _, err := checkPlan(data); err != nil {
		resp.Diagnostics.AddError("check plan", err.Error())
		return
	}

	checkServerResource := func(data ClusterResourceModel) (bool, error) {
		if data.ServerResource.IsNull() {
			return false, fmt.Errorf("ServerResource is required")
		}

		switch client.ServerResource(data.ServerResource.ValueString()) {
		case client.ServerResourceAWST3XLarge, client.ServerResourceAWSM7ILarge, client.ServerResourceAWSR7ILarge, client.ServerResourceAWSR7IXLarge:
			return true, nil
		default:
			return false, fmt.Errorf("Invalid ServerResource: %s", data.ServerResource.ValueString())
		}
	}

	if _, err := checkServerResource(data); err != nil {
		resp.Diagnostics.AddError("check ServerResource", err.Error())
		return
	}

	var response *client.CNPGCluster
	var err error

	response, err = r.client.CreateCluster(client.CNPGClusterSpec{
		Name:           data.ClusterName.ValueString(),
		Plan:           client.CNPGClusterPlan(data.Plan.ValueString()),
		ServerResource: client.ServerResource(data.ServerResource.ValueString()),
		ClusterProvider: client.ClusterProvider{
			Type:   client.AWSCloudProvider,
			Region: data.Region.ValueString(),
		},
		PostgreSQLConfig: client.PostgreSQLConfig{
			PGDataDiskSize: data.PGDataDiskSize.ValueString(),
			VectorConfig: client.VectorConfig{
				DatabaseName: data.DatabaseName.ValueString(),
			},
		},
	}, data.AccountId.ValueString())

	if err != nil {
		err := client.Error{}
		if errors.As(err, &client.Error{}) {
			resp.Diagnostics.AddError("Failed to create cluster", err.Message)
			return
		}
		resp.Diagnostics.AddError("Failed to create cluster", err.Error())
		return
	}

	data.ClusterId = types.StringValue(response.Spec.ID)
	data.ClusterName = types.StringValue(response.Spec.Name)
	data.Plan = types.StringValue(string(response.Spec.Plan))
	data.ServerResource = types.StringValue(string(response.Spec.ServerResource))
	data.Region = types.StringValue(response.Spec.ClusterProvider.Region)
	data.ClusterProvider = types.StringValue(string(response.Spec.ClusterProvider.Type))
	data.Status = types.StringValue(string(response.Status.Status))
	data.ConnectEndpoint = types.StringValue(response.Status.VectorUserEndpoint)
	normalized := strings.TrimFunc(response.Spec.PostgreSQLConfig.PGDataDiskSize, func(r rune) bool {
		return r < '0' || r > '9'
	})
	data.PGDataDiskSize = types.StringValue(normalized)
	data.DatabaseName = types.StringValue(response.Spec.PostgreSQLConfig.VectorConfig.DatabaseName)

	// Wait for cluster to be RUNNING
	// Create() is passed a default timeout to use if no value
	// has been supplied in the Terraform configuration.
	createTimeout, diags := data.Timeouts.Create(ctx, defaultClusterCreateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.waitForStatus(ctx, createTimeout, r.client, string(client.CNPGClusterStatusReady))...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.refresh(r.client)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Read Cluster...")
	var state ClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(state.refresh(r.client)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Update Cluster...")

	var plan ClusterResourceModel
	var state ClusterResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Only support changes of cuSize - all other attributes are set to ForceNew
	_, err := r.client.UpgradeCluster(state.AccountId.String(), state.ClusterId.ValueString(), client.CNPGClusterUpgradeRequest{
		Plan:           client.CNPGClusterPlan(plan.Plan.ValueString()),
		ServerResource: client.ServerResource(plan.ServerResource.ValueString()),
		PGDataDiskSize: plan.PGDataDiskSize.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Failed to upgrade cluster", err.Error())
		return
	}

	// Wait for cluster to be RUNNING
	// Update() is passed a default timeout to use if no value
	// has been supplied in the Terraform configuration.
	updateTimeout, diags := plan.Timeouts.Update(ctx, defaultClusterUpdateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(state.waitForStatus(ctx, updateTimeout, r.client, string(client.CNPGClusterStatusReady))...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.refresh(r.client)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Delete Cluster...")
	var data ClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCluster(data.AccountId.String(), data.ClusterId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete cluster", err.Error())
		return
	}
}

func (r *ClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 1 || idParts[0] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: clusterId. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[0])...)
}

// ClusterResourceModel describes the resource data model.
type ClusterResourceModel struct {
	ClusterId       types.String   `tfsdk:"id"`
	AccountId       types.String   `tfsdk:"account_id"`
	ClusterName     types.String   `tfsdk:"cluster_name"`
	Plan            types.String   `tfsdk:"plan"`
	Region          types.String   `tfsdk:"region"`
	ServerResource  types.String   `tfsdk:"server_resource"`
	ClusterProvider types.String   `tfsdk:"cluster_provider"`
	Status          types.String   `tfsdk:"status"`
	ConnectEndpoint types.String   `tfsdk:"connect_endpoint"`
	PGDataDiskSize  types.String   `tfsdk:"pg_data_disk_size"`
	DatabaseName    types.String   `tfsdk:"database_name"`
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
}

func (data *ClusterResourceModel) refresh(client *client.Client) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	c, err := client.GetCluster(data.AccountId.String(), data.ClusterId.ValueString())
	if err != nil {
		diags.AddError("Client Error", fmt.Sprintf("Unable to GetCluster, got error: %s", err))
		return diags
	}

	// Save data into Terraform state
	data.ClusterId = types.StringValue(c.Spec.ID)
	data.ClusterName = types.StringValue(c.Spec.Name)
	data.Plan = types.StringValue(string(c.Spec.Plan))
	data.ServerResource = types.StringValue(string(c.Spec.ServerResource))
	data.Region = types.StringValue(c.Spec.ClusterProvider.Region)
	data.ClusterProvider = types.StringValue(string(c.Spec.ClusterProvider.Type))
	data.Status = types.StringValue(string(c.Status.Status))
	data.ConnectEndpoint = types.StringValue(c.Status.VectorUserEndpoint)
	data.PGDataDiskSize = types.StringValue(c.Spec.PostgreSQLConfig.PGDataDiskSize)
	data.DatabaseName = types.StringValue(c.Spec.PostgreSQLConfig.VectorConfig.DatabaseName)

	return diags
}

func (data *ClusterResourceModel) waitForStatus(ctx context.Context, timeout time.Duration, client *client.Client, status string) diag.Diagnostics {
	var diags diag.Diagnostics

	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		cluster, err := client.GetCluster(data.AccountId.String(), data.ClusterId.ValueString())
		if err != nil {
			return retry.NonRetryableError(err)
		}
		if string(cluster.Status.Status) != status {
			return retry.RetryableError(fmt.Errorf("cluster not yet in the %s state. Current state: %s", status, cluster.Status))
		}
		return nil
	})
	if err != nil {
		diags.AddError("Failed to wait for cluster to enter the RUNNING state.", err.Error())
	}

	return diags
}
