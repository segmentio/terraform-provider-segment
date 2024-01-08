package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type TransformationPlan struct {
	ID                           types.String `tfsdk:"id"`
	SourceID                     types.String `tfsdk:"source_id"`
	DestinationMetadataID        types.String `tfsdk:"destination_metadata_id"`
	Name                         types.String `tfsdk:"name"`
	Enabled                      types.Bool   `tfsdk:"enabled"`
	If                           types.String `tfsdk:"if"`
	NewEventName                 types.String `tfsdk:"new_event_name"`
	PropertyRenames              types.Set    `tfsdk:"property_renames"`
	PropertyValueTransformations types.Set    `tfsdk:"property_value_transformations"`
	FQLDefinedProperties         types.Set    `tfsdk:"fql_defined_properties"`
}

type TransformationState struct {
	ID                           types.String             `tfsdk:"id"`
	SourceID                     types.String             `tfsdk:"source_id"`
	DestinationMetadataID        types.String             `tfsdk:"destination_metadata_id"`
	Name                         types.String             `tfsdk:"name"`
	Enabled                      types.Bool               `tfsdk:"enabled"`
	If                           types.String             `tfsdk:"if"`
	NewEventName                 types.String             `tfsdk:"new_event_name"`
	PropertyRenames              []PropertyRename         `tfsdk:"property_renames"`
	PropertyValueTransformations []PropertyValueTransform `tfsdk:"property_value_transformations"`
	FQLDefinedProperties         []FQLDefinedProperty     `tfsdk:"fql_defined_properties"`
}

type PropertyRename struct {
	OldName types.String `tfsdk:"old_name"`
	NewName types.String `tfsdk:"new_name"`
}

type PropertyValueTransform struct {
	PropertyPaths []types.String `tfsdk:"property_paths"`
	PropertyValue types.String   `tfsdk:"property_value"`
}

type PropertyValueTransformPlan struct {
	PropertyPaths types.Set    `tfsdk:"property_paths"`
	PropertyValue types.String `tfsdk:"property_value"`
}

type FQLDefinedProperty struct {
	FQL          types.String `tfsdk:"fql"`
	PropertyName types.String `tfsdk:"property_name"`
}

func (t *TransformationState) Fill(transformation api.TransformationV1) {
	t.ID = types.StringValue(transformation.Id)
	t.SourceID = types.StringValue(transformation.SourceId)
	t.DestinationMetadataID = types.StringPointerValue(transformation.DestinationMetadataId)
	t.Name = types.StringValue(transformation.Name)
	t.Enabled = types.BoolValue(transformation.Enabled)
	t.If = types.StringValue(transformation.If)
	t.NewEventName = types.StringPointerValue(transformation.NewEventName)

	// Fill PropertyRenames
	t.PropertyRenames = make([]PropertyRename, len(transformation.PropertyRenames))
	for i, pr := range transformation.PropertyRenames {
		t.PropertyRenames[i] = PropertyRename{
			OldName: types.StringValue(pr.OldName),
			NewName: types.StringValue(pr.NewName),
		}
	}

	// Fill PropertyValueTransformations
	t.PropertyValueTransformations = make([]PropertyValueTransform, len(transformation.PropertyValueTransformations))
	for i, pvt := range transformation.PropertyValueTransformations {
		var paths []types.String
		for _, path := range pvt.PropertyPaths {
			paths = append(paths, types.StringValue(path))
		}

		t.PropertyValueTransformations[i] = PropertyValueTransform{
			PropertyPaths: paths,
			PropertyValue: types.StringValue(pvt.PropertyValue),
		}
	}

	// Fill FQLDefinedProperties
	t.FQLDefinedProperties = make([]FQLDefinedProperty, len(transformation.FqlDefinedProperties))
	for i, fdp := range transformation.FqlDefinedProperties {
		t.FQLDefinedProperties[i] = FQLDefinedProperty{
			FQL:          types.StringValue(fdp.Fql),
			PropertyName: types.StringValue(fdp.PropertyName),
		}
	}
}

func PropertyRenamesPlanToAPIValue(ctx context.Context, renames types.Set) ([]api.PropertyRenameV1, diag.Diagnostics) {
	apiRenames := []api.PropertyRenameV1{}

	if !renames.IsNull() && !renames.IsUnknown() {
		stateRenames := []PropertyRename{}
		diags := renames.ElementsAs(ctx, &stateRenames, false)
		if diags.HasError() {
			return apiRenames, diags
		}
		for _, rename := range stateRenames {
			apiRenames = append(apiRenames, api.PropertyRenameV1{
				OldName: rename.OldName.ValueString(),
				NewName: rename.NewName.ValueString(),
			})
		}
	}

	return apiRenames, diag.Diagnostics{}
}

func PropertyValueTransformationsPlanToAPIValue(ctx context.Context, transforms types.Set) ([]api.PropertyValueTransformationV1, diag.Diagnostics) {
	apiTransforms := []api.PropertyValueTransformationV1{}

	if !transforms.IsNull() && !transforms.IsUnknown() {
		stateTransforms := []PropertyValueTransformPlan{}
		diags := transforms.ElementsAs(ctx, &stateTransforms, false)
		if diags.HasError() {
			return apiTransforms, diags
		}
		for _, transform := range stateTransforms {
			paths := []string{}
			diags := transform.PropertyPaths.ElementsAs(ctx, &paths, false)
			if diags.HasError() {
				return apiTransforms, diags
			}

			apiTransforms = append(apiTransforms, api.PropertyValueTransformationV1{
				PropertyPaths: paths,
				PropertyValue: transform.PropertyValue.ValueString(),
			})
		}
	}

	return apiTransforms, diag.Diagnostics{}
}

func FQLDefinedPropertiesPlanToAPIValue(ctx context.Context, properties types.Set) ([]api.FQLDefinedPropertyV1, diag.Diagnostics) {
	apiPreperties := []api.FQLDefinedPropertyV1{}

	if !properties.IsNull() && !properties.IsUnknown() {
		stateProperties := []FQLDefinedProperty{}
		diags := properties.ElementsAs(ctx, &stateProperties, false)
		if diags.HasError() {
			return apiPreperties, diags
		}
		for _, property := range stateProperties {
			apiPreperties = append(apiPreperties, api.FQLDefinedPropertyV1{
				Fql:          property.FQL.ValueString(),
				PropertyName: property.PropertyName.ValueString(),
			})
		}
	}

	return apiPreperties, diag.Diagnostics{}
}
