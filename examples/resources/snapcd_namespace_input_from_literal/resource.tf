data "snapcd_stack" "default" {
  name = "default"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}


resource "snapcd_namespace_input_from_literal" "myparam" {
  input_kind    = "Param"
  name          = "myvar"
  literal_value = "This will be the value of 'var.myvar'!"
  namespace_id  = snapcd_namespace.mynamespace.id
  usage_mode    = "UseByDefault"
}
