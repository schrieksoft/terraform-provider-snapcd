
data "snapcd_stack" "mystack" {
  name = "MyStack"
}

data "snapcd_service_principal" "mysp" {
  client_id = "MyServicePrincipal"
}

resource "snapcd_stack_role_assignment" "mysp_contributor" {
  stack_id                = data.snapcd_stack.mystack.id
  principal_id            = data.snapcd_service_principal.mysp.id
  principal_discriminator = "ServicePrincipal"
  role_name               = "Contributor"
}