package app

type Host struct {
	Path string `json:"path"`
	IP   string `json:"ip"`
}

func (a *app) GetHosts() ([]Host, error) {
	ingresses, err := a.Kubernetes.GetIngressList()
	if err != nil {
		return nil, err
	}

	hosts := make([]Host, 0)
	for _, ingress := range ingresses {
		ingressHosts, ip, err := a.Kubernetes.GetHostsForIngress(ingress)
		if err != nil {
			return nil, err
		}

		for _, host := range ingressHosts {
			hosts = append(hosts, Host{
				Path: host,
				IP:   ip,
			})
		}
	}

	return hosts, nil
}
