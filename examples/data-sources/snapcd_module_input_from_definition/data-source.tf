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
data "snapcd_module_input_from_definition" "myvar" {
  input_kind = "Param"
  name       = "myvar"
  module_id  = data.snapcd_module.mymodule.id
}
