---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "segment_warehouse Data Source - terraform-provider-segment"
subcategory: ""
description: |-
  The warehouse
---

# segment_warehouse (Data Source)

The warehouse

## Example Usage

```terraform
# Gets the warehouse info
data "segment_warehouse" "my_warehouse" {
  id = "abc123"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `enabled` (Boolean) When set to true, this Warehouse receives data.
- `id` (String) The id of the Warehouse.
- `metadata` (Attributes) The metadata for the Warehouse. (see [below for nested schema](#nestedatt--metadata))
- `workspace_id` (String) The id of the Workspace that owns this Warehouse.

<a id="nestedatt--metadata"></a>
### Nested Schema for `metadata`

Read-Only:

- `description` (String) A description, in English, of this object.
- `id` (String) The id of this object.
- `logos` (Attributes) Logo information for this object. (see [below for nested schema](#nestedatt--metadata--logos))
- `name` (String) The name of this object.
- `options` (Attributes List) The Integration options for this object. (see [below for nested schema](#nestedatt--metadata--options))
- `slug` (String) A human-readable, unique identifier for object.

<a id="nestedatt--metadata--logos"></a>
### Nested Schema for `metadata.logos`

Optional:

- `alt` (String) The alternative text for this logo.
- `mark` (String) The logo mark.

Read-Only:

- `default` (String) The default URL for this logo.


<a id="nestedatt--metadata--options"></a>
### Nested Schema for `metadata.options`

Optional:

- `description` (String) An optional short text description of the field.
- `label` (String) An optional label for this field.

Read-Only:

- `name` (String) The name identifying this option in the context of a Segment Integration.
- `required` (Boolean) Whether this is a required option when setting up the Integration.
- `type` (String) Defines the type for this option in the schema.