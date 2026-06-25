resource "snapcd_integration_namespace_supply" "alerts" {
  integration_id = data.snapcd_integration.alerts.id
  namespace_id   = snapcd_namespace.example.id
}
