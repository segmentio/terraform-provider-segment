package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/avast/retry-go/v4"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/segmentio/public-api-sdk-go/api"
	"github.com/segmentio/terraform-provider-segment/internal/provider/docs"
	"github.com/segmentio/terraform-provider-segment/internal/provider/models"
)

var (
	_ resource.Resource              = &sourceTrackingPlanConnectionResource{}
	_ resource.ResourceWithConfigure = &sourceTrackingPlanConnectionResource{}
)

func NewSourceTrackingPlanConnectionResource() resource.Resource {
	return &sourceTrackingPlanConnectionResource{}
}

type sourceTrackingPlanConnectionResource struct {
	client      *api.APIClient
	authContext context.Context
}

func (r *sourceTrackingPlanConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_tracking_plan_connection"
}

func (r *sourceTrackingPlanConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <source_id>:<tracking_plan_id>. Got: %q", req.ID),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("source_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("tracking_plan_id"), idParts[1])...)
}

func (r *sourceTrackingPlanConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Configures a connection between a Source and a Tracking Plan. For more information, visit the [Segment docs](https://segment.com/docs/protocols/validate/connect-sources/).\n\n" +
			docs.GenerateImportDocs("<source_id>:<tracking_plan_id>", "segment_source_tracking_plan_connection"),
		Attributes: map[string]schema.Attribute{
			"source_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the Source.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tracking_plan_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the Tracking Plan.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"schema_settings": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "The schema settings associated with the Source. Upon import, this field will be empty even if the settings have already been configured due to Terraform limitations, but will be populated on the first apply. Fields not present in the config will not be managed by Terraform.",
				Attributes: map[string]schema.Attribute{
					"track": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "Track settings.",
						Attributes: map[string]schema.Attribute{
							"allow_unplanned_events": schema.BoolAttribute{
								Optional:    true,
								Description: "Enable to allow unplanned track events.",
							},
							"allow_unplanned_event_properties": schema.BoolAttribute{
								Optional:    true,
								Description: "Enable to allow unplanned track event properties.",
							},
							"allow_event_on_violations": schema.BoolAttribute{
								Optional:    true,
								Description: "Allow track event on violations.",
							},
							"allow_properties_on_violations": schema.BoolAttribute{
								Optional:    true,
								Description: "Enable to allow track properties on violations.",
							},
							"common_event_on_violations": schema.StringAttribute{
								Optional:    true,
								Description: "The common track event on violations.",
							},
						},
					},
					"identify": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "Identify settings.",
						Attributes: map[string]schema.Attribute{
							"allow_unplanned_traits": schema.BoolAttribute{
								Optional:    true,
								Description: "Enable to allow unplanned identify traits.",
							},
							"allow_traits_on_violations": schema.BoolAttribute{
								Optional:    true,
								Description: "Enable to allow identify traits on violations.",
							},
							"common_event_on_violations": schema.StringAttribute{
								Optional:    true,
								Description: "The common identify event on violations.",
							},
						},
					},
					"group": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "Group settings.",
						Attributes: map[string]schema.Attribute{
							"allow_unplanned_traits": schema.BoolAttribute{
								Optional:    true,
								Description: "Enable to allow unplanned group traits.",
							},
							"allow_traits_on_violations": schema.BoolAttribute{
								Optional:    true,
								Description: "Enable to allow group traits on violations.",
							},
							"common_event_on_violations": schema.StringAttribute{
								Optional:    true,
								Description: "The common group event on violations.",
							},
						},
					},
					"forwarding_violations_to": schema.StringAttribute{
						Optional:    true,
						Description: "Source id to forward violations to.",
					},
					"forwarding_blocked_events_to": schema.StringAttribute{
						Optional:    true,
						Description: "Source id to forward blocked events to.",
					},
				},
			},
		},
	}
}

