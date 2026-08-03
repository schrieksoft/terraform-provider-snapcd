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

resource "snapcd_module_terraform_inline_policy" "mypolicy" {
  name      = "mypolicy"
  module_id = snapcd_module.mymodule.id

  policy_content = <<-EOF
    package terraform.security

    import rego.v1

    deny contains msg if {
      some r in input.resource_changes

      r.type == "aws_s3_bucket_public_access_block"
      r.change.after != null

      not r.change.after.block_public_acls

      msg := sprintf("%s allows public ACLs", [r.address])
    }
  EOF
}
