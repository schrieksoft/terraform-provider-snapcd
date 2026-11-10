# Lets a Kubernetes ServiceAccount authenticate as the Service Principal with its projected
# token, so no client secret is deployed with the Runner.
resource "snapcd_service_principal" "runner" {
  client_id = "runner-prod"
}

resource "snapcd_service_principal_federated_credential" "runner_cluster" {
  service_principal_id = snapcd_service_principal.runner.id
  name                 = "prod-cluster"

  # The cluster's OIDC issuer. The Snap CD Server must be able to reach it.
  issuer   = "https://oidc.eks.eu-west-1.amazonaws.com/id/EXAMPLE"
  subject  = "system:serviceaccount:snapcd/snapcd-runner"
  audience = "snapcd"
}
