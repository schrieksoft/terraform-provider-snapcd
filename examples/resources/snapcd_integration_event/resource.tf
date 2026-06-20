# Demand: fire the integration on a trigger for a scope (Organization | Stack | Namespace | Module).
resource "snapcd_integration_event" "notify_on_failure" {
  scope          = "Module"
  scope_id       = snapcd_module.example.id
  integration_id = data.snapcd_integration.alerts.id
  trigger        = "JobFailed" # see the IntegrationTrigger catalog
  template       = "❌ {{jobType}} failed on *{{moduleName}}* (job {{jobId}})."
  is_disabled    = false
}
