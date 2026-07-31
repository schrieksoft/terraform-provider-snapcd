data "snapcd_stack" "mystack" {
  name = "mystack"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.mystack.id
}

// A Namespace Additional Trigger Path joins the trigger watch set of every module in the
// namespace that has the trigger path filter enabled — useful for a directory that all modules
// in the namespace depend on, such as shared configuration.

resource "snapcd_namespace_additional_trigger_path" "shared_config" {
  path         = "shared/config"
  namespace_id = snapcd_namespace.mynamespace.id
}
