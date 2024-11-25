// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/henrywhitaker3/terraform-provider-cronitor/pkg/cronitor"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &MonitorDataSource{}

func NewMonitorDataSource() datasource.DataSource {
	return &MonitorDataSource{}
}

// MonitorDataSource defines the data source implementation.
type MonitorDataSource struct {
	client *cronitor.Client
}

type RequestModel struct {
	Url             types.String            `tfsdk:"url"`
	Headers         map[string]types.String `tfsdk:"headers"`
	Cookies         map[string]types.String `tfsdk:"cookies"`
	Body            types.String            `tfsdk:"body"`
	Method          types.String            `tfsdk:"method"`
	TimeoutSeconds  types.Int32             `tfsdk:"timeout_seconds"`
	Regions         []types.String          `tfsdk:"regions"`
	FollowRedirects types.Bool              `tfsdk:"follow_redirects"`
	VerifySsl       types.Bool              `tfsdk:"verify_ssl"`
}

// MonitorDataSourceModel describes the data source data model.
type MonitorModel struct {
	Id                types.String   `tfsdk:"id"`
	Key               types.String   `tfsdk:"key"`
	Name              types.String   `tfsdk:"name"`
	Assertions        []types.String `tfsdk:"assertions"`
	Disabled          types.Bool     `tfsdk:"disabled"`
	FailureTolerance  types.Int32    `tfsdk:"failure_tolerance"`
	GraceSeconds      types.Int32    `tfsdk:"grace_seconds"`
	Notify            []types.String `tfsdk:"notify"`
	Paused            types.Bool     `tfsdk:"paused"`
	Platform          types.String   `tfsdk:"platform"`
	RealertInterval   types.String   `tfsdk:"realert_interval"`
	Request           *RequestModel  `tfsdk:"request"`
	Schedule          types.String   `tfsdk:"schedule"`
	ScheduleTolerance types.Int32    `tfsdk:"schedule_tolerance"`
	Tags              []types.String `tfsdk:"tags"`
	Timezone          types.String   `tfsdk:"timezone"`
	Type              types.String   `tfsdk:"type"`
	Environments      []types.String `tfsdk:"environments"`
}

func (m *MonitorModel) hydrate(sdk *cronitor.Monitor) {
	m.Key = types.StringValue(sdk.Key)
	m.Name = types.StringValue(sdk.Name)
	for _, a := range sdk.Assertions {
		m.Assertions = append(m.Assertions, types.StringValue(a))
	}
	m.Disabled = types.BoolValue(sdk.Disabled)
	if sdk.FailureTolerance != nil {
		m.FailureTolerance = types.Int32Value(int32(*sdk.FailureTolerance))
	}
	if sdk.GraceSeconds != nil {
		m.GraceSeconds = types.Int32Value(int32(*sdk.GraceSeconds))
	}
	for _, n := range sdk.Notify {
		m.Notify = append(m.Notify, types.StringValue(n))
	}
	m.Paused = types.BoolValue(sdk.Paused)
	m.Platform = types.StringValue(sdk.Platform)
	m.RealertInterval = types.StringValue(sdk.RealertInterval)
	m.Request = &RequestModel{
		Url:             types.StringValue(sdk.Request.URL),
		Method:          types.StringValue(sdk.Request.Method),
		Body:            types.StringValue(sdk.Request.Body),
		TimeoutSeconds:  types.Int32Value(int32(sdk.Request.TimeoutSeconds)),
		FollowRedirects: types.BoolValue(sdk.Request.FollowRedirects),
		VerifySsl:       types.BoolValue(sdk.Request.VerifySsl),
		Headers:         make(map[string]basetypes.StringValue),
		Cookies:         make(map[string]basetypes.StringValue),
	}
	for _, r := range sdk.Request.Regions {
		m.Request.Regions = append(m.Request.Regions, types.StringValue(r))
	}
	for key, val := range sdk.Request.Headers {
		m.Request.Headers[key] = types.StringValue(val)
	}
	for key, val := range sdk.Request.Cookies {
		m.Request.Cookies[key] = types.StringValue(val)
	}
	m.Schedule = types.StringValue(sdk.Schedule)
	if sdk.ScheduleTolerance != nil {
		m.ScheduleTolerance = types.Int32Value(int32(*sdk.ScheduleTolerance))
	}
	for _, t := range sdk.Tags {
		m.Tags = append(m.Tags, types.StringValue(t))
	}
	if sdk.Timezone != nil {
		m.Timezone = types.StringValue(*sdk.Timezone)
	}
	m.Type = types.StringValue(sdk.Type)
	for _, e := range sdk.Environments {
		m.Environments = append(m.Environments, types.StringValue(e))
	}
}

