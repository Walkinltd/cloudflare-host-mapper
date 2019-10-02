package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// Load supports pulling the config from the cobra command or the environment variable under the key 'config'
func Load(cmd *cobra.Command, config interface{}) {
	flag := cmd.Flag("config")
	configStr := flag.Value.String()
	if configStr == "" {
		log.Printf("Config cmd flag not set. Using env variable")
		configStr = os.Getenv("CONFIG")
		if configStr == "" {
			log.Println("Config env not set. Using default config")
			return
		}
	}

	err := json.Unmarshal([]byte(configStr), config)
	if err != nil {
		log.Fatalf("Error loading config json: %s", err.Error())
	}
}
