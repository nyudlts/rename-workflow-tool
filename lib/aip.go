package lib

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/nyudlts/go-aspace"
	bagit "github.com/nyudlts/go-bagit"
)

var transferPtn = regexp.MustCompile("transfer-info.txt")

func PrepAIP() error {
	fmt.Println("prepping aip")

	//load the config
	if err := loadConfig(); err != nil {
		return err
	}

	//set metadata dir
	mdDir = filepath.Join(config.SIPLoc, "metadata")

	//set the workorder path
	if err := getWorkOrderFile(mdDir); err != nil {
		return err
	}

	//load the workorder
	if err := validateWorkOrder(); err != nil {
		return err
	}

	var err error
	sipDirs, err = os.ReadDir(config.SIPLoc)
	if err != nil {
		return err
	}

	// for each line in the workorder
	for _, row := range workOrder.Rows {
		//get the corresponding directory in sip
		fmt.Println("finding target dir", row.GetComponentID())
		sourceDir, err := sipDirs.get(row.GetComponentID())
		if err != nil {
			return err
		}

		//create a directory in the aip directory with a UUID appended
		id := uuid.New().String()
		targetPath := filepath.Join(config.AIPLoc, sourceDir.Name()+"-"+id)
		if err := os.Mkdir(targetPath, 0755); err != nil {
			return err
		}

		//create a metadata directory
		aipMdDir := filepath.Join(targetPath, "metadata")
		if err := os.Mkdir(aipMdDir, 0755); err != nil {
			return err
		}

		//copy the transfer-info.txt to metadata

		transferInfo := filepath.Join(config.SIPLoc, "metadata", "transfer-info.txt")

		transferInfoBytes, err := os.ReadFile(transferInfo)
		if err != nil {
			return err
		}

		desc := "Internal-sender-description: " + row.GetTitle() + "\n"
		transferInfoBytes = append(transferInfoBytes, []byte(desc)...)

		targetTransferInfo := filepath.Join(aipMdDir, "transfer-info.txt")
		targetTransferInfoFile, err := os.Create(targetTransferInfo)
		if err != nil {
			return err
		}
		defer targetTransferInfoFile.Close()

		writer := bufio.NewWriter(targetTransferInfoFile)
		writer.Write(transferInfoBytes)
		writer.Flush()

		//generate new workorder for current line in workorder to metadatadir
		woName := filepath.Base(woPath)
		targetWorkOrderPath := filepath.Join(aipMdDir, woName)
		if _, err := os.Create(targetWorkOrderPath); err != nil {
			return err
		}

		woBody := strings.Join(aspace.HEADER_ROW, "\t")
		woBody = woBody + "\n"
		woBody = woBody + row.String()

		if err := os.WriteFile(targetWorkOrderPath, []byte(woBody), 0755); err != nil {
			return err
		}

		//move target to payload
		oldDir := filepath.Join(config.SIPLoc, sourceDir.Name())
		newDir := filepath.Join(targetPath, sourceDir.Name())
		if err := os.Rename(oldDir, newDir); err != nil {
			return err
		}

	}

	return nil
}

func BagAIP() error {
	if err := loadConfig(); err != nil {
		return err
	}

	aipDirs, err := os.ReadDir(config.AIPLoc)
	if err != nil {
		return err
	}

	for _, aipDir := range aipDirs {
		targetDir := filepath.Join(config.AIPLoc, aipDir.Name())
		fmt.Println("bagging", targetDir)

		logName := filepath.Join(config.LogLoc, "bagit", fmt.Sprintf("%s_bagit.log", aipDir.Name()))
		if _, err := os.Create(logName); err != nil {
			return err
		}
		bagCmd := exec.Command("bagit.py", "--sha256", targetDir)
		cmdOut, err := bagCmd.CombinedOutput()
		if err != nil {
			return err
		}

		if err := os.WriteFile(logName, cmdOut, 0644); err != nil {
			return err
		}

	}
	return nil
}

