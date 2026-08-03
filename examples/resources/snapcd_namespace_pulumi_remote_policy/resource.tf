data "snapcd_stack" "mystack" {
  name = "mystack"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.mystack.id
}

resource "snapcd_namespace_pulumi_remote_policy" "mypolicy" {
  name         = "mypolicy"
  namespace_id = snapcd_namespace.mynamespace.id

  repo_url = "https://github.com/myorg/policy-repo.git"
  revision = "v1.0.0"
  path     = "packs/aws-guardrails"
}
