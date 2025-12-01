
data "snapcd_group" "mygroup" {
  name = "mygroup"
}



## Service Principal

data "snapcd_service_principal" "mysp" {
  client_id = "MyServicePrincipal"
}

resource "snapcd_group_member" "mysp_contributor" {
  group_id                = data.snapcd_group.mygroup.id
  principal_id            = data.snapcd_service_principal.mysp.id
  principal_discriminator = "ServicePrincipal"
}


## User

data "snapcd_user" "myuser" {
  user_name = "myuser@somedomain.com"
}

resource "snapcd_group_member" "myuser_contributor" {
  group_id                = data.snapcd_group.mygroup.id
  principal_id            = data.snapcd_user.myuser.id
  principal_discriminator = "User"
}


## (Nested) Group

data "snapcd_group" "mygroup" {
  user_name = "MyGroup"
}

resource "snapcd_group_member" "mygroup_contributor" {
  group_id                = data.snapcd_group.mygroup.id
  principal_id            = data.snapcd_group.mygroup.id
  principal_discriminator = "Group"
}
