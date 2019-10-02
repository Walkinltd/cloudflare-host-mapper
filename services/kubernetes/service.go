package kubernetes

import (
	"errors"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Service interface {
	GetIngressList() ([]string, error)
	GetHostsForIngress(ingress string) ([]string, string, error)
}

type service struct {
	Client    *kubernetes.Clientset
	Namespace string
}

func New(namespace string) (Service, error) {
	// Pull the k8s config from the cluster
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// Creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &service{
		Client:    clientset,
		Namespace: namespace,
	}, nil
}

func (s *service) GetIngressList() ([]string, error) {
	client := s.Client.ExtensionsV1beta1().Ingresses(s.Namespace)

	res, err := client.List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	ingresses := make([]string, len(res.Items))
	for i, ingress := range res.Items {
		ingresses[i] = ingress.Name
	}

	return ingresses, nil
}

func (s *service) GetHostsForIngress(ingress string) ([]string, string, error) {
	client := s.Client.ExtensionsV1beta1().Ingresses(s.Namespace)

	ing, err := client.Get(ingress, v1.GetOptions{})
	if err != nil {
		return nil, "", err
	}

	hosts := make([]string, len(ing.Spec.Rules))
	for i, rule := range ing.Spec.Rules {
		hosts[i] = rule.Host
	}

	if len(ing.Status.LoadBalancer.Ingress) == 0 {
		return nil, "", errors.New("no_ingress_available")
	}

	return hosts, ing.Status.LoadBalancer.Ingress[0].IP, nil
}
