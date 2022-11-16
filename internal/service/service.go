package service

type HealthCheck interface {
	Check() (bool, error)
}

type Service struct {
	HealthCheck
}
