package domain

type Statistics struct {
	CertificateTotal      int `json:"certificateTotal"`
	CertificateExpireSoon int `json:"certificateExpireSoon"`
	CertificateExpired    int `json:"certificateExpired"`

	WorkflowTotal    int `json:"workflowTotal"`
	WorkflowEnabled  int `json:"workflowEnabled"`
	WorkflowDisabled int `json:"workflowDisabled"`
}
