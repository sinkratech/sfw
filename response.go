package sfw

// BaseResponseData is a wrapper for api response
type BaseResponseData[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// EmptyResponseData is an alias for *struct{}
type EmptyResponseData = *struct{}

// ResponseData is a client response. This must have Body property to fullfill huma.API standard
type ResponseData[T any] struct {
	Body BaseResponseData[T]
}

// RespondOKWith will generate response along with data
func RespondOKWith[T any](response T) *ResponseData[T] {
	return &ResponseData[T]{
		BaseResponseData[T]{
			Message: "OK",
			Data:    response,
		},
	}
}

// RespondOKWith will generate response but only {"message": "Your message here"}
func RespondMsg(msg string) *ResponseData[EmptyResponseData] {
	return &ResponseData[EmptyResponseData]{
		BaseResponseData[EmptyResponseData]{
			Message: msg,
			Data:    nil,
		},
	}
}

// RespondOK will generate ok but only {"message": "OK"}
func RespondOK() *ResponseData[EmptyResponseData] {
	return &ResponseData[EmptyResponseData]{
		BaseResponseData[EmptyResponseData]{
			Message: "OK",
			Data:    nil,
		},
	}
}
