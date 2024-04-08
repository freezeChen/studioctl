package main

import (
	"fmt"
	"github.com/freezeChen/studioctl/core/genCode"
	"github.com/urfave/cli"
	"os"
	"runtime"
)

var (
	buildVersion = "0.0.1"
	commands     = []cli.Command{
		{
			Name:  "model",
			Usage: "generate model code",
			Subcommands: []cli.Command{
				{
					Name:  "mysql",
					Usage: "generate mysql model",
					Subcommands: []cli.Command{
						{
							Name:  "datasource",
							Usage: "generate model from datasource",
							Flags: []cli.Flag{
								cli.StringFlag{
									Name:     "url",
									Usage:    `the data source of database, like "root:password@tcp(127.0.0.1:3306)/database"`,
									Required: true,
								},
								cli.StringFlag{
									Name:     "table,t",
									Usage:    "the table or table globbing patterns in the database",
									Required: true,
								},
								cli.StringFlag{
									Name:     "dir,d",
									Usage:    "the target dir",
									Required: true,
								},
								cli.StringFlag{
									Name:  "style,s",
									Usage: "the table columns naming format like snake(default),same",
								},
								cli.StringFlag{
									Name:  "tag",
									Usage: "tag style default(xorm),xorm,gorm",
									Value: "xorm",
								},
								cli.BoolTFlag{
									Name:  "curd,c",
									Usage: "gen simple curd function with gen model",
								},
								cli.BoolFlag{
									Name:  "proto,p",
									Usage: "print protobuf model",
								},
							},

							Action: genCode.CodeHandler,
						},
					},
				},
			},
		},
		{
			Name: "web",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  genCode.FlagUrl,
					Usage: `the data source of database, like "root:password@tcp(127.0.0.1:3306)/database"`,
				},
				cli.StringFlag{
					Name:  genCode.FlagPort,
					Usage: `web port`,
					Value: ":8888",
				},
				cli.StringFlag{Name: genCode.FlagType,
					Usage: `database type  like mysql,mssql`,
					Value: "mysql",
				},
			},
			Action: genCode.CodeHandler,
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Version = fmt.Sprintf("%s %s/%s", buildVersion, runtime.GOOS, runtime.GOARCH)
	app.Commands = commands
	if err := app.Run(os.Args); err != nil {
		fmt.Println("error:", err)
		return
	}

}
