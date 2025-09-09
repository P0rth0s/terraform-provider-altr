package datasources

import (
	"context"
	"fmt"
	"terraform-provider-altr/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &DataSourceResource{}
var _ resource.ResourceWithImportState = &DataSourceResource{}

func NewAccessManagementOltpPolicyDataResource() resource.Resource {
	return &DataSourceResource{}
}

// DataSourceResource is the data source implementation for fetching databases.
type DataSourceResource struct {
	client *client.Client
}

type DataSourceResourceModel struct {
	FriendlyDatabaseName types.String `tfsdk:"friendly_database_name"`
	DatabaseType         types.String `tfsdk:"database_type"`
	DatabaseName         types.String `tfsdk:"database_name"`
	DatabaseUsername     types.String `tfsdk:"database_username"`
	SFCount              types.Int64  `tfsdk:"sf_count"`
	InProgress           types.Int64  `tfsdk:"in_progress"`
	LastConnectedTime    types.String `tfsdk:"last_connected_time"`
	Hostname             types.String `tfsdk:"hostname"`
	DatabasePort         types.Int64  `tfsdk:"database_port"`
	ClientID             types.String `tfsdk:"client_id"`
	ID                   types.Int64  `tfsdk:"id"`
}

func NewDataSourceResource() resource.Resource {
	return &DataSourceResource{}
}

// Metadata returns the data source type name.
func (d *DataSourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "altr_datasource"
}

// ImportState imports the state of an existing resource.
func (d *DataSourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the data source.
func (d *DataSourceResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about databases connected through the ALTR platform.",
		Attributes: map[string]schema.Attribute{
			"databases": schema.ListNestedAttribute{
				Description: "List of databases.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                     schema.Int64Attribute{Description: "Database ID.", Computed: true},
						"client_id":              schema.StringAttribute{Description: "Client ID.", Computed: true},
						"friendly_database_name": schema.StringAttribute{Description: "Friendly database name.", Computed: true},
						"database_type":          schema.StringAttribute{Description: "Type of the database.", Computed: true},
						"database_name":          schema.StringAttribute{Description: "Name of the database.", Computed: true},
						"database_username":      schema.StringAttribute{Description: "Database username.", Computed: true},
						"sf_count":               schema.Int64Attribute{Description: "Number of columns protected.", Computed: true},
						"in_progress":            schema.Int64Attribute{Description: "Status of the database.", Computed: true},
						"last_connected_time":    schema.StringAttribute{Description: "Last connected time.", Computed: true},
					},
				},
				Computed: true,
			},
		},
	}
}

// Read fetches the data from the API and populates the Terraform state.
func (d *DataSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DataSourceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch databases from the API
	database, err := d.client.GetDataSource(state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to retrieve databases",
			"An error occurred while calling the ALTR API: "+err.Error(),
		)
		return
	}

	if database == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	d.mapPolicyToModel(database, &state)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (d *DataSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state struct {
		ID types.Int64 `tfsdk:"id"`
	}

	// Read the current state to get the database ID
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API to delete the database
	err := d.client.DeleteDataSource(state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Database",
			fmt.Sprintf("Could not delete database with ID %d: %s", state.ID.ValueInt64(), err.Error()),
		)
		return
	}

	// Remove the resource from the state
	resp.State.RemoveResource(ctx)
}

func (d *DataSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan struct {
		FriendlyDatabaseName types.String `tfsdk:"friendly_database_name"`
		DatabaseType         types.String `tfsdk:"database_type"`
		DatabaseName         types.String `tfsdk:"database_name"`
		DatabaseUsername     types.String `tfsdk:"database_username"`
		DatabasePassword     types.String `tfsdk:"database_password"`
		Hostname             types.String `tfsdk:"hostname"`
		DatabasePort         types.Int64  `tfsdk:"database_port"`
	}

	// Read the plan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the input for the API call
	input := client.CreateDataSourceInput{
		FriendlyDatabaseName: plan.FriendlyDatabaseName.ValueString(),
		DatabaseType:         plan.DatabaseType.ValueString(),
		DatabaseName:         plan.DatabaseName.ValueString(),
		DatabaseUsername:     plan.DatabaseUsername.ValueString(),
		DatabasePassword:     plan.DatabasePassword.ValueString(),
		Hostname:             plan.Hostname.ValueString(),
		DatabasePort:         plan.DatabasePort.ValueInt64(),
	}

	// Call the API to create the database
	dataSource, err := d.client.CreateDataSource(input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Database",
			fmt.Sprintf("Could not create database: %s", err.Error()),
		)
		return
	}

	// Map the response to the state
	var state struct {
		ID                   types.Int64  `tfsdk:"id"`
		FriendlyDatabaseName types.String `tfsdk:"friendly_database_name"`
		DatabaseType         types.String `tfsdk:"database_type"`
		DatabaseName         types.String `tfsdk:"database_name"`
		DatabaseUsername     types.String `tfsdk:"database_username"`
		Hostname             types.String `tfsdk:"hostname"`
		DatabasePort         types.Int64  `tfsdk:"database_port"`
	}

	state.ID = types.Int64Value(dataSource.ID)
	state.FriendlyDatabaseName = types.StringValue(dataSource.FriendlyDatabaseName)
	state.DatabaseType = types.StringValue(dataSource.DatabaseType)
	state.DatabaseName = types.StringValue(dataSource.DatabaseName)
	state.DatabaseUsername = types.StringValue(dataSource.DatabaseUsername)
	state.Hostname = types.StringValue(dataSource.Hostname)
	state.DatabasePort = types.Int64Value(dataSource.DatabasePort)

	// Set the state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *DataSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Not Implemented",
		"The Update operation is not implemented for this resource.",
	)
}

func (d *DataSourceResource) mapPolicyToModel(dataSource *client.DataSource, model *DataSourceResourceModel) {
	model.FriendlyDatabaseName = types.StringValue(dataSource.FriendlyDatabaseName)
	model.DatabaseType = types.StringValue(dataSource.DatabaseType)
	model.DatabaseName = types.StringValue(dataSource.DatabaseName)
	model.DatabaseUsername = types.StringValue(dataSource.DatabaseUsername)
	model.SFCount = types.Int64Value(dataSource.SFCount)
	model.InProgress = types.Int64Value(dataSource.InProgress)
	model.LastConnectedTime = types.StringValue(dataSource.LastConnectedTime)
	model.Hostname = types.StringValue(dataSource.Hostname)
	model.DatabasePort = types.Int64Value(dataSource.DatabasePort)
	model.ClientID = types.StringValue(dataSource.ClientID)
}
