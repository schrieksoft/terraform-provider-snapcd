data "snapcd_stack" "mystack" {
  name = "mystack"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.mystack.id
}

data "snapcd_runner" "myrunner" {
  name = "myrunner"
}

resource "snapcd_module" "mymodule" {
  name            = "mymodule"
  namespace_id    = snapcd_namespace.mynamespace.id
  source_revision = "main"
  source_url      = "https://github.com/schrieksoft/snapcd-samples.git"
  runner_id       = data.snapcd_runner.myrunner.id
}

resource "snapcd_module_terraform_remote_policy" "mypolicy" {
  name      = "mypolicy"
  module_id = snapcd_module.mymodule.id

  repo_url = "https://github.com/myorg/policy-repo.git"
  revision = "v1.0.0"
  path     = "policies/terraform"
}
