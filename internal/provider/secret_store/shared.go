package secret_store

const (
	// Existing constants (maintaining backward compatibility)
	DescId                    = "Unique ID of the Secret Store."
	DescName                  = "Unique Name of the Secret Store."
	DescIsAssignedToAllScopes = "If set to true, secrets scoped to any resource in the system (any Stack, Namespace, Module or Output) can be assigned to this Secret Store"
	DescKeyVaultUrl           = "URL of the Key Vault backing the Secret Store."
	DescRegion                = "Region of the AWS Secrets Manager backing the Secret Store."

	// HCP-specific constants
	DescProjectId             = "The HCP project ID where the secret store is located."
	DescOrganizationId        = "The HCP organization ID."
	DescAppName               = "The HCP application name."
)
