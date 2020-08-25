package commands

import "flag"

type (

	// TorikumiCommand command and flagset to be executed.
	TorikumiCommand struct {
		TorikumiFlagSet *flag.FlagSet

		bashoID int
		day     int
	}
)

//NewTorikumiCommand returns a new ToeikumiCommand type
func NewTorikumiCommand() *TorikumiCommand {

	cmd := &TorikumiCommand{
		TorikumiFlagSet: flag.NewFlagSet("torikumi", flag.ExitOnError),
	}

	cmd.TorikumiFlagSet.IntVar(&cmd.bashoID, "basho-id", -1, "The basho to target <YYYYMM>")
	cmd.TorikumiFlagSet.IntVar(&cmd.day, "day", -1, "the day to get bouts for must be a value between 1-16")

	return cmd
}

//CommandName returns command name.
func (cmd *TorikumiCommand) CommandName() string {
	return cmd.TorikumiFlagSet.Name()
}

//Parse parses command arguments and returns an error if any of the values are invalid.
func (cmd *TorikumiCommand) Parse(osArgs []string) error {
	cmd.TorikumiFlagSet.Parse(osArgs)
	return nil
}

//Run logic to be executed when the torikumi command is called.
func (cmd *TorikumiCommand) Run() error {
	//TODO: add logic
	return nil
}
