package namespace_env_var

const (
	DescSharedId             = "Unique ID of the "
	DescSharedName1          = "Name of the "
	DescSharedName2          = " Must be unique in combination with `namespace_id`."
	DescSharedNamespaceId1   = "ID of the "
	DescSharedNamespaceId2   = "'s parent Namespace."
	DescSharedDefinitionName = "Name of the Definition from which to get take the input."
	DescSharedLiteralType    = "Type of literal input."
	DescSharedLiteralValue   = "Literal value of the input."
	DescSharedSecretName     = "Name of the Secret to take as input."
	DescSharedSecretType     = "Type of literal input the secret value should be formatted as."
	DescSharedSecretScope    = "Scope on which the Secret to take as input as been stored (Stack or Namespace)."
	DescSharedUsage          = "Whether the input should be used by default on all Modules, or only when explicitly selected on the Module itself."
)
