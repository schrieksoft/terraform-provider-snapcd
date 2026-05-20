data "snapcd_agent" "myagent" {
  name = "myagent"
}


## Service Principal

data "snapcd_service_principal" "mysp" {
  client_id = "MyServicePrincipal"
}

resource "snapcd_agent_role_assignment" "mysp_contributor" {
  agent_id                = data.snapcd_agent.myagent.id
  principal_id            = data.snapcd_service_principal.mysp.id
  principal_discriminator = "ServicePrincipal"
  role_name               = "Contributor"
}


## User

data "snapcd_user" "myuser" {
  user_name = "myuser@somedomain.com"
}

resource "snapcd_agent_role_assignment" "myuser_contributor" {
  agent_id                = data.snapcd_agent.myagent.id
  principal_id            = data.snapcd_user.myuser.id
  principal_discriminator = "User"
  role_name               = "Contributor"
}


## Group

data "snapcd_group" "mygroup" {
  name = "MyGroup"
}

resource "snapcd_agent_role_assignment" "mygroup_contributor" {
  agent_id                = data.snapcd_agent.myagent.id
  principal_id            = data.snapcd_group.mygroup.id
  principal_discriminator = "Group"
  role_name               = "Contributor"
}
