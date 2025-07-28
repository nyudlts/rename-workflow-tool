package lib

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nyudlts/go-aspace"
)

var (
	woPath       string
	workOrder    aspace.WorkOrder
	componentIDs []string
)

func getWorkOrderFile(path string) error {
	mdFiles, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, mdFile := range mdFiles {
		name := mdFile.Name()
		if strings.Contains(name, "_aspace_wo.tsv") {
			woPath = filepath.Join(path, name)
			return nil
		}
	}
	return fmt.Errorf("%s does not contain a work order", path)
}

func validateWorkOrder() error {
	wof, err := os.Open(woPath)
	if err != nil {
		return err
	}
	defer wof.Close()

	if err := workOrder.Load(wof); err != nil {
		return err
	}

	return nil
}

func checkWorkOrderForDuplicates() error {
	componentIDs = []string{}
	//get an array of componentIDs
	dupeCount := 0
	for _, row := range workOrder.Rows {
		if contains(row.GetComponentID(), componentIDs) {
			log.Printf("[ERROR] duplicate componentID, %s, found in workorder\n", row.GetComponentID())
			dupeCount++
		} else {
			componentIDs = append(componentIDs, row.GetComponentID())
		}
	}
	sort.Strings(componentIDs)
	log.Printf("[INFO] check 5. %s contains %d duplicate cuids \n", filepath.Base(woPath), dupeCount)
	if dupeCount > 0 {
		return fmt.Errorf("work order contains %d duplicate component IDs", dupeCount)
	} else {
		return nil
	}
}

func contains(s string, sl []string) bool {
	for _, sls := range sl {
		if s == sls {
			return true
		}
	}
	return false
}
