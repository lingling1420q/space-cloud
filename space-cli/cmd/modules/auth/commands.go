package auth

import (
	"github.com/spaceuptech/space-cli/cmd/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetSubCommands is the list of commands the auth module exposes
func GetSubCommands() []*cobra.Command {
	var getAuthProviders = &cobra.Command{
		Use:  "auth-providers",
		RunE: actionGetAuthProviders,
	}
	return []*cobra.Command{getAuthProviders}
}

func actionGetAuthProviders(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project := viper.GetString("project")
	commandName := cmd.Use

	params := map[string]string{}
	if len(args) != 0 {
		params["id"] = args[0]
	}

	objs, err := GetAuthProviders(project, commandName, params)
	if err != nil {
		return nil
	}

	if err := utils.PrintYaml(objs); err != nil {
		return nil
	}
	return nil
}
