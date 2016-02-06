package cli

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

// groupCommand describe all dependencies of a group command
type groupCommand struct {
	cmd          *kingpin.CmdClause
	delimiter    *string
	transformers *Transformers
	parser       *Parser
	sorters      *Sorters
	positions    *Positions
	ignore       *bool
}

//  tokenizeCommand describe all dependencies of a tokenize command
type tokenizeCommand struct {
	cmd    *kingpin.CmdClause
	parser *Parser
}

// command list all available cli commands
type command map[string]*kingpin.CmdClause

func initApp() *kingpin.Application {
	return kingpin.New("goller", "Aggregate log fields and count occurences")
}

func initCmd(app *kingpin.Application) map[string]*kingpin.CmdClause {
	return map[string]*kingpin.CmdClause{
		"group":    app.Command("group", "Group occurence of field"),
		"tokenize": app.Command("tokenize", "Show how first log line is tokenized"),
	}
}

func initGroupArgs(groupCmd *kingpin.CmdClause) *groupCommand {
	return &groupCommand{
		delimiter:    groupCmd.Flag("delimiter", "Separator between results").Short('d').Default(" | ").String(),
		ignore:       groupCmd.Flag("ignore", "Ignore lines wrongly parsed").Short('i').Bool(),
		transformers: TransformersWrapper(groupCmd.Flag("transformer", "Transformers applied to every fields").Short('t')),
		sorters:      SortersWrapper(groupCmd.Flag("sort", "Sort lines").Short('s')),
		parser:       ParserWrapper(groupCmd.Arg("parser", "Log line parser to use").Required()),
		positions:    PositionsWrapper(groupCmd.Arg("positions", "Field positions").Required()),
	}
}

func initTokenizeArgs(tokenizeCmd *kingpin.CmdClause) *tokenizeCommand {
	return &tokenizeCommand{
		parser: ParserWrapper(tokenizeCmd.Arg("parser", "Log line parser to use").Required()),
	}
}

// Run commmand line arguments parsing
func Run(version string) {
	app := initApp()
	cmd := initCmd(app)
	groupArgs := initGroupArgs(cmd["group"])
	tokenizeArgs := initTokenizeArgs(cmd["tokenize"])

	app.Version(version)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case cmd["group"].FullCommand():
		triggerFatalError(groupArgs.sorters.ValidatePositions(groupArgs.positions.Get()))
		triggerFatalError(groupArgs.transformers.ValidatePositions(groupArgs.positions.Get()))

		group := NewGroup(groupArgs)
		group.Consume()
		group.Sort()
		group.Dispatch()
	case cmd["tokenize"].FullCommand():
		tokenize := NewTokenize(tokenizeArgs)
		tokenize.Tokenize()
		tokenize.Render()
	}
}
