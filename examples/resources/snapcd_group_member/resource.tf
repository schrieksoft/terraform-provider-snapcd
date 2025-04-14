
resource "snapcd_group" "mygroup" {
  name = "mygroup"
}

data "snapcd_service_principal" "mysp" {
  client_id = "MyServicePrincipal"
}


resource "snapcd_group_member" "mygroup_mysp" {
  group_id                = snapcd_group.mygroup.id
  principal_id            = data.snapcd_service_principal.mysp.id
  principal_discriminator = "ServicePrincipal"
}
