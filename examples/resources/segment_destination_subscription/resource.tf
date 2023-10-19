# Configures a specific destination
resource "segment_destination_subscription" "send_to_webhook" {
  destination_id = segment_destination.webhook.id
  name = "Webhook send subscription"
  enabled = true
  action_id = "abc123"
  trigger = "type = \"track\""
  settings = jsonencode({
    "url" : "https://webhook.site/abc-123",
    "data" : {
      "@path" : "$.context.app"
    },
    "enable_batching" : false,
    "method" : "POST"
  })
}
