resource "snapcd_integration_module_supply" "alerts" {
  integration_id = data.snapcd_integration.alerts.id
  module_id      = snapcd_module.example.id
}
