package vo

var _errors = map[int]*Error{}

func _newError(code int, sets ...SetErrorOption) *Error {
	_errors[code] = NewError(code, sets...)
	return _errors[code]
}

var (
	ErrSystem = _newError(500)
)
