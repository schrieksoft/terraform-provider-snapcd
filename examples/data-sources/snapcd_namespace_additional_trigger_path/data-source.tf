data "snapcd_stack" "mystack" {
  name = "mystack"
}

data "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.mystack.id
}
data "snapcd_namespace_additional_trigger_path" "mytriggerpath" {
  path         = "shared/config"
  namespace_id = data.snapcd_namespace.mynamespace.id
}
