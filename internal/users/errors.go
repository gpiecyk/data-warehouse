package users

type ErrRecordNotFound struct {
	Message    string
	StatusCode int
	Err        error
}

func (e *ErrRecordNotFound) Error() string {
	return "Error: " + e.Err.Error() + ". Message: " + e.Message
}
