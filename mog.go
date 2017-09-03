package main

import (
	"os"
	"runtime"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

// VERSION is the cli version
const VERSION = "v0.1.1"

func main() {
	app := cli.NewApp()
	app.Version = VERSION
	app.Usage = "A CLI Tool for Digdag"
	app.Commands = Commands
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "host, H",
			Usage: "digdag server's hostname or ip address",
			Value: "localhost",
		},
		cli.IntFlag{
			Name:  "port, P",
			Usage: "digdag server's port number",
			Value: 65432,
		},
		cli.BoolFlag{
			Name:  "ssl",
			Usage: "make `https` request",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "verbose output",
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.GlobalBool("verbose") == true {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	}
	app.OnUsageError = CustomOnUsageError

	cpu := runtime.NumCPU()
	runtime.GOMAXPROCS(cpu)
	app.Run(os.Args)
}

// CustomOnUsageError is show custom message and show usage
func CustomOnUsageError(c *cli.Context, err error, isSubcommand bool) error {
	fmt.Fprintf(c.App.Writer, "Error: %s is required \n", err)
	cli.ShowCommandHelp(c, c.Command.Name)
	os.Exit(1)
	return nil
}
