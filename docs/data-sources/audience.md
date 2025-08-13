# segment_audience Data Source

Fetches a Segment Audience by space and ID.

## Example Usage

```hcl
resource "segment_audience" "example" {
  space_id    = "your_space_id"
  name        = "Test Audience"
  description = "Created by Terraform"
  definition  = {
    query = "event('Shoes Bought').count() >= 1"
  }
}

data "segment_audience" "example" {
  space_id = segment_audience.example.space_id
  id       = segment_audience.example.id
}
```

## Argument Reference

- `space_id` (Required) — The Segment Space ID.
- `id` (Required) — The Audience ID.

## Attributes Reference

The following attributes are exported:

- `name` — The name of the Audience.
- `description` — The description of the Audience.
- `key` — The key of the Audience.
- `enabled` — Whether the Audience is enabled.
- `definition` — The definition of the Audience (map).
- `status` — The status of the Audience.
- `options` — Additional options for the Audience (map).

## Advanced Usage & Test Scenarios

### Error Scenario: Not Found

If you attempt to look up an Audience that does not exist, the data source will return an error. This is useful for negative testing and validation.

**Example:**

```hcl
# This will fail if the audience does not exist

data "segment_audience" "not_found" {
  space_id = "your_space_id"
  id       = "audience_does_not_exist"
}
```

**Expected error:**

```
Error: Unable to read Audience: ...
```

---

### Complex Attribute Validation

The data source supports complex/nested attributes, such as advanced queries and boolean options. You can validate these in your tests or usage.

**Example:**

```hcl
resource "segment_audience" "complex" {
  space_id    = "your_space_id"
  name        = "Complex Audience"
  description = "Testing complex attributes"
  definition  = {
    query = "event('Shoes Bought').count() >= 1 && event('Shirt Bought').count() >= 2"
  }
  options = {
    includeHistoricalData = true
  }
}

data "segment_audience" "complex" {
  space_id = segment_audience.complex.space_id
  id       = segment_audience.complex.id
}
```

You can then assert on these attributes in your tests or outputs:

- `definition.query` should match the complex query string.
- `options.includeHistoricalData` should be `true`.

---

For more test examples, see the provider's `internal/provider/audience_data_source_test.go` file. 