package service

type Authorization interface {
}

type Service struct {
	Authorization
}

func NewService() *Service {
	return &Service{}
}