func (d *MonitorDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (d *MonitorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Monitor data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The monitor id",
				Required:            true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "The monitor key",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The monitor name",
				Computed:            true,
			},
			"assertions": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The monitor assertions",
				Computed:            true,
			},
			"disabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the monitor is disabled",
				Computed:            true,
			},
			"failure_tolerance": schema.Int32Attribute{
				MarkdownDescription: "The number of times the monitor can fail before triggering an alert",
				Computed:            true,
			},
			"grace_seconds": schema.Int32Attribute{
				MarkdownDescription: "The number of seconds to wait after failure before triggering an alert",
				Computed:            true,
			},
			"notify": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Where the alerts are sent when a failure occurs",
				Computed:            true,
			},
			"paused": schema.BoolAttribute{
				MarkdownDescription: "Whether the monitor is paused",
				Computed:            true,
			},
			"platform": schema.StringAttribute{
				MarkdownDescription: "Used in conjuction with type to detemrine how/where the monitor runs",
				Computed:            true,
			},
			"realert_interval": schema.StringAttribute{
				MarkdownDescription: "The interval that alerts are re-sent at",
				Computed:            true,
			},
			"request": schema.SingleNestedAttribute{
				MarkdownDescription: "The details of the outgoing request",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						MarkdownDescription: "The url of the resource to monitor",
						Computed:            true,
					},
					"headers": schema.MapAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "The headers sent with the request",
						Computed:            true,
						Sensitive:           true,
					},
					"cookies": schema.MapAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "The cookies sent with the request",
						Computed:            true,
						Sensitive:           true,
					},
					"body": schema.StringAttribute{
						MarkdownDescription: "The body sent with the request",
						Computed:            true,
						Sensitive:           true,
					},
					"method": schema.StringAttribute{
						MarkdownDescription: "The method of the request",
						Computed:            true,
					},
					"timeout_seconds": schema.Int32Attribute{
						MarkdownDescription: "The numbers of seconds to wait for a response",
						Computed:            true,
					},
					"regions": schema.ListAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "The regions to run the test from",
						Computed:            true,
					},
					"follow_redirects": schema.BoolAttribute{
						MarkdownDescription: "Whether to follow redirects of the response",
						Computed:            true,
					},
					"verify_ssl": schema.BoolAttribute{
						MarkdownDescription: "Whether to verify the ssl certificate of the response",
						Computed:            true,
					},
				},
			},
			"schedule": schema.StringAttribute{
				MarkdownDescription: "The schedule the monitor runs on",
				Computed:            true,
			},
			"schedule_tolerance": schema.Int32Attribute{
				MarkdownDescription: "The number of missed scheduled executions before triggering an alert",
				Computed:            true,
			},
			"tags": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The monitor tags",
				Computed:            true,
			},
			"timezone": schema.StringAttribute{
				MarkdownDescription: "The timezone of the schedule",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the monitor",
				Computed:            true,
			},
			"environments": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The environments the monitor runs in",
				Computed:            true,
			},
		},
	}
}

func (d *MonitorDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cronitor.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *MonitorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data MonitorModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	monitor, err := d.client.Get(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to retreive monitor", err.Error())
		return
	}

	tflog.Debug(ctx, "retreived monitor details", map[string]interface{}{"monitor": monitor})

	data.hydrate(monitor)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
