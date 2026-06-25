resource "snapcd_integration_stack_supply" "alerts" {
  integration_id = data.snapcd_integration.alerts.id
  stack_id       = snapcd_stack.example.id
}
