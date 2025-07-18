package lib

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func PrintSourceSize() error {
	if err := loadConfig(); err != nil {
		return err
	}

	if err := printPackageSize(config.SourceLoc); err != nil {
		return err
	}

	return nil
}

func TransferSourcePackage() error {
	//create a log file for source transfer
	if err := loadConfig(); err != nil {
		return err
	}
	logFileName := filepath.Join(config.LogLoc, "rsync", fmt.Sprintf("%s-source-transfer-rsync.txt", config.CollectionCode))
	logFile, err := os.Create(logFileName)
	if err != nil {
		return err
	}
	defer logFile.Close()

	// create a buffered writer for the log file
	writer := bufio.NewWriter(logFile)

	fmt.Printf("  * Transferring %s to sip directory\n", config.SourceLoc)

	// run the rsync command to transfer the source package
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("robocopy.exe", config.SourceLoc, config.SIPLoc, "/DCOPY:DAT", "/E")
	} else {
		cmd = exec.Command("rsync", "-av", config.SourceLoc, config.SIPLoc)
	}
	b, err := cmd.CombinedOutput()
	if err != nil {
		return nil
	}

	if _, err := writer.Write(b); err != nil {
		return err
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	fmt.Printf("  * Transfer complete\n")

	mdDir := filepath.Join(config.SIPLoc, "metadata")
	if _, err := os.Stat(mdDir); os.IsNotExist(err) {
		fmt.Printf("  * Creating metadata directory in SIP\n")
		if err := os.Mkdir(mdDir, 0755); err != nil {
			return err
		}
	}

	return nil
}
