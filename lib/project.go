package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func InitializeProject(collectionCode string, sourceLocation string) error {
	// check if collectionCode and sourceLocation are provided
	if collectionCode == "" || sourceLocation == "" {
		return fmt.Errorf("Collection code and source location must be provided")
	}

	//generate and print the configuration
	fmt.Println("* generating project configuration:")
	if err := generateConfig(collectionCode, sourceLocation); err != nil {
		return err
	}

	if err := PrintFormattedJson(config); err != nil {
		return err
	}

	//create project directories
	if err := mkProjectDir(); err != nil {
		return err
	}

	//write the config to the project directory
	if err := writeConfig(); err != nil {
		return err
	}

	fmt.Println("* Project initialized successfully at:", config.ProjectLoc)
	//return
	return nil
}

func generateConfig(collectionCode string, sourceLocation string) error {
	config = &Config{}

	// get the working directory
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	// update members
	config.PartnerCode = strings.Split(collectionCode, "_")[0]
	config.CollectionCode = collectionCode
	config.ProjectLoc = filepath.Join(wd, collectionCode)
	config.SIPLoc = filepath.Join(config.ProjectLoc, "sip")
	config.AIPLoc = filepath.Join(config.ProjectLoc, "aips")
	config.LogLoc = filepath.Join(config.ProjectLoc, "logs")
	config.SourceLoc = sourceLocation
	return nil
}

func mkProjectDir() error {
	//check if the project directory already exists
	if _, err := os.Stat(config.ProjectLoc); !os.IsNotExist(err) {
		return fmt.Errorf("Project directory %s already exists", config.ProjectLoc)
	}

	//create the project directory
	fmt.Println("* Creating project directory:", config.ProjectLoc)
	if err := os.Mkdir(config.ProjectLoc, 0775); err != nil {
		return err
	}

	//create the aips directory
	fmt.Println("* Creating aips directory:", config.AIPLoc)
	if err := os.Mkdir(filepath.Join(config.ProjectLoc, "aips"), 0775); err != nil {
		return err
	}

	//create the logs directory
	fmt.Println("* Creating logs directory:", config.LogLoc)
	if err := os.Mkdir(filepath.Join(config.ProjectLoc, "logs"), 0775); err != nil {
		return err
	}

	//create the rsync output directory
	fmt.Println("* Creating rsync directory:", config.LogLoc+"/rsync")
	if err := os.Mkdir(filepath.Join(config.ProjectLoc, "logs", "rsync"), 0775); err != nil {
		return err
	}

	//create the sip output directory
	fmt.Println("* Creating sip directory:", config.SIPLoc)
	if err := os.Mkdir(filepath.Join(config.ProjectLoc, "sip"), 0775); err != nil {
		return err
	}
	return nil
}

func writeConfig() error {

	//marshall the updated config
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}

	//write the config to the project directory
	if err := os.WriteFile(filepath.Join(config.ProjectLoc, "config.json"), b, 0755); err != nil {
		return err
	}

	return nil
}
