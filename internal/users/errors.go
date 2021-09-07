package users

// maybe api error repack somewhere?
// it should be some api error struct to display a nice json
type ErrRecordNotFound struct { // should be a more common name, not so specific, it's better to put a specific type into "err" var
	Message    string
	StatusCode int
	Err        error
}

// there is a template for this in goapp on github
func (e *ErrRecordNotFound) Error() string {
	return "Error: " + e.Err.Error() + ". Message: " + e.Message
}
