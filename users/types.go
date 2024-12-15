package main

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Error   string `json:"error,omitempty"`
	Code    int    `json:"code"`
}

func NewResponse(message string, data any, err *string, code int) *Response {
	res := &Response{
		Message: message,
		Code:    code,
	}

	if err != nil {
		res.Error = *err
	}

	if data != nil {
		res.Data = data
	}

	return res
}
