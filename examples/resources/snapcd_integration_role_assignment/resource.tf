# Grant an integration role to a principal (User | Group | ServicePrincipal).
resource "snapcd_integration_role_assignment" "alerts_owner" {
  integration_id          = data.snapcd_integration.alerts.id
  principal_id            = var.user_id
  principal_discriminator = "User"
  role_name               = "Owner" # Owner | Contributor | Reader | IdentityAccessManager
}
