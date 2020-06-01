package internal

type Command interface {
	Validate() error
	Apply() error
}
