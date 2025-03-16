package api

type notifierService struct {
}

// NewNotifierService ...
func NewNotifierService() Service {
	return &notifierService{}
}
