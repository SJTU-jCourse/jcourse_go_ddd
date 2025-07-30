package app

import "jcourse_go/internal/config"

type ServiceContainer struct{}

func NewServiceContainer(conf config.Config) *ServiceContainer {
	return &ServiceContainer{}
}
