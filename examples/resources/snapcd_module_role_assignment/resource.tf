
data "snapcd_stack" "mystack" {
  name = "MyStack"
}

data "snapcd_namespace" "mynamespace" {
  stack_id = data.snapcd_stack.mystack.id
  name     = "MyNamespace"
}

data "snapcd_module" "mymodule" {
  namespace_id = data.snapcd_namespace.mynamespace.id
  name         = "MyModule"
}

data "snapcd_service_principal" "mysp" {
  client_id = "MyServicePrincipal"
}

resource "snapcd_module_role_assignment" "mysp_owner" {
  module_id               = data.snapcd_module.mymodule.id
  principal_id            = data.snapcd_service_principal.mysp.id
  principal_discriminator = "ServicePrincipal"
  role_name               = "Owner"
}