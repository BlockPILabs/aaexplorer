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

	MonitorNotExist = _newError(20001, SetErrorMessage("The address is not in monitor list."))
	MonitorExist    = _newError(20002, SetErrorMessage("The address is already in monitor list."))
	NewUserErr      = _newError(20003, SetErrorMessage("The asset information is being synchronized."))
)
