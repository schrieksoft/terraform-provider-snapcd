package integration_events

import (
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
)

var IntegrationDataSourceConfig = `
data "snapcd_integration" "this" {
  name = "debug-slack"
}`

var OrganizationIntegrationEventCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_organization_integration_event" "this" {
  integration_id = data.snapcd_integration.this.id
  trigger        = "MissionFaulted"
  template       = "Mission faulted: {{jobName}} %s"
}`)

var OrganizationIntegrationEventUpdateConfig = providerconfig.AppendRandomString(`
resource "snapcd_organization_integration_event" "this" {
  integration_id = data.snapcd_integration.this.id
  trigger        = "MissionFaulted"
  template       = "Mission faulted (updated): {{jobName}} %s"
}`)

var StackIntegrationEventCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_stack_integration_event" "this" {
  integration_id = data.snapcd_integration.this.id
  stack_id       = snapcd_stack.this.id
  trigger        = "JobFailed"
  template       = "Job failed on stack: {{jobName}} %s"
}`)

var StackIntegrationEventUpdateConfig = providerconfig.AppendRandomString(`
resource "snapcd_stack_integration_event" "this" {
  integration_id = data.snapcd_integration.this.id
  stack_id       = snapcd_stack.this.id
  trigger        = "JobFailed"
  template       = "Job failed on stack (updated): {{jobName}} %s"
}`)

var NamespaceIntegrationEventCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_namespace_integration_event" "this" {
  integration_id = data.snapcd_integration.this.id
  namespace_id   = snapcd_namespace.this.id
  trigger        = "JobFailed"
  template       = "Job failed on namespace: {{jobName}} %s"
}`)

var NamespaceIntegrationEventUpdateConfig = providerconfig.AppendRandomString(`
resource "snapcd_namespace_integration_event" "this" {
  integration_id = data.snapcd_integration.this.id
  namespace_id   = snapcd_namespace.this.id
  trigger        = "JobFailed"
  template       = "Job failed on namespace (updated): {{jobName}} %s"
}`)

var ModuleIntegrationEventCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_module_integration_event" "this" {
  integration_id = data.snapcd_integration.this.id
  module_id      = snapcd_module.this.id
  trigger        = "JobFailed"
  template       = "Job failed on module: {{jobName}} %s"
}`)

var ModuleIntegrationEventUpdateConfig = providerconfig.AppendRandomString(`
resource "snapcd_module_integration_event" "this" {
  integration_id = data.snapcd_integration.this.id
  module_id      = snapcd_module.this.id
  trigger        = "JobFailed"
  template       = "Job failed on module (updated): {{jobName}} %s"
}`)
