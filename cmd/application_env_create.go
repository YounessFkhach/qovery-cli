package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/pterm/pterm"
	"github.com/qovery/qovery-cli/utils"
	"github.com/spf13/cobra"
)

var applicationEnvCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create application environment variable or secret",
	Run: func(cmd *cobra.Command, args []string) {
		utils.Capture(cmd)

		tokenType, token, err := utils.GetAccessToken()
		if err != nil {
			utils.PrintlnError(err)
			os.Exit(1)
			panic("unreachable") // staticcheck false positive: https://staticcheck.io/docs/checks#SA5011
		}

		client := utils.GetQoveryClient(tokenType, token)
		_, projectId, envId, err := getContextResourcesId(client)

		if err != nil {
			utils.PrintlnError(err)
			os.Exit(1)
			panic("unreachable") // staticcheck false positive: https://staticcheck.io/docs/checks#SA5011
		}

		applications, _, err := client.ApplicationsApi.ListApplication(context.Background(), envId).Execute()

		if err != nil {
			utils.PrintlnError(err)
			os.Exit(1)
			panic("unreachable") // staticcheck false positive: https://staticcheck.io/docs/checks#SA5011
		}

		application := utils.FindByApplicationName(applications.GetResults(), applicationName)

		if application == nil {
			utils.PrintlnError(fmt.Errorf("application %s not found", applicationName))
			utils.PrintlnInfo("You can list all applications with: qovery application list")
			os.Exit(1)
			panic("unreachable") // staticcheck false positive: https://staticcheck.io/docs/checks#SA5011
		}

		if utils.IsSecret {
			err = utils.CreateSecret(client, projectId, envId, application.Id, utils.Key, utils.Value, utils.ApplicationScope)

			if err != nil {
				utils.PrintlnError(err)
				os.Exit(1)
				panic("unreachable") // staticcheck false positive: https://staticcheck.io/docs/checks#SA5011
			}

			utils.Println(fmt.Sprintf("Secret %s has been created", pterm.FgBlue.Sprintf(utils.Key)))
			return
		}

		err = utils.CreateEnvironmentVariable(client, projectId, envId, application.Id, utils.Key, utils.Value, utils.ApplicationScope)

		if err != nil {
			utils.PrintlnError(err)
			os.Exit(1)
			panic("unreachable") // staticcheck false positive: https://staticcheck.io/docs/checks#SA5011
		}

		utils.Println(fmt.Sprintf("Environment variable %s has been created", pterm.FgBlue.Sprintf(utils.Key)))
	},
}

func init() {
	applicationEnvCmd.AddCommand(applicationEnvCreateCmd)
	applicationEnvCreateCmd.Flags().StringVarP(&organizationName, "organization", "", "", "Organization Name")
	applicationEnvCreateCmd.Flags().StringVarP(&projectName, "project", "", "", "Project Name")
	applicationEnvCreateCmd.Flags().StringVarP(&environmentName, "environment", "", "", "Environment Name")
	applicationEnvCreateCmd.Flags().StringVarP(&applicationName, "application", "n", "", "Application Name")
	applicationEnvCreateCmd.Flags().StringVarP(&utils.Key, "key", "k", "", "Environment variable or secret key")
	applicationEnvCreateCmd.Flags().StringVarP(&utils.Value, "value", "v", "", "Environment variable or secret value")
	applicationEnvCreateCmd.Flags().StringVarP(&utils.ApplicationScope, "scope", "", "APPLICATION", "Scope of this env var <PROJECT|ENVIRONMENT|APPLICATION>")
	applicationEnvCreateCmd.Flags().BoolVarP(&utils.IsSecret, "secret", "", false, "This environment variable is a secret")

	_ = applicationEnvCreateCmd.MarkFlagRequired("key")
	_ = applicationEnvCreateCmd.MarkFlagRequired("value")
	_ = applicationEnvCreateCmd.MarkFlagRequired("application")
}
