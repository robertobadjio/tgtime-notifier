package endpoints

// LivenessResponse ...
type LivenessResponse struct {
	Code int    `json:"status"`
	Err  string `json:"err,omitempty"`
}

// ReadinessResponse ...
type ReadinessResponse struct {
	Code int    `json:"status"`
	Err  string `json:"err,omitempty"`
}
