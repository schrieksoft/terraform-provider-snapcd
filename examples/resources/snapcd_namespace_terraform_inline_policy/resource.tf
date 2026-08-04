data "snapcd_stack" "mystack" {
  name = "mystack"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.mystack.id
}

resource "snapcd_namespace_terraform_inline_policy" "mypolicy" {
  name         = "mypolicy"
  namespace_id = snapcd_namespace.mynamespace.id

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
