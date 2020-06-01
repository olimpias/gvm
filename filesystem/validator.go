package filesystem

func ValidateOperation() error {
	_, err := getGORoot()
	return err
}
