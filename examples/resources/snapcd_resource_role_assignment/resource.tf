data "snapcd_service_principal" "mysp" {
  client_id = "MyServicePrincipal"
}

data "snapcd_stack" "default" {
  name = "default"
}


resource "snapcd_resource_role_assignment" "mysp_administrator" {
  principal_id            = snapcd_service_principal.mysp.id
  resource_id             = data.snapcd_stack.default.id
  principal_discriminator = "ServicePrincipal"
  resource_discriminator  = "Stack"
  role_name               = "Owner"
}
