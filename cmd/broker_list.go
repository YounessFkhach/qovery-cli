package cmd

import (
	"fmt"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
	"os"
	"qovery.go/api"
	"qovery.go/util"
	"strings"
)

var brokerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List brokers",
	Long: `LIST show all available brokers within a project and environment. For example:

	qovery broker list`,
	Run: func(cmd *cobra.Command, args []string) {
		if !hasFlagChanged(cmd) {
			BranchName = util.CurrentBranchName()
			qoveryYML, err := util.CurrentQoveryYML()
			if err != nil {
				util.PrintError("No qovery configuration file found")
				os.Exit(1)
			}
			ProjectName = qoveryYML.Application.Project
		}

		ShowBrokerList(ProjectName, BranchName)
	},
}

func init() {
	brokerListCmd.PersistentFlags().StringVarP(&ProjectName, "project", "p", "", "Your project name")
	brokerListCmd.PersistentFlags().StringVarP(&BranchName, "branch", "b", "", "Your branch name")

	brokerCmd.AddCommand(brokerListCmd)
}

func ShowBrokerList(projectName string, branchName string) {
	output := []string{
		"name | status | type | version | endpoint | port | username | password | application",
	}

	projectId := api.GetProjectByName(projectName).Id
	environment := api.GetEnvironmentByName(projectId, branchName)

	services := api.ListBrokers(projectId, environment.Id)

	if services.Results == nil || len(services.Results) == 0 {
		fmt.Println(columnize.SimpleFormat(output))
		return
	}

	for _, a := range services.Results {
		applicationName := "none"

		if a.Applications != nil {
			applicationName = strings.Join(a.GetApplicationNames(), ", ")
		}

		output = append(output, strings.Join([]string{
			a.Name,
			a.Status.CodeMessage,
			a.Type,
			a.Version,
			a.FQDN,
			intPointerValue(a.Port),
			a.Username,
			a.Password,
			applicationName,
		}, " | "))
	}

	fmt.Println(columnize.SimpleFormat(output))
}
