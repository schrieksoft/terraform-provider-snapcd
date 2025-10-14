data "snapcd_stack" "default" {
  name = "default"
}

data "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}

data "snapcd_namespace_secret" "mysecret" {
  name         = "mynamespace"
  namespace_id = data.snapcd_stack.mynamespace.id
}




// Map the contents of "my-secret" to the input parameter var.myvar
resource "snapcd_namespace_input_from_secret" "myparam" {
  input_kind   = "Param"
  name         = "myvar"
  secret_id    = data.snapcd_namespace_secret.mysecret.id
  namespace_id = data.snapcd_namespace.mynamespace.id
  usage_mode   = "UseByDefault"
}
