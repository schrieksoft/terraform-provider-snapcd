
# Basic runner pool
resource "snapcd_runner" "myrunnerpool" {
  name = "myrunnerpool"
}

# Runner pool with custom command approval threshold
resource "snapcd_runner" "secure_pool" {
  name                              = "secure-pool"
  custom_command_approval_threshold = 2 # Requires 2 pre-approvals for custom commands
}
