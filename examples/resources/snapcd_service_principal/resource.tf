# Created without a client secret. Give it a federated credential so the workload can
# authenticate with a token its own platform issues.
resource "snapcd_service_principal" "runner" {
  client_id = "runner-prod"
}

resource "snapcd_service_principal_federated_credential" "runner_cluster" {
  service_principal_id = snapcd_service_principal.runner.id
  name                 = "prod-cluster"

  issuer   = "https://oidc.eks.eu-west-1.amazonaws.com/id/EXAMPLE"
  subject  = "system:serviceaccount:snapcd/snapcd-runner"
  audience = "snapcd"
}
