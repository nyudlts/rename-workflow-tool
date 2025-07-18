package lib

func GetSourceSize() error {
	if err := loadConfig(); err != nil {
		return err
	}
	return nil
}
