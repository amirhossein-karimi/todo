package response

type response struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
}

func CreateResponse(code int, message string, status string, data interface{}, err interface{}) response {
	return response{
		Status:  status,
		Code:    code,
		Message: message,
		Data:    data,
		Errors:  err,
	}
}
