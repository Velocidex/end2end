package main

import (
	"os"

	"github.com/alecthomas/kingpin"
)

type CommandHandler func(command string) bool

var (
	app = kingpin.New("velo_e2e", "Velociraptor end to end tests")

	command_handlers []CommandHandler
)

func main() {
	app.HelpFlag.Short('h')
	app.UsageTemplate(kingpin.CompactUsageTemplate)
	args := os.Args[1:]

	command := kingpin.MustParse(app.Parse(args))

	for _, command_handler := range command_handlers {
		if command_handler(command) {
			break
		}
	}

}
