package cmd

type Command interface {
	Validate() error
	Apply() error
}
