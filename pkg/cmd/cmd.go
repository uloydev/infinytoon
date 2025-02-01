package cmd

type Command interface {
	// Name returns the name of the command.
	Name() string
	// Run the command.
	Run() error
	// Shutdown the command.
	Shutdown() error
}
