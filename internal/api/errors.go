package api

type APIError struct {
	Path       string
	Message    string
	StatusCode int
	Err        error // check in Stripe API error handling
}
