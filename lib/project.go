package lib

import "fmt"

func InitializeProject(collectionCode string, sourceLocation string) error {
	// Logic to initialize a project with the given collection code and source location
	if collectionCode == "" || sourceLocation == "" {
		return fmt.Errorf("Collection code and source location must be provided")
	}
	return nil
}
