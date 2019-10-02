package main

import (
	"log"
	"os"
	"strings"
	"time"

	"hostmapper/app"
	"hostmapper/config"
	"hostmapper/services/cloudflare"
	"hostmapper/services/kubernetes"

	"github.com/spf13/cobra"
)

var rootCMD = &cobra.Command{
	Use:   "host-mapper",
	Short: "Host mapper maps hosts to DNS records",
	Long:  "Cloudflare host mapper pulls hosts from a K8s ingress and maps them to DNS",
	RunE:  runE,
}

type CloudflareSettings struct {
	AccountID string `json:"account_id"`
	ZoneID    string `json:"zone_id"`
	Token     string `json:"token"`
}

type Config struct {
	Cloudflare    CloudflareSettings `json:"cloudflare"`
	Namespace     string             `json:"kube_namespace"`
	MappableHosts []string           `json:"mappable_hosts"`
}

func main() {
	rootCMD.PersistentFlags().String("config", "", "Config variable json")

	// Run command and exit with error code upon error
	if err := rootCMD.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}

func runE(cmd *cobra.Command, args []string) error {
	defaultConfig := &Config{
		Cloudflare: CloudflareSettings{
			AccountID: "",
			ZoneID:    "",
			Token:     "",
		},

		Namespace: "default",

		MappableHosts: []string{},
	}
	config.Load(cmd, defaultConfig)

	cloudflareSvc, err := cloudflare.New(defaultConfig.Cloudflare.ZoneID, defaultConfig.Cloudflare.AccountID, defaultConfig.Cloudflare.Token)
	if err != nil {
		return err
	}

	kubernetesSvc, err := kubernetes.New(defaultConfig.Namespace)
	if err != nil {
		return err
	}

	app := app.New(cloudflareSvc, kubernetesSvc)

	for {
		log.Println("Running ingress scan")

		err = mapHosts(app, defaultConfig.MappableHosts)
		if err != nil {
			log.Fatal(err.Error())
		}

		// Sleep for 30 seconds before running again
		time.Sleep(30 * time.Second)
	}

	return nil
}

func mapHosts(a app.App, mappableHosts []string) error {
	// Get the hosts from k8s
	hosts, err := a.GetHosts()
	if err != nil {
		return err
	}

	// List matching hosts only
	availableHosts := make([]app.Host, 0)
	for _, host := range hosts {
		found := false
		for _, mappable := range mappableHosts {
			if strings.HasSuffix(host.Path, mappable) {
				availableHosts = append(availableHosts, host)
				found = true
				break
			}
		}

		if !found {
			// Log that the host doesn't match any available mappable hosts
			log.Printf("Skipping: %s. Does not match any suffix", host.Path)
		}
	}

	if len(availableHosts) == 0 {
		// No hosts found so skip creating records
		return nil
	}

	// Create the records within cloudflare
	_, err = a.CreateRecords(availableHosts)
	return err
}
