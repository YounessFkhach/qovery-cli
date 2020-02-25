package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"qovery.go/api"
	"qovery.go/util"
)

var deployListCmd = &cobra.Command{
	Use:   "list",
	Short: "List deployments",
	Long: `LIST show all deployable environment. For example:

	qovery deploy list`,
	Run: func(cmd *cobra.Command, args []string) {
		if !hasFlagChanged(cmd) {
			BranchName = util.CurrentBranchName()
			qoveryYML := util.CurrentQoveryYML()
			ProjectName = qoveryYML.Application.Project
			ApplicationName = qoveryYML.Application.Name

			if BranchName == "" || ProjectName == "" || ApplicationName == "" {
				fmt.Println("The current directory is not a Qovery project (-h for help)")
				os.Exit(0)
			}
		}

		ShowDeploymentList(ProjectName, BranchName, ApplicationName)
	},
}

func init() {
	deployListCmd.PersistentFlags().StringVarP(&ApplicationName, "application", "a", "", "Your application name")
	deployListCmd.PersistentFlags().StringVarP(&ProjectName, "project", "p", "", "Your project name")
	deployListCmd.PersistentFlags().StringVarP(&BranchName, "branch", "b", "", "Your branch name")

	deployCmd.AddCommand(deployListCmd)
}

func ShowDeploymentList(projectName string, branchName string, applicationName string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"branch", "commit date", "commit id", "commit author", "deployed"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorders(tablewriter.Border{Left: false, Top: true, Right: false, Bottom: true})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.BgMagentaColor, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.BgMagentaColor, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.BgMagentaColor, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.BgMagentaColor, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.BgMagentaColor, tablewriter.FgWhiteColor})

	environments := api.GetBranchByName(api.GetProjectByName(projectName).Id, branchName).Environments
	if len(environments) == 0 {
		table.Append([]string{"", "", "", "", ""})
		table.Render()
		return
	}

	var environment api.Environment
	for _, e := range environments {
		if e.Application.Name == applicationName {
			environment = e
		}
	}

	if environment.Id == "" {
		table.Append([]string{"", "", "", ""})
		table.Render()
		return
	}

	for _, commit := range util.ListCommits(10) {
		if environment.CommitId == commit.ID().String() {
			table.Append([]string{branchName, commit.Author.When.String(), commit.ID().String(), commit.Author.Name, "✓"})
		} else {
			table.Append([]string{branchName, commit.Author.When.String(), commit.ID().String(), commit.Author.Name, ""})
		}
	}
	table.Render()
}
