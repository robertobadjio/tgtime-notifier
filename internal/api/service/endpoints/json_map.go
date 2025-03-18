package endpoints

// LivenessRequest ...
type LivenessRequest struct{}

// LivenessResponse ...
type LivenessResponse struct {
	Code int    `json:"status"`
	Err  string `json:"err,omitempty"`
}

// ReadinessRequest ...
type ReadinessRequest struct{}

// ReadinessResponse ...
type ReadinessResponse struct {
	Code int    `json:"status"`
	Err  string `json:"err,omitempty"`
}
