package lib

import (
	"encoding/json"
	"fmt"
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
