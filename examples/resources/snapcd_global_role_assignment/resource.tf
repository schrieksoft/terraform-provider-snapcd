data "snapcd_service_principal" "mysp" {
  client_id = "MyServicePrincipal"
}

resource "snapcd_global_role_assignment" "mysp_administrator" {
  principal_id            = snapcd_service_principal.mysp.id
  principal_discriminator = "ServicePrincipal"
  role_name               = "Administrator"
}
