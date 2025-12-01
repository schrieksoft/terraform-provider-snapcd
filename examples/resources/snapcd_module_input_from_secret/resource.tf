data "snapcd_stack" "mystack" {
  name = "mystack"
}

data "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.mystack.id
}


data "snapcd_module" "mymodule" {
  name         = "mymodule"
  namespace_id = data.snapcd_namespace.mynamespace.id
}

data "snapcd_module_secret" "mysecret" {
  name      = "mymodule"
  module_id = data.snapcd_module.mymodule.id
}

resource "snapcd_module_input_from_secret" "myparam" {
  input_kind = "Param"
  name       = "myvar"
  secret_id  = data.snapcd_module_secret.mysecret.id
  module_id  = data.snapcd_module.mymodule.id
}
