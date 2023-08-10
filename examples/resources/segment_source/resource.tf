resource "segment_source" "test" {
  id           = "source-id"
  slug         = "source-slug"
  workspace_id = "w-123"
  enabled      = true
  metadata = {
    id = "m-123"
  }
}