func (r *sourceTrackingPlanConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.SourceTrackingPlanConnectionPlan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.TrackingPlanID.String() == "" || plan.SourceID.String() == "" {
		resp.Diagnostics.AddError("Unable to create connection between Source and Tracking Plan", "At least one ID is empty")

		return
	}

	err := retry.Do(
		func() error {
			_, body, err := r.client.TrackingPlansAPI.AddSourceToTrackingPlan(r.authContext, plan.TrackingPlanID.ValueString()).AddSourceToTrackingPlanV1Input(api.AddSourceToTrackingPlanV1Input{
				SourceId: plan.SourceID.ValueString(),
			}).Execute()
			if body != nil {
				defer body.Body.Close()
			}
			if err != nil {
				return errors.New(getError(err, body))
			}

			return nil
		},
		retry.Delay(DefaultRetryDelay),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Source-Tracking Plan connection",
			err.Error(),
		)

		return
	}

	var schemaSettings *api.SourceSettingsOutputV1
	apiSchemaSettings, diags := models.GetSchemaSettingsFromPlan(ctx, plan.SchemaSettings)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if apiSchemaSettings != nil {
		settingsOut, body, err := r.client.SourcesAPI.UpdateSchemaSettingsInSource(r.authContext, plan.SourceID.ValueString()).UpdateSchemaSettingsInSourceV1Input(api.UpdateSchemaSettingsInSourceV1Input{
			Track:                     apiSchemaSettings.Track,
			Identify:                  apiSchemaSettings.Identify,
			Group:                     apiSchemaSettings.Group,
			ForwardingViolationsTo:    apiSchemaSettings.ForwardingViolationsTo,
			ForwardingBlockedEventsTo: apiSchemaSettings.ForwardingBlockedEventsTo,
		}).Execute()
		if body != nil {
			defer body.Body.Close()
		}
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unable to update Source schema settings (ID: %s)", plan.SourceID.ValueString()),
				getError(err, body),
			)

			return
		}

		schemaSettings = &settingsOut.Data.Settings
	}

	var state models.SourceTrackingPlanConnectionState
	state.Fill(plan.SourceID.ValueString(), plan.TrackingPlanID.ValueString(), schemaSettings)

	// This is to satisfy terraform requirements that the returned fields must match the input ones because new settings can be generated in the response
	plannedSchemaSettings, diags := models.SchemaSettingsPlanToState(ctx, plan.SchemaSettings)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.SchemaSettings = filterOmittedSchemaSettings(plannedSchemaSettings, state.SchemaSettings)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *sourceTrackingPlanConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var previousState models.SourceTrackingPlanConnectionState

	diags := req.State.Get(ctx, &previousState)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, body, err := r.client.SourcesAPI.GetSource(r.authContext, previousState.SourceID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		if body.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read Source (ID: %s)", previousState.SourceID.ValueString()),
			getError(err, body),
		)

		return
	}

	if !out.Data.TrackingPlanId.IsSet() || out.Data.TrackingPlanId.Get() == nil || *out.Data.TrackingPlanId.Get() != previousState.TrackingPlanID.ValueString() {
		diags = resp.State.Set(ctx, &models.SourceTrackingPlanConnectionState{
			SourceID:       previousState.SourceID,
			TrackingPlanID: types.StringValue(""),
		})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		return
	}

	var schemaSettings *api.SourceSettingsOutputV1
	if previousState.SchemaSettings != nil {
		settingsOut, body, err := r.client.SourcesAPI.ListSchemaSettingsInSource(r.authContext, previousState.SourceID.ValueString()).Execute()
		if body != nil {
			defer body.Body.Close()
		}
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unable to read Source schema settings (ID: %s)", previousState.SourceID.ValueString()),
				getError(err, body),
			)

			return
		}

		schemaSettings = &settingsOut.Data.Settings
	}

	var state models.SourceTrackingPlanConnectionState
	state.Fill(previousState.SourceID.ValueString(), previousState.TrackingPlanID.ValueString(), schemaSettings)

	// This is to satisfy terraform requirements that the returned fields must match the input ones because new settings can be generated in the response
	previousState.SchemaSettings = filterOmittedSchemaSettings(previousState.SchemaSettings, previousState.SchemaSettings)

	diags = resp.State.Set(ctx, &previousState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *sourceTrackingPlanConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.SourceTrackingPlanConnectionPlan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, body, err := r.client.SourcesAPI.GetSource(r.authContext, plan.SourceID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read Source (ID: %s)", plan.SourceID.ValueString()),
			getError(err, body),
		)

		return
	}

	if !out.Data.TrackingPlanId.IsSet() || out.Data.TrackingPlanId.Get() == nil || *out.Data.TrackingPlanId.Get() != plan.TrackingPlanID.ValueString() {
		resp.Diagnostics.AddError(
			"Source is not connected to specified Tracking Plan",
			fmt.Sprintf("Source ID: '%s', Tracking Plan ID: '%s'", plan.SourceID.ValueString(), plan.TrackingPlanID.ValueString()),
		)

		return
	}

	var schemaSettings *api.SourceSettingsOutputV1
	apiSchemaSettings, diags := models.GetSchemaSettingsFromPlan(ctx, plan.SchemaSettings)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if apiSchemaSettings != nil {
		settingsOut, body, err := r.client.SourcesAPI.UpdateSchemaSettingsInSource(r.authContext, plan.SourceID.ValueString()).UpdateSchemaSettingsInSourceV1Input(api.UpdateSchemaSettingsInSourceV1Input{
			Track:                     apiSchemaSettings.Track,
			Identify:                  apiSchemaSettings.Identify,
			Group:                     apiSchemaSettings.Group,
			ForwardingViolationsTo:    apiSchemaSettings.ForwardingViolationsTo,
			ForwardingBlockedEventsTo: apiSchemaSettings.ForwardingBlockedEventsTo,
		}).Execute()
		if body != nil {
			defer body.Body.Close()
		}
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unable to update Source schema settings (ID: %s)", plan.SourceID.ValueString()),
				getError(err, body),
			)

			return
		}

		schemaSettings = &settingsOut.Data.Settings
	}

	var state models.SourceTrackingPlanConnectionState
	state.Fill(plan.SourceID.ValueString(), plan.TrackingPlanID.ValueString(), schemaSettings)

	// This is to satisfy terraform requirements that the returned fields must match the input ones because new settings can be generated in the response
	plannedSchemaSettings, diags := models.SchemaSettingsPlanToState(ctx, plan.SchemaSettings)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.SchemaSettings = filterOmittedSchemaSettings(plannedSchemaSettings, state.SchemaSettings)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *sourceTrackingPlanConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var config models.SourceTrackingPlanConnectionPlan
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If either of these are empty, it means the connection no longer exists
	if config.TrackingPlanID.ValueString() != "" && config.SourceID.ValueString() != "" {
		_, body, err := r.client.TrackingPlansAPI.RemoveSourceFromTrackingPlan(r.authContext, config.TrackingPlanID.ValueString()).SourceId(config.SourceID.ValueString()).Execute()
		if body != nil {
			defer body.Body.Close()
		}
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to remove Source-Tracking Plan connection",
				getError(err, body),
			)

			return
		}
	}
}

func (r *sourceTrackingPlanConnectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*ClientInfo)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected ClientInfo, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = config.client
	r.authContext = config.authContext
}

// Filters out fields that were omitted from the plan to ensure consistent terraform state.
func filterOmittedSchemaSettings(plannedState *models.SchemaSettingsState, returnedState *models.SchemaSettingsState) *models.SchemaSettingsState {
	if plannedState == nil || returnedState == nil {
		return nil
	}

	out := models.SchemaSettingsState{}

	if plannedState.Track != nil {
		out.Track = &models.TrackSettings{}

		if !plannedState.Track.AllowEventOnViolations.IsNull() && !plannedState.Track.AllowEventOnViolations.IsUnknown() {
			out.Track.AllowEventOnViolations = returnedState.Track.AllowEventOnViolations
		}
		if !plannedState.Track.AllowPropertiesOnViolations.IsNull() && !plannedState.Track.AllowPropertiesOnViolations.IsUnknown() {
			out.Track.AllowPropertiesOnViolations = returnedState.Track.AllowPropertiesOnViolations
		}
		if !plannedState.Track.AllowUnplannedEvents.IsNull() && !plannedState.Track.AllowUnplannedEvents.IsUnknown() {
			out.Track.AllowUnplannedEvents = returnedState.Track.AllowUnplannedEvents
		}
		if !plannedState.Track.AllowUnplannedEventProperties.IsNull() && !plannedState.Track.AllowUnplannedEventProperties.IsUnknown() {
			out.Track.AllowUnplannedEventProperties = returnedState.Track.AllowUnplannedEventProperties
		}
		if !plannedState.Track.CommonEventOnViolations.IsNull() && !plannedState.Track.CommonEventOnViolations.IsUnknown() {
			out.Track.CommonEventOnViolations = returnedState.Track.CommonEventOnViolations
		}
	}

	if plannedState.Identify != nil {
		out.Identify = &models.IdentifySettings{}

		if !plannedState.Identify.AllowTraitsOnViolations.IsNull() && !plannedState.Identify.AllowTraitsOnViolations.IsUnknown() {
			out.Identify.AllowTraitsOnViolations = returnedState.Identify.AllowTraitsOnViolations
		}
		if !plannedState.Identify.AllowUnplannedTraits.IsNull() && !plannedState.Identify.AllowUnplannedTraits.IsUnknown() {
			out.Identify.AllowUnplannedTraits = returnedState.Identify.AllowUnplannedTraits
		}
		if !plannedState.Identify.CommonEventOnViolations.IsNull() && !plannedState.Identify.CommonEventOnViolations.IsUnknown() {
			out.Identify.CommonEventOnViolations = returnedState.Identify.CommonEventOnViolations
		}
	}

	if plannedState.Group != nil {
		out.Group = &models.GroupSettings{}

		if !plannedState.Group.AllowTraitsOnViolations.IsNull() && !plannedState.Group.AllowTraitsOnViolations.IsUnknown() {
			out.Group.AllowTraitsOnViolations = returnedState.Group.AllowTraitsOnViolations
		}
		if !plannedState.Group.AllowUnplannedTraits.IsNull() && !plannedState.Group.AllowUnplannedTraits.IsUnknown() {
			out.Group.AllowUnplannedTraits = returnedState.Group.AllowUnplannedTraits
		}
		if !plannedState.Group.CommonEventOnViolations.IsNull() && !plannedState.Group.CommonEventOnViolations.IsUnknown() {
			out.Group.CommonEventOnViolations = returnedState.Group.CommonEventOnViolations
		}
	}

	if !plannedState.ForwardingBlockedEventsTo.IsNull() && !plannedState.ForwardingBlockedEventsTo.IsUnknown() {
		out.ForwardingBlockedEventsTo = returnedState.ForwardingBlockedEventsTo
	}

	if !plannedState.ForwardingViolationsTo.IsNull() && !plannedState.ForwardingViolationsTo.IsUnknown() {
		out.ForwardingViolationsTo = returnedState.ForwardingViolationsTo
	}

	return &out
}
