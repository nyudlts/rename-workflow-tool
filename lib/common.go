package lib

import (
	"encoding/json"
	"fmt"
	"os"
)

var config *Config

type Config struct {
	SIPLoc         string `json:"sip-location"`
	SourceLoc      string `json:"source-location"`
	PartnerCode    string `json:"partner-code"`
	CollectionCode string `json:"collection-code"`
	ProjectLoc     string `json:"project-location"`
	LogLoc         string `json:"log-location"`
	AIPLoc         string `json:"aip-location"`
}

func PrintFormattedJson(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("Error formatting JSON: %s", err.Error())

	}

	fmt.Println(string(data))
	return nil
}

func loadConfig() error {
	//read the adoc-config
	b, err := os.ReadFile("config.json")
	if err != nil {
		return err
	}

	config = &Config{}
	//unmarshal to config options
	if err := json.Unmarshal(b, &config); err != nil {
		return err
	}

	return nil
}
