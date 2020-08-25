package commands

// Command required interface to be implemented by commands in order to function.
type Command interface {

	// CommandName returns the name of the command.
	CommandName() string
	// Parse should parse the command and return any errors that occur during the process
	Parse([]string) error
	// Run should contain the logic to execute the desired command
	Run() error
}
