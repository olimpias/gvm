package common

func ValidateOperation() error {
	_, err := getGORoot()
	return err
}
