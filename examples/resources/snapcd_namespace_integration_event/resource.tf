resource "snapcd_namespace_integration_event" "notify_on_failure" {
  integration_id = data.snapcd_integration.alerts.id
  namespace_id   = snapcd_namespace.example.id
  trigger        = "JobFailed"
  template       = "Job failed on namespace: {{jobName}}"
}
