package endpoints

// LivenessRequest ...
type LivenessRequest struct{}

// LivenessResponse ...
type LivenessResponse struct {
	Code int    `json:"status"`
	Err  string `json:"err,omitempty"`
}

/*func (r LivenessResponse) Error() error {
	return r.Err
}*/

// ReadinessRequest ...
type ReadinessRequest struct{}

// ReadinessResponse ...
type ReadinessResponse struct {
	Code int    `json:"status"`
	Err  string `json:"err,omitempty"`
}

/*func (r ReadinessResponse) Error() error {
	return r.Err
}*/
