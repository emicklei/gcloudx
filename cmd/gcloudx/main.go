package main

import (
	"log"
	"os"
	"sort"

	"github.com/emicklei/gcloudx/bq"
	"github.com/emicklei/gcloudx/im"
	"github.com/emicklei/gcloudx/ps"
	"github.com/urfave/cli/v2"
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
	cli.VersionFlag = &cli.BoolFlag{
		Name:  "print-version, V",
		Usage: "print only the version",
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "v",
			Usage: "verbose logging",
		},
		&cli.BoolFlag{
			Name:  "q",
			Usage: "quiet mode, accept any prompt",
		},
	}
	projectFlag := &cli.StringFlag{
		Name:  "p",
		Usage: `GCP project identifier`,
	}
	topicFlag := &cli.StringFlag{
		Name:  "t",
		Usage: `PubSub topic identifier (short name)`,
	}
	subscriptionFlag := &cli.StringFlag{
		Name:  "s",
		Usage: `PubSub subscription identifier (short name)`,
	}
	fileFlag := &cli.StringFlag{
		Name:  "f",
		Usage: `File containing the payload`,
	}
	bqDotOutputFlag := &cli.StringFlag{
		Name:  "o",
		Usage: `output file with DOT notation`,
		Value: "bigquery.dot",
	}
	app.Commands = []*cli.Command{
		{
			Name:  "pubsub",
			Usage: "Work with Pub/Sub",
			Subcommands: []*cli.Command{
				{
					Name:  "publish",
					Usage: "publish a document from file",
					Action: func(c *cli.Context) error {
						defer logBegin(c)()
						args := ps.PubSubArguments{
							Project: c.String("p"),
							File:    c.String("f"),
							Topic:   c.String("t"),
						}
						return ps.Publish(args)
					},
					Flags: []cli.Flag{projectFlag, topicFlag, fileFlag},
				},
				{
					Name:  "create-topic",
					Usage: "create a new topic",
					Action: func(c *cli.Context) error {
						defer logBegin(c)()
						args := ps.PubSubArguments{
							Project: c.String("p"),
							Topic:   c.String("t"),
						}
						return ps.CreateTopic(args)
					},
					Flags: []cli.Flag{projectFlag, topicFlag},
				},
				{
					Name:  "create-subscription",
					Usage: "create a new subscription",
					Action: func(c *cli.Context) error {
						defer logBegin(c)()
						args := ps.PubSubArguments{
							Project:      c.String("p"),
							Topic:        c.String("t"),
							Subscription: c.String("s"),
						}
						return ps.CreateSubscription(args)
					},
					Flags: []cli.Flag{projectFlag, topicFlag, subscriptionFlag},
				},
			},
		},
		{
			Name:  "iam",
			Usage: "Work with IAM",
			Subcommands: []*cli.Command{
				{
					Name:  "roles",
					Usage: "list all permissions assigned to a member",
					Action: func(c *cli.Context) error {
						defer logBegin(c)()
						args := im.IAMArguments{
							Verbose: c.Bool("v"),
							Member:  c.Args().First(),
						}
						return im.Roles(args)
					},
				},
				{
					Name:  "owners",
					Usage: "list all members thata have Project Owner permission on a project",
					Action: func(c *cli.Context) error {
						defer logBegin(c)()
						args := im.IAMArguments{
							Verbose: c.Bool("v"),
							Member:  c.Args().First(),
						}
						return im.Owners(args)
					},
				},
			},
		},
		{
			Name:  "bq",
			Usage: "Work with BigQuery",
			Subcommands: []*cli.Command{
				{
					Name:  "deps",
					Usage: "bq deps PROJECT(.|:)DATASET.VIEW,...",
					Action: func(c *cli.Context) error {
						defer logBegin(c)()
						args := bq.BigQueryArguments{
							Verbose: c.Bool("v"),
							Output:  c.String("o"),
						}
						for i := 0; i < c.NArg(); i++ {
							args.TableSources = append(args.TableSources, c.Args().Get(i))
						}
						return bq.ExportViewDepencyGraph(args)
					},
					Flags: []cli.Flag{bqDotOutputFlag},
				},
			},
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	return app
}
