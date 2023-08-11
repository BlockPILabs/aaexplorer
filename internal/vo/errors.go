package vo

var _errors = map[int]*Error{}

func _newError(code int, sets ...SetErrorOption) *Error {
	_errors[code] = NewError(code, sets...)
	return _errors[code]
}

var (
	// todo common error
	ErrSystem          = _newError(10001, SetErrorMessage("system error"))
	ErrParams          = _newError(10002, SetErrorMessage("params error"))
	ErrDataNotFound    = _newError(10004, SetErrorMessage("data not found"))
	ErrNetworkNotFound = ErrDataNotFound.SetMessage("network not found")
)
