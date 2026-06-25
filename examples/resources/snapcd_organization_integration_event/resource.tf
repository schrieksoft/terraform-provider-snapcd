resource "snapcd_organization_integration_event" "notify_on_failure" {
  integration_id = data.snapcd_integration.alerts.id
  trigger        = "JobFailed"
  template       = "Job failed: {{jobName}}"
}
