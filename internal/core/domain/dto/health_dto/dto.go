package health_dto

type HealthStatus struct {
	Status string `json:"status"`
	Err    string `json:"err,omitempty"`
}

func New(err error) *HealthStatus {
	status := "ok"
	errMsg := ""
	if err != nil {
		status = "error"
		errMsg = err.Error()
	}

	return &HealthStatus{
		Status: status,
		Err:    errMsg,
	}
}
