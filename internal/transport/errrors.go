package transport

type ErrorType string

type SlugError struct {
	error     string
	slug      string
	errorType ErrorType
}

func (s SlugError) Error() string {
	return s.error
}
func (s SlugError) Slug() string {
	return s.slug
}
func (s SlugError) ErrorType() ErrorType {
	return s.errorType
}
