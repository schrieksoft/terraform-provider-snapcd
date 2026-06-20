# Supply the integration to a specific module (scope = Stack | Namespace | Module).
resource "snapcd_integration_supply" "alerts_to_module" {
  integration_id = data.snapcd_integration.alerts.id
  scope          = "Module"
  scope_id       = snapcd_module.example.id
}
