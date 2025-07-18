package lib

func GetSourceSize() error {
	if err := loadConfig(); err != nil {
		return err
	}

	if err := getPackageSize(config.SourceLoc); err != nil {
		return err
	}

	return nil
}
