package lib

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
	return nil
}
