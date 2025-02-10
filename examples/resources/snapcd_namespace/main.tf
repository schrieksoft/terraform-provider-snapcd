# Copyright (c) HashiCorp, Inc.

data "snapcd_stack" "default" {
  name = "default"
}

resource "snapcd_stack" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}
