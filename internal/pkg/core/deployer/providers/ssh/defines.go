package ssh

type OutputFormatType string

const (
	OUTPUT_FORMAT_PEM = OutputFormatType("PEM")
	OUTPUT_FORMAT_PFX = OutputFormatType("PFX")
	OUTPUT_FORMAT_JKS = OutputFormatType("JKS")
)
