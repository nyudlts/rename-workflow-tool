package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nyudlts/bytemath"
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

func printPackageSize(pkgPath string) error {
	numFiles := 0
	numDirectories := 0
	sizeFiles := int64(0)

	if err := filepath.Walk(pkgPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			numDirectories++
		} else {
			numFiles++
			sizeFiles += info.Size()
		}
		return nil
	}); err != nil {
		return err
	}

	fmt.Printf("%s: %d files in %d directories, %s\n", pkgPath, numFiles, numDirectories, bytemath.ConvertBytesToHumanReadable(sizeFiles))
	return nil
}
