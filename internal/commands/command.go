package commands

type Command interface {
	Validate() error
	Apply() error
}
