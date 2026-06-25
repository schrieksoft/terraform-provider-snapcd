resource "snapcd_stack_integration_event" "notify_on_failure" {
  integration_id = data.snapcd_integration.alerts.id
  stack_id       = snapcd_stack.example.id
  trigger        = "JobFailed"
  template       = "Job failed on stack: {{jobName}}"
}
