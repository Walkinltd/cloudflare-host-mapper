package app

func (a *app) CreateRecords(hosts []Host) ([]string, error) {
	ids := make([]string, len(hosts))

	for i, host := range hosts {
		id, err := a.Cloudflare.CreateRecord(host.Path, host.IP)
		if err != nil {
			return nil, err
		}

		ids[i] = id
	}
	return ids, nil
}
