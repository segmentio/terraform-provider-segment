---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "segment_destination_filter Resource - terraform-provider-segment"
subcategory: ""
description: |-
  Destination filters let you send or block events, properties, and user traits from reaching a destination. Enabled filters are applied on every matching event in transit to this destination.
      To import a Destination filter into Terraform, use the following format: 'destination-id:filter-id'
---

# segment_destination_filter (Resource)

Destination filters let you send or block events, properties, and user traits from reaching a destination. Enabled filters are applied on every matching event in transit to this destination.
		
		To import a Destination filter into Terraform, use the following format: 'destination-id:filter-id'

## Example Usage

```terraform
# Configures a specific destination filter
resource "segment_destination_filter" "sample" {
  if             = "type = \"identify\""
  destination_id = "abc123"
  source_id      = "xyz321"
  title          = "Identify event sampling filter"
  enabled        = true
  description    = "Samples identify events at 5%"
  actions = [
    {
      type    = "SAMPLE"
      percent = 0.05
    },
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `actions` (Attributes Set) Actions for the Destination filter. (see [below for nested schema](#nestedatt--actions))
- `destination_id` (String) The id of the Destination associated with this filter.
- `enabled` (Boolean) When set to true, the Destination filter is active.
- `if` (String) The filter's condition.
- `source_id` (String) The id of the Source associated with this filter.
- `title` (String) The title of the filter.

### Optional

- `description` (String) The description of the filter.

### Read-Only

- `id` (String) The unique id of this filter.

<a id="nestedatt--actions"></a>
### Nested Schema for `actions`

Required:

- `type` (String) The kind of Transformation to apply to any matched properties.

								Enum: "ALLOW_PROPERTIES" "DROP" "DROP_PROPERTIES" "SAMPLE"

Optional:

- `fields` (String) A dictionary of paths to object keys that this filter applies to. The literal string '' represents the top level of the object.
- `path` (String) The JSON path to a property within a payload object from which Segment generates a deterministic sampling rate.
- `percent` (Number) A decimal between 0 and 1 used for 'sample' type events and influences the likelihood of sampling to occur.