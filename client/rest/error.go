package rest_client

type ErrorResponse struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Fields  map[string]interface{} `json:"fields"`
}

func (e *ErrorResponse) Error() string {
	return e.Code
}
