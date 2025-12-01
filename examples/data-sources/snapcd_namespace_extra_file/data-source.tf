data "snapcd_stack" "mystack" {
  name = "mystack"
}

data "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.mystack.id
}
data "snapcd_namespace_extra_file" "myextrafile" {
  file_name    = "myextrafile.tf"
  namespace_id = data.snapcd_namespace.mynamespace.id
}
