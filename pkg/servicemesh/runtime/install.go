package runtime

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/prereq"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/logging"
)

var Namespace string
var Filename = "runtime-manifests.tar"
var ManifestsDir = "manifets-nexus-runtime"

var prerequisites []prereq.Prerequiste = []prereq.Prerequiste{
	prereq.KUBERNETES,
	prereq.KUBERNETES_VERSION,
}

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs the Nexus runtime on the specified namespace",
	//Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return prereq.PreReqVerifyOnDemand(prerequisites)
	},
	RunE: Install,
}

func Install(cmd *cobra.Command, args []string) error {
	files, err := DownloadRuntimeFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			utils.SystemCommand(cmd, utils.RUNTIME_INSTALL_FAILED, []string{}, "kubectl", "apply", "-f", filepath.Join(ManifestsDir, "runtime-manifests", file.Name()), "-n", Namespace)
		}
	}
	for _, label := range common.PodLabels {
		utils.CheckPodRunning(cmd, utils.RUNTIME_INSTALL_FAILED, label, Namespace)
	}
	fmt.Printf("\u2713 Runtime installation successful on namespace %s\n", Namespace)
	os.Remove(Filename)
	os.RemoveAll(ManifestsDir)
	return nil
}

func init() {
	InstallCmd.Flags().StringVarP(&Namespace, "namespace",
		"n", "", "name of the namespace to be created")
	err := cobra.MarkFlagRequired(InstallCmd.Flags(), "namespace")
	if err != nil {
		logging.Debugf("Runtime install err: %v", err)
	}
}

func DownloadRuntimeFiles() ([]fs.FileInfo, error) {
	err := utils.DownloadFile(common.RUNTIME_MANIFESTS_URL, Filename)
	if err != nil {
		return nil, utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("could not download the runtime manifests due to %s", err)).Print().ExitIfFatalOrReturn()
	}
	file, err := os.Open(Filename)
	if err != nil {
		return nil, utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("accessing downloaded runtime manifests directory failed dwith error: %s", err)).Print().ExitIfFatalOrReturn()
	}
	defer file.Close()
	_, err = os.Stat(ManifestsDir)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(ManifestsDir, os.ModePerm)
		if err != nil {
			return nil, utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
				fmt.Errorf("could not create the runtime manifests directory due to %s", err)).Print().ExitIfFatalOrReturn()
		}
	}
	err = utils.Untar(ManifestsDir, file)
	if err != nil {
		return nil, utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("unarchive of runtime manifests directory failed with error %s", err)).Print().ExitIfFatalOrReturn()
	}

	files, err := ioutil.ReadDir(filepath.Join(ManifestsDir, "runtime-manifests"))
	if err != nil {
		return nil, utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("accessing runtime manifests directory failed with error to %s", err)).Print().ExitIfFatalOrReturn()
	}
	return files, nil
}
