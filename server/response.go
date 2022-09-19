package server

type Response[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func EmptyResponse() Response[any] {
	return Response[any]{
		Message: "EMPTY",
		Data:    nil,
	}
}

func OKResponse[T any](data T) Response[T] {
	return Response[T]{
		Message: "OK",
		Data:    data,
	}
}

func ErrorResponse(err error) Response[string] {
	return Response[string]{
		Message: "ERROR",
		Data:    err.Error(),
	}
}
