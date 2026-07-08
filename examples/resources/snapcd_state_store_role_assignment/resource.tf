resource "snapcd_state_store" "my_state_store" {
  name = "MyStateStore"
}


## Service Principal

data "snapcd_service_principal" "mysp" {
  client_id = "MyServicePrincipal"
}

resource "snapcd_state_store_role_assignment" "mysp_contributor" {
  state_store_id          = snapcd_state_store.my_state_store.id
  principal_id            = data.snapcd_service_principal.mysp.id
  principal_discriminator = "ServicePrincipal"
  role_name               = "Contributor"
}


## User

data "snapcd_user" "myuser" {
  user_name = "myuser@somedomain.com"
}

resource "snapcd_state_store_role_assignment" "myuser_contributor" {
  state_store_id          = snapcd_state_store.my_state_store.id
  principal_id            = data.snapcd_user.myuser.id
  principal_discriminator = "User"
  role_name               = "Contributor"
}


## Group

data "snapcd_group" "mygroup" {
  name = "MyGroup"
}

resource "snapcd_state_store_role_assignment" "mygroup_contributor" {
  state_store_id          = snapcd_state_store.my_state_store.id
  principal_id            = data.snapcd_group.mygroup.id
  principal_discriminator = "Group"
  role_name               = "Contributor"
}
