package engine_flags

// Pulumi flags

var pulumiCommandTaskValues = []string{
	"Init",
	"Plan",
	"Apply",
	"Destroy",
	"Output",
}

var pulumiFlagValues = []string{
	"CloudUrl",
	"LoginLocal",
	"LoginCloud",
	"DefaultOrg",
	"Insecure",
	"StackName",
	"SecretsProvider",
	"CreateStack",
	"OidcExpiration",
	"OidcOrg",
	"OidcTeam",
	"OidcToken",
	"OidcUser",
	"ConfigFile",
	"Debug",
	"Diff",
	"ExpectNoChanges",
	"Json",
	"Message",
	"Parallel",
	"Refresh",
	"RunProgram",
	"ShowConfig",
	"ShowFullOutput",
	"ShowReads",
	"ShowReplacementSteps",
	"ShowSames",
	"ShowSecrets",
	"SuppressOutputs",
	"SuppressProgress",
	"TargetDependents",
	"ExcludeDependents",
	"Neo",
	"ImportFile",
	"ContinueOnError",
	"SkipPreview",
	"Strict",
	"ExcludeProtected",
	"Remove",
	"Shell",
	"Color",
	"Verbose",
	"Emoji",
}

var pulumiArrayFlagValues = []string{
	"PolicyPack",
	"PolicyPackConfig",
	"AttachDebugger",
	"Target",
	"Replace",
	"Exclude",
	"TargetReplace",
	"Config",
}

// Terraform flags

var terraformCommandTaskValues = []string{
	"Init",
	"Plan",
	"Apply",
	"Destroy",
	"Output",
}

var terraformFlagValues = []string{
	"ForceCopy",
	"FromModule",
	"GetPlugins",
	"LockTimeout",
	"Lockfile",
	"MigrateState",
	"Plugin",
	"Reconfigure",
	"TestDirectory",
	"Upgrade",
	"CompactWarnings",
	"Concurrency",
	"Lock",
	"NoColor",
	"Parallelism",
	"Refresh",
	"RefreshOnly",
	"DetailedExitcode",
	"GenerateConfigOut",
	"CreateBeforeDestroy",
	"Raw",
}

var terraformArrayFlagValues = []string{
	"Target",
	"Replace",
	"Exclude",
	"Var",
	"BackendConfig",
}
