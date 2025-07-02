
# Basic runner pool
resource "snapcd_runner_pool" "myrunnerpool" {
  name = "myrunnerpool"
}

# Runner pool with custom command approval threshold
resource "snapcd_runner_pool" "secure_pool" {
  name                              = "secure-pool"
  custom_command_approval_threshold = 2 # Requires 2 pre-approvals for custom commands
}
