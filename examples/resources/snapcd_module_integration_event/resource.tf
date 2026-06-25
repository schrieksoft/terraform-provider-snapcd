resource "snapcd_module_integration_event" "notify_on_failure" {
  integration_id = data.snapcd_integration.alerts.id
  module_id      = snapcd_module.example.id
  trigger        = "JobFailed"
  template       = "Job failed on module: {{jobName}}"
}
