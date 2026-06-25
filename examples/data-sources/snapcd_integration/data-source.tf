# Integrations are created/managed in the SnapCd UI (the bot token is write-only).
# Reference one by name to wire up assignments, role assignments, and events.
data "snapcd_integration" "alerts" {
  name = "alerts"
}
