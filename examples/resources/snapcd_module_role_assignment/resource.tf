
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



## Service Principal

data "snapcd_service_principal" "mysp" {
  client_id = "MyServicePrincipal"
}

resource "snapcd_module_role_assignment" "mysp_contributor" {
  module_id               = data.snapcd_module.mymodule.id
  principal_id            = data.snapcd_service_principal.mysp.id
  principal_discriminator = "ServicePrincipal"
  role_name               = "Contributor"
}


## User

data "snapcd_user" "myuser" {
  user_name = "myuser@somedomain.com"
}

resource "snapcd_module_role_assignment" "myuser_contributor" {
  module_id               = data.snapcd_module.mymodule.id
  principal_id            = data.snapcd_user.myuser.id
  principal_discriminator = "User"
  role_name               = "Contributor"
}


## Group

data "snapcd_group" "mygroup" {
  user_name = "MyGroup"
}

resource "snapcd_module_role_assignment" "mygroup_contributor" {
  module_id               = data.snapcd_module.mymodule.id
  principal_id            = data.snapcd_group.mygroup.id
  principal_discriminator = "Group"
  role_name               = "Contributor"
}
