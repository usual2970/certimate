package dtos

import "time"

type PurgeResponse struct {
	ID                 string    `json:"id"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	CompletedAt        time.Time `json:"completed_at"`
	ProgressPercentage float32   `json:"progress_percentage"`
}

type PurgeRequest struct {
	EnvironmentID string   `json:"environment_id"`
	PurgeType     string   `json:"purge_type"`
	Values        []string `json:"values"`
	Hostname      *string  `json:"hostname"`
}
