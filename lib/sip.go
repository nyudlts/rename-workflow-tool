package lib

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type SIPDirs []os.DirEntry

var (
	mdDir   string
	sipDirs SIPDirs
)

func (sds SIPDirs) get(s string) (os.DirEntry, error) {
	for _, sipDir := range sds {
		if strings.Contains(sipDir.Name(), s) {
			return sipDir, nil
		}
	}
	return nil, fmt.Errorf("no corresponding sip found")
}

func (sds SIPDirs) contains(s string) bool {
	for _, sipDir := range sds {
		if strings.Contains(sipDir.Name(), s) {
			return true
		}
	}
	return false
}

func GetSipSize() error {
	fmt.Printf("rwt sip size, %s\n", VERSION)
	if err := loadConfig(); err != nil {
		return err
	}

	if err := printPackageSize(config.SIPLoc); err != nil {
		return err
	}

	return nil
}

func ValidateSip() error {
	fmt.Printf("rwt sip validate, %s\n", VERSION)
	validationError := false
	//load the project configuration
	if err := loadConfig(); err != nil {
		return fmt.Errorf("error loading configuration: %v", err)
	}

	//create a logger
	logFile, err := os.Create(filepath.Join("logs", fmt.Sprintf("%s-sip-validate.log", config.CollectionCode)))
	if err != nil {
		return fmt.Errorf("error creating log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	fmt.Printf("  * validating SIP transfer package at %s\n", config.SIPLoc)
	log.Printf("[INFO] validating SIP transfer package at %s\n", config.SIPLoc)

	//check if the SIP Directory exists
	fmt.Printf("  * checking SIP directory exists\n")
	if err := checkSIPDirExists(); err != nil {
		log.Printf("[ERROR] %s\n", err)
		fmt.Printf("  [ERROR] %s\n", err)
		validationError = true
	}

	//check that there is a metadata directory
	fmt.Printf("  * checking metadata directory exists\n")
	if err := checkMetadataDirExists(); err != nil {
		log.Printf("[ERROR] %s\n", err)
		fmt.Printf("  [ERROR] %s\n", err)
		validationError = true
	}

	//check if a work order exists
	fmt.Printf("  * checking work order exists\n")
	if err := getWorkOrderFile(mdDir); err != nil {
		log.Printf("[ERROR] %s\n", err)
		fmt.Printf("  [ERROR] %s\n", err)
		validationError = true
	}

	//check that work order is valid
	fmt.Printf("  * checking work order is valid\n")
	if err := validateWorkOrder(); err != nil {
		log.Printf("[ERROR] %s\n", err)
		fmt.Printf("  [ERROR] %s\n", err)
		validationError = true
	}

	//check work order for duplicates
	fmt.Printf("  * checking work order for duplicates\n")
	if err := checkWorkOrderForDuplicates(); err != nil {
		log.Printf("[ERROR] %s\n", err)
		fmt.Printf("  [ERROR] %s\n", err)
		validationError = true
	}

	//check sip for missing directories
	fmt.Printf("  * checking SIP for missing directories\n")
	if err := checkSIPForMissingDirectories(); err != nil {
		log.Printf("[ERROR] %s\n", err)
		validationError = true
	}

	//check for transfer info.txt
	fmt.Println("  * checking valid transfer-info.txt exists")
	if err := validateTransferInfo(); err != nil {
		log.Printf("[ERROR] %s\n", err)
		validationError = true
	}

	if validationError {
		return fmt.Errorf("SIP is not valid, see logs for details")
	}
	return nil
}

func checkSIPDirExists() error {
	if _, err := os.Stat(config.SIPLoc); os.IsNotExist(err) {
		return fmt.Errorf("SIP directory does not exist: %s", config.SIPLoc)
	}
	return nil
}

func checkMetadataDirExists() error {
	mdDir = filepath.Join(config.SIPLoc, "metadata")
	if _, err := os.Stat(mdDir); os.IsNotExist(err) {
		return fmt.Errorf("metadata directory does not exist: %s", mdDir)
	}
	return nil
}

func checkSIPForMissingDirectories() error {
	//get a list of dirs in payload
	var err error
	sipDirs, err = os.ReadDir(config.SIPLoc)
	if err != nil {
		return err
	}

	missingDirs := 0
	for _, row := range workOrder.Rows {
		componentID := row.GetComponentID()
		if !sipDirs.contains(componentID) {
			missingDirs++
			log.Printf("[ERROR] componentID, %s is missing in transfered directories\n", componentID)
			fmt.Printf("    * cuid %s is missing from transferred directories\n", componentID)
		}
	}

	if missingDirs > 0 {
		return fmt.Errorf("SIP contains %d missing directories, see validation log", missingDirs)
	} else {
		return nil
	}
}

func checkSIPForExtraDirectories() error {
	sourceDirs, err := os.ReadDir(config.SIPLoc)
	if err != nil {
		log.Printf("[ERROR] could not read SIP directory: %s", err)
		fmt.Printf("  [ERROR] could not read SIP directory: %s\n", err)
		return err
	} else {
		extraDirs := 0
		for _, sourceDir := range sourceDirs {
			if sourceDir.Name() != "metadata" {
				if !contains(sourceDir.Name(), componentIDs) {
					extraDirs++
					log.Printf("[ERROR] %s is not listed on workorder\n", sourceDir.Name())
				}
			}
		}

		if extraDirs > 0 {
			return fmt.Errorf("SIP contains %d extra directories", extraDirs)
		} else {
			return nil
		}
	}

}

func checkTransferInfoExists() error {
	return nil
}
