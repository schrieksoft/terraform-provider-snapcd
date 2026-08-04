data "snapcd_stack" "mystack" {
  name = "mystack"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.mystack.id
}

resource "snapcd_namespace_pulumi_inline_policy" "mypolicy" {
  name         = "mypolicy"
  namespace_id = snapcd_namespace.mynamespace.id

  runtime = "Python"

  policy_content = <<-EOF
    from pulumi_policy import EnforcementLevel, PolicyPack, ResourceValidationPolicy

    def no_public_buckets(args, report_violation):
        if args.resource_type == "aws:s3/bucket:Bucket":
            if args.props.get("acl") == "public-read":
                report_violation(f"{args.name} must not be public")

    PolicyPack(
        name="no-public-buckets",
        enforcement_level=EnforcementLevel.MANDATORY,
        policies=[
            ResourceValidationPolicy(
                name="no-public-buckets",
                description="S3 buckets must not be public.",
                validate=no_public_buckets,
            ),
        ],
    )
  EOF
}
