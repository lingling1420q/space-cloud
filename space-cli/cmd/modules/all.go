package modules

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/modules/auth"
	"github.com/spaceuptech/space-cli/cmd/modules/database"
	"github.com/spaceuptech/space-cli/cmd/modules/eventing"
	"github.com/spaceuptech/space-cli/cmd/modules/filestore"
	"github.com/spaceuptech/space-cli/cmd/modules/ingress"
	"github.com/spaceuptech/space-cli/cmd/modules/letsencrypt"
	"github.com/spaceuptech/space-cli/cmd/modules/project"
	remoteservices "github.com/spaceuptech/space-cli/cmd/modules/remote-services"
	"github.com/spaceuptech/space-cli/cmd/modules/services"
	"github.com/spaceuptech/space-cli/cmd/utils"
)

func getSubCommands() []*cobra.Command {

	var getProjects = &cobra.Command{
		Use:   "all",
		Short: "Gets entire project config",
		RunE:  getAllProjects,
	}

	return []*cobra.Command{getProjects}
}

func getAllProjects(cmd *cobra.Command, args []string) error {
	projectName := viper.GetString("project")

	if len(args) == 0 {
		_ = utils.LogError("Directory not specified as an arguement to store config files", nil)
		return nil
	}
	dir := args[0]
	// create directory if directory doesn't exists
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil
	}

	if err := os.Chdir(dir); err != nil {
		return nil
	}

	obj, err := project.GetProjectConfig(projectName, "project", map[string]string{})
	if err != nil {
		return nil
	}
	if err := createConfigFile("1", "project", obj); err != nil {
		return nil
	}

	objs, err := database.GetDbConfig(projectName, "db-config", map[string]string{})
	if err != nil {
		return nil
	}
	if err := createConfigFile("2", "db-config", objs); err != nil {
		return nil
	}

	objs, err = database.GetDbRule(projectName, "db-rules", map[string]string{})
	if err != nil {
		return nil
	}
	if err := createConfigFile("3", "db-rules", objs); err != nil {
		return nil
	}

	objs, err = database.GetDbSchema(projectName, "db-schema", map[string]string{})
	if err != nil {
		return nil
	}
	if err := createConfigFile("4", "db-schema", objs); err != nil {
		return nil
	}

	obj, err = filestore.GetFileStoreConfig(projectName, "filestore-config", map[string]string{})
	if err != nil {
		return nil
	}
	if err := createConfigFile("5", "filestore-config", obj); err != nil {
		return nil
	}

	objs, err = filestore.GetFileStoreRule(projectName, "filestore-rule", map[string]string{})
	if err != nil {
		return nil
	}
	if err := createConfigFile("6", "filestore-rule", objs); err != nil {
		return nil
	}

	obj, err = eventing.GetEventingConfig(projectName, "eventing-config", map[string]string{})
	if err != nil {
		return nil
	}
	if err := createConfigFile("7", "eventing-config", obj); err != nil {
		return nil
	}

	objs, err = eventing.GetEventingTrigger(projectName, "eventing-triggers", map[string]string{})
	if err != nil {
		return nil
	}
	if err := createConfigFile("8", "eventing-triggers", objs); err != nil {
		return nil
	}

	objs, err = eventing.GetEventingSecurityRule(projectName, "eventing-rule", map[string]string{})
	if err != nil {
		return nil
	}
	if err := createConfigFile("9", "eventing-rule", objs); err != nil {
		return nil
	}

	objs, err = eventing.GetEventingSchema(projectName, "eventing-schema", map[string]string{})
	if err != nil {
		return nil
	}
	if err := createConfigFile("10", "eventing-schema", objs); err != nil {
		return nil
	}

	objs, err = remoteservices.GetRemoteServices(projectName, "remote-services", map[string]string{})
	if err != nil {
		return nil
	}
	if err := createConfigFile("11", "remote-services", objs); err != nil {
		return nil
	}

	objs, err = services.GetServices(projectName, "service", map[string]string{})
	if err != nil {
		return nil
	}
	if err := createConfigFile("12", "services", objs); err != nil {
		return nil
	}

	objs, err = services.GetServicesRoutes(projectName, "services-routes", map[string]string{})
	if err != nil {
		return nil
	}
	if err := createConfigFile("13", "services-routes", objs); err != nil {
		return nil
	}

	// objs, _ = services.GetServicesSecrets(projectName, "services-secrets", map[string]string{})
	// if _ != nil {
	// 	return _
	// }
	// _ = createConfigFile("14", "services-secrets", objs); _ != nil {
	// 	return _
	// }

	objs, err = ingress.GetIngressRoutes(projectName, "ingress-routes", map[string]string{})
	if err != nil {
		return nil
	}
	if err := createConfigFile("15", "ingress-routes", objs); err != nil {
		return nil
	}

	objs, err = auth.GetAuthProviders(projectName, "auth-providers", map[string]string{})
	if err != nil {
		return nil
	}
	if err := createConfigFile("16", "auth-providers", objs); err != nil {
		return nil
	}

	obj, err = letsencrypt.GetLetsEncryptDomain(projectName, "letsencrypt", map[string]string{})
	if err != nil {
		return nil
	}
	if err := createConfigFile("17", "letsencrypt", obj); err != nil {
		return nil
	}

	return nil
}

func createConfigFile(pos, commandName string, objs []*model.SpecObject) error {
	message := ""
	for _, val := range objs {
		data, err := yaml.Marshal(val)
		if err != nil {
			return err
		}
		message = message + string(data) + "---" + "\n"
	}
	fileName := fmt.Sprintf("%s-%s.yaml", pos, commandName)
	if err := ioutil.WriteFile(fileName, []byte(message), 0755); err != nil {
		return err
	}
	return nil
}
