package cloudflare

import (
	"fmt"

	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
)

type Service interface {
	CreateRecord(host, ip string) (string, error)
}

type service struct {
	ZoneID    string
	AccountID string
	Client    *cloudflare.API
}

func New(zoneID string, accountID string, token string) (Service, error) {
	api, err := cloudflare.NewWithAPIToken(token)
	if err != nil {
		return nil, err
	}

	return &service{
		ZoneID:    zoneID,
		AccountID: accountID,
		Client:    api,
	}, nil
}

func (c *service) CreateRecord(host, ip string) (string, error) {
	record, err := c.getRecord(host)
	if err != nil {
		return "", err
	}

	if record == nil {
		record = &cloudflare.DNSRecord{Name: host}
	} else if record.Type == "A" && record.Content == ip && record.Proxied == true {
		// Return record as data is already matching
		return record.ID, nil
	}

	record.Type = "A"
	record.Content = ip
	record.Proxied = true
	record.TTL = 1
	record.ZoneID = c.ZoneID

	if record.ID == "" {
		resp, err := c.Client.CreateDNSRecord(c.ZoneID, *record)
		if err != nil {
			return "", err
		}

		if !resp.Success {
			return "", errors.New(fmt.Sprintf("%v: %s", resp.Response.Errors[0].Code, resp.Response.Errors[0].Message))
		}

		return resp.Result.ID, nil
	}

	err = c.Client.UpdateDNSRecord(c.ZoneID, record.ID, *record)
	if err != nil {
		return "", err
	}

	return record.ID, nil
}

func (c *service) getRecord(host string) (*cloudflare.DNSRecord, error) {
	record := cloudflare.DNSRecord{
		Type: "A",
		Name: host,
	}

	recs, err := c.Client.DNSRecords(c.ZoneID, record)
	if err != nil {
		return nil, err
	}

	if len(recs) != 0 {
		return &recs[0], nil
	}

	return nil, nil
}
