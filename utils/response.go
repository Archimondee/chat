package utils

type response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func ResponseData(status string, message string, data any) *response {
	return &response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
