
data "snapcd_group" "mygroup" {
  name = "mygroup"
}



## Service Principal

data "snapcd_service_principal" "mysp" {
  client_id = "MyServicePrincipal"
}

resource "snapcd_group_member" "mysp_contributor" {
  group_id                   = data.snapcd_group.mygroup.id
  principal_id               = data.snapcd_service_principal.mysp.id
  group_member_discriminator = "ServicePrincipal"
}


## User

data "snapcd_user" "myuser" {
  user_name = "myuser@somedomain.com"
}

resource "snapcd_group_member" "myuser_contributor" {
  group_id                   = data.snapcd_group.mygroup.id
  principal_id               = data.snapcd_user.myuser.id
  group_member_discriminator = "User"
}


## (Nested) Group

data "snapcd_group" "myothergroup" {
  name = "MyOtherGroup"
}

resource "snapcd_group_member" "mygroup_contributor" {
  group_id                   = data.snapcd_group.mygroup.id
  principal_id               = data.snapcd_group.myothergroup.id
  group_member_discriminator = "Group"
}
