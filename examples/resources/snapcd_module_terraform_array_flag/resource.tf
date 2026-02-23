resource "snapcd_module_terraform_array_flag" "targets" {
  for_each  = toset(["aws_vpc.main", "aws_subnet.main"])
  module_id = snapcd_module.mymodule.id
  task      = "Plan"
  flag      = "Target"
  value     = each.value
}