# Configures a specific invite/user
resource "segment_user" "example" {
  email = "example@segment.com"
  permissions = [
    {
      role_id = "abc123"
      resources = [
        {
          id   = "abc123"
          type = "WORKSPACE"
          labels = [
            {
              key   = "env"
              value = "stage"
            }
          ]
        }
      ]
    }
  ]
}
