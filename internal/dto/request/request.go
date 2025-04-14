package request

type Request interface {
	ValidateRequest() error
}
