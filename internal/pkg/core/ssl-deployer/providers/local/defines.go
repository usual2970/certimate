package local

type OutputFormatType string

const (
	OUTPUT_FORMAT_PEM = OutputFormatType("PEM")
	OUTPUT_FORMAT_PFX = OutputFormatType("PFX")
	OUTPUT_FORMAT_JKS = OutputFormatType("JKS")
)

type ShellEnvType string

const (
	SHELL_ENV_SH         = ShellEnvType("sh")
	SHELL_ENV_CMD        = ShellEnvType("cmd")
	SHELL_ENV_POWERSHELL = ShellEnvType("powershell")
)
