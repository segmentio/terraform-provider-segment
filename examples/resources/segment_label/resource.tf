# Configures a label
resource "segment_label" "dev" {
  key         = "environment"
  value       = "dev"
  description = "dev environment"
}
