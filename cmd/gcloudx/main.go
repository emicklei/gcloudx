package main

import (
	"log"
	"os"
	"sort"

	"github.com/emicklei/gcloudx/ps"
	"github.com/urfave/cli"
)

var Version string

func main() {
	if err := newApp().Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Version = Version
	app.EnableBashCompletion = true
	app.Name = "gcloudx"
	app.Usage = "Extra features to manage Google Cloud Platform"

	// override -v
	cli.VersionFlag = cli.BoolFlag{
		Name:  "print-version, V",
		Usage: "print only the version",
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "v",
			Usage: "verbose logging",
		},
		cli.BoolFlag{
			Name:  "q",
			Usage: "quiet mode, accept any prompt",
		},
	}
	projectFlag := cli.StringFlag{
		Name:  "p",
		Usage: `GCP project identifier`,
	}
	topicFlag := cli.StringFlag{
		Name:  "t",
		Usage: `PubSub topic identifier (short name)`,
	}
	fileFlag := cli.StringFlag{
		Name:  "f",
		Usage: `File containing the payload`,
	}
	app.Commands = []cli.Command{
		{
			Name:  "pubsub",
			Usage: "Work with Pub/Sub",
			Subcommands: []cli.Command{
				{
					Name:  "publish",
					Usage: "publish a document from file",
					Action: func(c *cli.Context) error {
						defer started(c, "publish")()
						args := ps.PubSubArguments{
							Project: c.String("p"),
							File:    c.String("f"),
							Topic:   c.String("t"),
						}
						return ps.Publish(args)
					},
					Flags: []cli.Flag{projectFlag, topicFlag, fileFlag},
				},
			},
		},
		{
			Name:  "iam",
			Usage: "Work with IAM",
		},
		{
			Name:  "bq",
			Usage: "Work with BigQuery",
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	return app
}
