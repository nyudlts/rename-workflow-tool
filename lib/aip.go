package lib

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/google/uuid"
	"github.com/nyudlts/go-aspace"
	bagit "github.com/nyudlts/go-bagit"
)

var transferPtn = regexp.MustCompile("transfer-info.txt")

func PrepAIPs() error {
	fmt.Printf("rwt aip prep, %s\n", VERSION)

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

		aipPath, err := prepareAIP(row)
		if err != nil {
			return err
		}

		if err := bagAIP(aipPath); err != nil {
			return err
		}

		if err := updateAIP(aipPath); err != nil {
			return err
		}

	}

	return nil

}

func prepareAIP(row aspace.WorkOrderRow) (string, error) {
	cuid := row.GetComponentID()
	fmt.Printf("  * preparing aip for %s\n", cuid)

	//get the corresponding directory in sip
	sourceDir, err := sipDirs.get(cuid)
	if err != nil {
		return "", err
	}

	//create a directory in the aip directory with a UUID appended
	id := uuid.New().String()
	fmt.Printf("  * creating aip directory with id: %s\n", id)
	targetPath := filepath.Join(config.AIPLoc, sourceDir.Name()+"-"+id)
	if err := os.Mkdir(targetPath, 0755); err != nil {
		return "", err
	}

	//copy the transfer-info.txt to metadata
	transferInfo := filepath.Join(config.SIPLoc, "metadata", "transfer-info.txt")

	transferInfoBytes, err := os.ReadFile(transferInfo)
	if err != nil {
		return "", err
	}

	desc := "Internal-sender-description: " + row.GetTitle() + "\n"
	transferInfoBytes = append(transferInfoBytes, []byte(desc)...)

	targetTransferInfo := filepath.Join(config.TmpLoc, fmt.Sprintf("%s-transfer-info.txt", id))
	if err := os.WriteFile(targetTransferInfo, transferInfoBytes, 0755); err != nil {
		return "", err
	}

	//move target to payload
	oldDir := filepath.Join(config.SIPLoc, sourceDir.Name())
	newDir := filepath.Join(targetPath, sourceDir.Name())
	if err := os.Rename(oldDir, newDir); err != nil {
		return "", err
	}

	return targetPath, nil
}

func bagAIP(pkgPath string) error {
	fmt.Println("  * bagging", filepath.Base(pkgPath))

	//bag the transfer
	var bagCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		bagCmd = exec.Command("python", "-m", "bagit", "--sha256", pkgPath)
	} else {
		bagCmd = exec.Command("bagit.py", "--sha256", pkgPath)
	}

	cmdOut, err := bagCmd.CombinedOutput()
	if err != nil {
		return err
	}

	logName := filepath.Join(config.LogLoc, "bagit", fmt.Sprintf("%s-bagit.log", config.CollectionCode))
	logFile, err := os.OpenFile(logName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	if _, err := logFile.Write(cmdOut); err != nil {
		return err
	}

	return nil
}

func BagAIPs() error {
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
		if err := bagAIP(targetDir); err != nil {
			return err
		}

	}
	return nil
}

func updateAIP(pkgPath string) error {

	fmt.Println("updating", filepath.Base(pkgPath))

	bag, err := bagit.GetExistingBag(pkgPath)
	if err != nil {
		return err
	}

	dirUUID := bag.Name[len(bag.Name)-36:]

	//Locate transfer-info.txt
	fmt.Println("  * Locating transfer-info.txt")
	tiFilename := fmt.Sprintf("%s-transfer-info.txt", dirUUID)

	//create a tag set from transfer-info.txt
	fmt.Printf("  * Creating new tag set from %s\n", "transfer-info.txt")
	transferInfo, err := bagit.NewTagSet(tiFilename, config.TmpLoc)
	if err != nil {
		return err
	}

	//Update the hostname
	fmt.Println("  * Adding hostname to tag set")
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	transferInfo.Tags["nyu-dl-hostname"] = hostname

	//add pathname to the tag-set
	fmt.Printf("  * Adding bag's path to tag set: ")
	path, err := filepath.Abs(pkgPath)
	if err != nil {
		return err
	}
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}
	transferInfo.Tags["nyu-dl-pathname"] = path

	//getting tagset from bag-info
	fmt.Printf("  * Creating new tag set from %s\n", "bag-info.txt")
	bagInfo, err := bagit.NewTagSet("bag-info.txt", pkgPath)
	if err != nil {
		return err
	}

	//merge tagsets
	fmt.Println("  * Merging Tag Sets")
	bagInfo.AddTags(transferInfo.Tags)

	fmt.Printf("  * Getting data as byte array\n")
	bagInfoBytes := bagInfo.GetTagSetAsByteSlice()

	fmt.Printf("  * Opening bag-info.txt\n")
	bagInfoLocation := filepath.Join(pkgPath, "bag-info.txt")
	bagInfoFile, err := os.Open(bagInfoLocation)
	if err != nil {
		return err
	}
	defer bagInfoFile.Close()

	fmt.Printf("  * Rewriting bag-info.txt\n")
	if err := os.WriteFile(bagInfoLocation, bagInfoBytes, 0777); err != nil {
		return err
	}
	//create new manifest object for tagmanifest-sha256.txt
	fmt.Printf("  * Creating new tagmanifest-sha256.txt\n")
	tagManifest, err := bagit.NewManifest(pkgPath, "tagmanifest-sha256.txt")
	if err != nil {
		return err
	}

	//update the checksum for bag-info.txt
	fmt.Printf("  * Updating checksum for bag-info.txt in tagmanifest-sha256.txt\n")
	if err := tagManifest.UpdateManifest("bag-info.txt"); err != nil {
		return err
	}

	fmt.Printf("  * Rewriting tagmanifest-sha256.txt\n")
	if err := tagManifest.Serialize(); err != nil {
		return err
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
		fmt.Println(aipDir.Name())
	}

	return nil
}

func ValidateAIPs() error {
	fmt.Printf("rwt aip validate, %s\n", VERSION)
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
			if runtime.GOOS == "windows" {
				valCmd := exec.Command("python", "-m", "bagit", "--validate", aipPath)
				cmdOut, err := valCmd.CombinedOutput()
				if err != nil {
					return err
				}
				fmt.Println(string(cmdOut))
			} else {
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
	}

	return nil
}

func TransferAIPs() error {
	fmt.Printf("rwt aip transfer, %s\n", VERSION)

	if err := loadConfig(); err != nil {
		return err
	}

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
