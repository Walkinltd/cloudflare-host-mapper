package app

import (
	"hostmapper/services/cloudflare"
	"hostmapper/services/kubernetes"
)

type App interface {
	GetHosts() ([]Host, error)
	CreateRecords(hosts []Host) ([]string, error)
}

type app struct {
	Cloudflare cloudflare.Service
	Kubernetes kubernetes.Service
}

func New(cloudflare cloudflare.Service, kubernetes kubernetes.Service) App {
	return &app{
		Cloudflare: cloudflare,
		Kubernetes: kubernetes,
	}
}
