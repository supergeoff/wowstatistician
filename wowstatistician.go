package main

import (
	"errors"
	"log"
	"os"
	"wowstatistician/cmd"
	"wowstatistician/controllers"
	"wowstatistician/helpers/databases"
	_ "wowstatistician/routers"

	"github.com/astaxie/beego"
	"github.com/davecgh/go-spew/spew"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Wow Statistician",
		Usage: "Cli tool for wow class repartition stats",
		Commands: []*cli.Command{
			{
				Name:    "compute",
				Aliases: []string{"c"},
				Usage:   "Compute stats",
				Subcommands: []*cli.Command{
					{
						Name:    "generate",
						Aliases: []string{"g"},
						Usage:   "Generate stats for a db",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "database",
								Aliases:  []string{"db"},
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							log.Println("[+] Generating stats for db: databases/" + c.String("database"))
							err := databases.WriteStatsForDb(c.String("database"))
							if err != nil {
								return err
							}
							log.Println("[-] Generating stats for db: databases/" + c.String("database"))
							return nil
						},
					},
					{
						Name:    "print",
						Aliases: []string{"p"},
						Usage:   "Print stats maps object on console",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "database",
								Aliases:  []string{"db"},
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							log.Println("[+] Printing stats for db: databases/" + c.String("database"))
							stats, err := databases.ReadStatsDb(c.String("database"))
							if err != nil {
								return err
							}
							spew.Dump(stats)
							log.Println("[-] Printing stats for db: databases/" + c.String("database"))
							return nil
						},
					},
				},
			},
			{
				Name:    "get",
				Aliases: []string{"g"},
				Usage:   "Retreive data from blizzard api",
				Subcommands: []*cli.Command{
					{
						Name:    "arena",
						Aliases: []string{"a"},
						Usage:   "Get and store arena leatherboards",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "region",
								Aliases: []string{"rg"},
								Value:   "eu",
								Usage:   "Region to query arena leatherboards from",
							},
						},
						Action: func(c *cli.Context) error {
							log.Println("[+] Saving arena profiles")
							err := cmd.SaveArenaProfiles(c.String("region"))
							if err != nil {
								return err
							}
							log.Println("[-] Saving arena profiles")
							return nil
						},
					},
					{
						Name:    "mythic",
						Aliases: []string{"m"},
						Usage:   "Get and store mythic+ leatherboards",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "region",
								Aliases: []string{"rg"},
								Value:   "eu",
								Usage:   "Region to query mythic+ leatherboards from",
							},
						},
						Action: func(c *cli.Context) error {
							log.Println("[+] Saving mythic+ profiles")
							err := cmd.SaveMythicProfiles(c.String("region"))
							if err != nil {
								return err
							}
							log.Println("[-] Saving mythic+ profiles")
							return nil
						},
					},
					{
						Name:    "raid",
						Aliases: []string{"r"},
						Usage:   "Get and store raid leatherboard",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "region",
								Aliases: []string{"rg"},
								Value:   "eu",
								Usage:   "Region to query raid leatherboard from",
							},
							&cli.StringFlag{
								Name:    "raid",
								Aliases: []string{"r"},
								Value:   "nyalotha-the-waking-city",
								Usage:   "Raid to query leatherboard from",
							},
						},
						Action: func(c *cli.Context) error {
							log.Println("[+] Saving raid profiles")
							err := cmd.SaveRaidProfiles(c.String("region"), c.String("raid"))
							if err != nil {
								return err
							}
							log.Println("[-] Saving raid profiles")
							return nil
						},
					},
					{
						Name:    "rbg",
						Aliases: []string{"rb"},
						Usage:   "Get and store rbg leatherboard",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "region",
								Aliases: []string{"rg"},
								Value:   "eu",
								Usage:   "Region to query rbg leatherboards from",
							},
						},
						Action: func(c *cli.Context) error {
							log.Println("[+] Saving rbg profiles")
							err := cmd.SaveRbgProfiles(c.String("region"))
							if err != nil {
								return err
							}
							log.Println("[-] Saving rbg profiles")
							return nil
						},
					},
				},
			},
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Usage:   "Serve results as html",
				Action: func(c *cli.Context) error {
					db, err := databases.OpenDB("databases/stats")
					controllers.Db = db
					if err != nil {
						return errors.New("main: could not open stats db - " + err.Error())
					}
					defer db.Close()
					beego.Run()
					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
