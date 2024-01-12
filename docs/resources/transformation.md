---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "segment_transformation Resource - terraform-provider-segment"
subcategory: ""
description: |-
  
---

# segment_transformation (Resource)



## Example Usage

```terraform
# Configures a specific transformation
resource "segment_transformation" "example" {
  source_id      = segment_source.example.id
  name           = "My transformation name"
  enabled        = true
  if             = "event = 'Bad Event'"
  new_event_name = "Good Event"
  property_renames = [
    {
      old_name = "old-name"
      new_name = "new-name"
    }
  ]
  property_value_transformations = [
    {
      property_paths = ["properties.some-property", "context.some-property"],
      property_value = "some property value"
    },
  ]
  fql_defined_properties = []
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `enabled` (Boolean) If the Transformation is enabled.
- `fql_defined_properties` (Attributes Set) Optional array for defining new properties in FQL. Currently limited to 1 property. (see [below for nested schema](#nestedatt--fql_defined_properties))
- `if` (String) If statement (FQL) to match events.

				For standard event matchers, use the following: Track -> "event='EVENT_NAME'" Identify -> "type='identify'" Group -> "type='group'"
- `name` (String) The name of the Transformation.
- `property_renames` (Attributes Set) Optional array for renaming properties collected by your events. (see [below for nested schema](#nestedatt--property_renames))
- `property_value_transformations` (Attributes Set) Optional array for transforming properties and values collected by your events. Limited to 10 properties. (see [below for nested schema](#nestedatt--property_value_transformations))
- `source_id` (String) The Source associated with the Transformation.

### Optional

- `destination_metadata_id` (String) The optional Destination metadata associated with the Transformation.
- `new_event_name` (String) Optional new event name for renaming events. Works only for 'track' event type.

### Read-Only

- `id` (String) The id of the Transformation.

<a id="nestedatt--fql_defined_properties"></a>
### Nested Schema for `fql_defined_properties`

Required:

- `fql` (String) The FQL expression used to compute the property.
- `property_name` (String) The new property name.


<a id="nestedatt--property_renames"></a>
### Nested Schema for `property_renames`

Required:

- `new_name` (String) The new name to rename the property.
- `old_name` (String) The old name of the property.


<a id="nestedatt--property_value_transformations"></a>
### Nested Schema for `property_value_transformations`

Required:

- `property_paths` (Set of String) The property paths. The maximum number of paths is 10.
- `property_value` (String) The new value of the property paths.