resource "snapcd_namespace_pulumi_array_flag" "targets" {
  for_each     = toset(["aws_vpc.main", "aws_subnet.main"])
  namespace_id = snapcd_namespace.mynamespace.id
  task         = "Plan"
  flag         = "Target"
  value        = each.value
}