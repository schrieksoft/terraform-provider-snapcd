// Fetch the Service Principal the Agent authenticates as.
data "snapcd_service_principal" "my_service_principal" {
  client_id = "MyAgentServicePrincipal"
}

resource "snapcd_agent" "my_agent" {
  name                     = "MyAgent"
  service_principal_id     = data.snapcd_service_principal.my_service_principal.id
  allow_multiple_instances = false
}
