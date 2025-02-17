package objects

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ErrorResponse(err error) Response {
	return Response{
		Message: err.Error(),
		Data:    nil,
	}
}
