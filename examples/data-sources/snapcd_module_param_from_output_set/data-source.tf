data "snapcd_stack" "default" {
  name = "default"
}

data "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}
data "snapcd_module" "mymodule" {
  name         = "mymodule"
  namespace_id = data.snapcd_namespace.mynamespace.id
}
data "snapcd_module_env_var_from_output_set" "myvar" {
  name      = "from_anothermodule"
  module_id = data.snapcd_module.mymodule.id
}
