package app

import (
	"github.com/spf13/cobra"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

func Package(cmd *cobra.Command, args []string) error {
	envList := []string{}
	err := utils.SystemCommand(envList, false, "make", "app_package")
	if err != nil {
		return err
	}
	return nil
}

var PackageCmd = &cobra.Command{
	Use:   "package",
	Short: "installs application package",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: Package,
}