func UpdateAIP() error {
	if err := loadConfig(); err != nil {
		return err
	}

	aipDirs, err := os.ReadDir(config.AIPLoc)
	if err != nil {
		return err
	}

	for _, aipDir := range aipDirs {
		fmt.Println("updating", aipDir.Name())
		bagPath := filepath.Join(config.AIPLoc, aipDir.Name())

		bag, err := bagit.GetExistingBag(bagPath)
		if err != nil {
			return err
		}

		//Locate transfer-info.txt
		fmt.Printf("  * Locating transfer-info.txt: ")
		matches := bag.Payload.FindFilesInPayload(transferPtn)
		if len(matches) != 1 {
			return fmt.Errorf("no transfer-info.txt found")
		}
		tiPath := matches[0].Path
		tiPath = strings.ReplaceAll(tiPath+"/", bagPath, "")
		fmt.Printf("OK\n")

		//create a tag set from transfer-info.txt
		fmt.Printf("  * Creating new tag set from %s: ", "transfer-info.txt")
		transferInfo, err := bagit.NewTagSet(tiPath, bagPath)
		if err != nil {
			return err
		}
		fmt.Printf("OK\n")

		//Update the hostname
		fmt.Printf("  * Adding hostname to tag set: ")
		hostname, err := os.Hostname()
		if err != nil {
			return err
		}
		transferInfo.Tags["nyu-dl-hostname"] = hostname
		fmt.Printf("OK\n")

		//add pathname to the tag-set
		fmt.Printf("  * Adding bag's path to tag set: ")
		path, err := filepath.Abs(bagPath)
		if err != nil {
			return err
		}
		path, err = filepath.EvalSymlinks(path)
		if err != nil {
			return err
		}
		transferInfo.Tags["nyu-dl-pathname"] = path
		fmt.Printf("OK\n")

		//getting tagset from bag-info
		fmt.Printf("  * Creating new tag set from %s: ", "bag-info.txt")
		bagInfo, err := bagit.NewTagSet("bag-info.txt", bagPath)
		if err != nil {
			return err
		}
		fmt.Printf("OK\n")

		//merge tagsets
		fmt.Printf("  * Merging Tag Sets: ")
		bagInfo.AddTags(transferInfo.Tags)
		fmt.Printf("OK\n")

		fmt.Printf("  * Getting data as byte array: ")
		bagInfoBytes := bagInfo.GetTagSetAsByteSlice()
		fmt.Printf("OK\n")

		fmt.Printf("  * Opening bag-info.txt: ")
		bagInfoLocation := filepath.Join(bagPath, "bag-info.txt")
		bagInfoFile, err := os.Open(bagInfoLocation)
		if err != nil {
			return err
		}
		defer bagInfoFile.Close()
		fmt.Printf("OK\n")

		fmt.Printf("  * Rewriting bag-info.txt: ")
		if err := os.WriteFile(bagInfoLocation, bagInfoBytes, 0777); err != nil {
			return err
		}
		fmt.Printf("OK\n")

		//create new manifest object for tagmanifest-sha256.txt
		fmt.Printf("  * Creating new tagmanifest-sha256.txt: ")
		tagManifest, err := bagit.NewManifest(bagPath, "tagmanifest-sha256.txt")
		if err != nil {
			return err
		}
		fmt.Printf("OK\n")

		//update the checksum for bag-info.txt
		fmt.Printf("  * Updating checksum for bag-info.txt in tagmanifest-sha256.txt: ")
		if err := tagManifest.UpdateManifest("bag-info.txt"); err != nil {
			return err
		}
		fmt.Printf("OK\n")

		fmt.Printf("  * Rewriting tagmanifest-sha256.txt: ")
		if err := tagManifest.Serialize(); err != nil {
			return err
		}
		fmt.Printf("OK\n")

	}

	return nil
}

func ValidateAIPs() error {
	if err := loadConfig(); err != nil {
		return err
	}

	logFile, err := os.Create(fmt.Sprintf("logs/%s-aip-validation.log", config.CollectionCode))
	if err != nil {
		return err
	}

	defer logFile.Close()
	log.SetOutput(logFile)

	directoryEntries, err := os.ReadDir(config.AIPLoc)
	if err != nil {
		return err
	}

	for _, entry := range directoryEntries {
		if entry.IsDir() {

			aipPath := filepath.Join(config.AIPLoc, entry.Name())
			bag, err := bagit.GetExistingBag(aipPath)
			if err != nil {
				return err
			}

			fmt.Printf("  * validating %s\n", entry.Name())
			if err := bag.ValidateBag(false, false); err != nil {
				return err
			}
		}

	}

	return nil
}

func TransferAIPs() error {

	directoryEntries, err := os.ReadDir(config.AIPLoc)
	if err != nil {
		return err
	}

	xferLog := filepath.Join("logs", fmt.Sprintf("%s-aip-transfer.txt", config.CollectionCode))
	_, err = os.Create(xferLog)
	if err != nil {
		return err
	}

	for _, entry := range directoryEntries {
		fmt.Printf("  * transferring %s\n", entry.Name())
		xferBag := filepath.Join(config.AIPLoc, entry.Name())
		xferCmd := exec.Command("rstar-scp.exp", xferBag)
		cmdOutput, err := xferCmd.CombinedOutput()
		if err != nil {
			return err
		}
		cmdOutput = append(cmdOutput, []byte("\n")...)

		f, err := os.OpenFile(xferLog, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0775)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if _, err = f.Write(cmdOutput); err != nil {
			panic(err)
		}
	}
	return nil

}
