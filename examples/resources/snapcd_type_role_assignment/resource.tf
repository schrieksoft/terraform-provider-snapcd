data "snapcd_service_principal" "mysp" {
  client_id = "MyServicePrincipal"
}

resource "snapcd_type_role_assignment" "mysp_administrator" {
  principal_id            = snapcd_service_principal.mysp.id
  principal_discriminator = "ServicePrincipal"
  resource_discriminator  = "RunnerPool"
  role_name               = "Owner"
}
