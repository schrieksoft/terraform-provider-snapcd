package secret_store

const (
	DescId                    = "Unique ID of the Secret Store."
	DescName                  = "Unique Name of the Secret Store."
	DescIsAssignedToAllScopes = "If set to true, secrets scoped to any resource in the system (any Stack, Namespace, Module or Output) can be assigned to this Secret Store"
	DescKeyVaultUrl           = "URL of the Key Vault backing the Secret Store."
)
