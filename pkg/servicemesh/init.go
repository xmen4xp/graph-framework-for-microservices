package servicemesh

import (
	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/lib-go/logging"
	"gitlab.eng.vmware.com/nexus/cli/pkg/servicemesh/app"
	"gitlab.eng.vmware.com/nexus/cli/pkg/servicemesh/apply"
	"gitlab.eng.vmware.com/nexus/cli/pkg/servicemesh/cluster"
	"gitlab.eng.vmware.com/nexus/cli/pkg/servicemesh/config"
	servicemesh_datamodel "gitlab.eng.vmware.com/nexus/cli/pkg/servicemesh/datamodel"
	"gitlab.eng.vmware.com/nexus/cli/pkg/servicemesh/gns"
	"gitlab.eng.vmware.com/nexus/cli/pkg/servicemesh/login"
	"gitlab.eng.vmware.com/nexus/cli/pkg/servicemesh/runtime"
)

// ClusterCmd ... Cluster command
var ClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Servicemesh cluster features",
}

// GnsCmd ... GNS command
var GnsCmd = &cobra.Command{
	Use:   "gns",
	Short: "Servicemesh global namespace features",
}

// ConfigCmd ... Config command
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Servicemesh configuration features",
}

// ApplyCmd ... Apply command
var ApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply Servicemesh configuration from file",
	RunE:  apply.ApplyResource,
}

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Servicemesh configuration from file",
	RunE:  apply.DeleteResource,
}

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to csp",
	RunE:  login.Login,
}

var RuntimeCmd = &cobra.Command{
	Use:   "runtime",
	Short: "Runtime installer and uninstaller",
}

var DataModelCmd = &cobra.Command{
	Use:   "datamodel",
	Short: "Datamodel installer and uninstaller",
}

var AppCmd = &cobra.Command{
	Use:   "app",
	Short: "Sample application installer",
}

func initCommands() {
	// Cluster commands
	ClusterCmd.AddCommand(cluster.GetCmd)
	ClusterCmd.AddCommand(cluster.ListCmd)
	ClusterCmd.AddCommand(cluster.CreateCmd)
	ClusterCmd.AddCommand(cluster.DeleteCmd)

	// GNS commands
	GnsCmd.AddCommand(gns.GetCmd)
	GnsCmd.AddCommand(gns.ListCmd)
	GnsCmd.AddCommand(gns.CreateCmd)
	GnsCmd.AddCommand(gns.DeleteCmd)
	// Config commands
	ConfigCmd.AddCommand(config.ViewCmd)

	ConfigCmd.AddCommand(config.ViewCmd)

	ApplyCmd.Flags().StringVarP(&apply.CreateResourceFile, "file",
		"f", "", "Resource file from which cluster is created.")

	err := cobra.MarkFlagRequired(ApplyCmd.Flags(), "file")
	if err != nil {
		logging.Debugf("init error: %v", err)
	}

	DeleteCmd.Flags().StringVarP(&apply.DeleteResourceFile, "file",
		"f", "", "Resource file from which cluster is created.")

	err = cobra.MarkFlagRequired(DeleteCmd.Flags(), "file")
	if err != nil {
		logging.Debugf("init error: %v", err)
	}

	LoginCmd.Flags().StringVarP(&login.ApiToken, "token",
		"t", "", "token for api access")

	err = cobra.MarkFlagRequired(LoginCmd.Flags(), "token")
	if err != nil {
		logging.Debugf("api token is mandatory for login")
	}

	LoginCmd.Flags().StringVarP(&login.Server, "server",
		"s", "", "saas server fqdn")

	err = cobra.MarkFlagRequired(LoginCmd.Flags(), "server")
	if err != nil {
		logging.Debugf("saas server fqdn name is mandatory for login")
	}
	RuntimeCmd.AddCommand(runtime.InstallCmd)
	RuntimeCmd.AddCommand(runtime.UninstallCmd)

	DataModelCmd.AddCommand(servicemesh_datamodel.InitCmd)
	DataModelCmd.AddCommand(servicemesh_datamodel.InstallCmd)
	DataModelCmd.AddCommand(servicemesh_datamodel.BuildCmd)

	AppCmd.AddCommand(app.InitCmd)
	AppCmd.AddCommand(app.PackageCmd)
	AppCmd.AddCommand(app.PublishCmd)
	AppCmd.AddCommand(app.DeployCmd)
	AppCmd.AddCommand(app.RunCmd)
}

func init() {
	initCommands()
}
