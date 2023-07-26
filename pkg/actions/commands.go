package actions

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

var (
	outfile, _ = os.Create("issueUpdater.log")
	l          = log.New(outfile, "", log.LstdFlags|log.Lshortfile)
)

func Commands() {
	app := &cli.App{
		Flags: []cli.Flag{
			cli.StringFlag{Name: "uid", Usage: "User ID"},
			cli.StringFlag{Name: "api_key", Usage: "API Key"},
			cli.StringFlag{Name: "scan_id", Usage: "[Optional] Scan ID to update issues. "},
			cli.StringFlag{Name: "issue_name", Usage: "Name of the issue to be updated"},
			cli.StringFlag{Name: "update", Usage: "Update type. USAGE: --update false_positive. Available options: {false_positive, accepted_risk, fixed_unconfirmed, fixed_cant_retest} Note: You should use proper 'fixed' according to the issue type. (Retestable or not)"},
		},

		Action: func(c *cli.Context) error {
			if c.String("uid") != "" && c.String("api_key") != "" && c.String("issue_name") != "" {
				// Get actions from user
				uid := c.String("uid")
				api_key := c.String("api_key")
				scan_id := c.String("scan_id")
				issue_name := c.String("issue_name")
				update_type := c.String("update")

				IssueActions(uid, api_key, scan_id, issue_name, update_type)
			} else {
				fmt.Println("The API Key, UID and issue_name parameters are necessary!")
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
