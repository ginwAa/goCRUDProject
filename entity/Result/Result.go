package Result

import "net/http"

type Result struct {
	Status int
	Msg    string
	Data   interface{}
}

func Error(status int, msg string) Result {
	return Result{
		Status: status,
		Msg:    msg,
		Data:   nil,
	}
}

func BadRequest() Result {
	return Error(http.StatusBadRequest, "Request error!")
}

func Success(data interface{}) Result {
	return Result{
		Status: http.StatusOK,
		Msg:    "",
		Data:   data,
	}
}
