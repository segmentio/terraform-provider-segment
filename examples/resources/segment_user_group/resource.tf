# Configures a specific user group
data "segment_workspace" "example" {}

resource "segment_label" "env_dev" {
  key   = "env"
  value = "dev"
}

resource "segment_user_group" "example" {
  name = "My user group"
  permissions = [
    {
      role_id = "abc123"
      resources = [
        {
          id   = data.segment_workspace.example.id
          type = "WORKSPACE"
          labels = [
            {
              key   = segment_label.env_dev.key
              value = segment_label.env_dev.value
            }
          ]
        }
      ]
    }
  ]
  members = ["example@segment.com"]
}
