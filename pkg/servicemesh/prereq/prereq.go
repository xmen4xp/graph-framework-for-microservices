package prereq

import (
	"bytes"
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

const goMinVersion = "1.17"

const (
	k8sMinVersion = "1.16"
	k8sMaxVersion = "1.21"
)

var All bool

type Prerequiste int

const (
	// Adding validation enum here to add multiple Prerequistes.
	DOCKER Prerequiste = iota
	KUBERNETES
	KUBERNETES_VERSION
	GOLANG_VERSION
	GOPATH
)

type PrerequisteMeta struct {
	what                  string
	verify                func() (bool, error)
	AdditionalDescription string
	Version               string
	Validator             Prerequiste
	Always                bool
}

var preReqs = map[Prerequiste]PrerequisteMeta{
	DOCKER: {
		what:                  "docker",
		Always:                true,
		AdditionalDescription: "docker daemon should be running on the host",
		verify: func() (bool, error) {
			_, err := exec.Command("docker", "ps").Output()
			if err != nil {
				return false, fmt.Errorf("verify if docker is running failed with error %v", err)
			}
			return true, nil
		},
	},
	GOLANG_VERSION: {
		what:    "go",
		Always:  true,
		Version: goMinVersion,
		verify: func() (bool, error) {
			out, err := exec.Command("go", "version").Output()
			if err != nil {
				return false, fmt.Errorf("verify go version failed with error %v", err)
			}
			re := regexp.MustCompile(`go[0-9][^\s]*`)
			match := re.FindStringSubmatch(string(out))
			if len(match) == 0 {
				return false, fmt.Errorf("unable to get go version from output: %v", string(out))
			}

			v1, errMinVersion := version.NewVersion(goMinVersion)
			if errMinVersion != nil {
				return false, fmt.Errorf("parse min go version failed with error %v", errMinVersion)
			}
			v2, errCurrVersion := version.NewVersion(strings.Trim(match[0], "go"))
			if errCurrVersion != nil {
				return false, fmt.Errorf("parse current go version failed with error %v", errCurrVersion)
			}

			if v2.LessThan(v1) {
				return false, fmt.Errorf("go version %s is less than %s", string(match[0]), goMinVersion)
			}

			return true, nil
		},
	},
	KUBERNETES: {
		what:                  "kubernetes",
		Always:                false,
		AdditionalDescription: "kubernetes cluster should be reachable via kubectl",
		verify: func() (bool, error) {
			_, err := exec.Command("kubectl", "get", "ns").Output()
			if err != nil {
				return false, fmt.Errorf("verifying running kubernetes cluster failed with error %v", err)
			}
			return true, nil
		},
	},
	KUBERNETES_VERSION: {
		what:                  "kubernetes version",
		Always:                false,
		AdditionalDescription: "Kubernetes version should be above 1.16 and below 1.20",
		verify: func() (bool, error) {
			versionStringBytes, _ := exec.Command("kubectl", "version", "--short=true").Output()
			cmd := exec.Command("tail", "-n", "1")
			cmd.Stdin = strings.NewReader(string(versionStringBytes))
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				return false, fmt.Errorf("could not get version string")
			}

			re := regexp.MustCompile(`Server Version: ([a-z0-9][^\s]*)`)
			serverVersion := re.FindStringSubmatch(out.String())
			if len(serverVersion) == 0 {
				return false, fmt.Errorf("unable to get k8s version from output: %v", out.String())
			}
			v1min, errMinVersion := version.NewVersion(k8sMinVersion)
			if errMinVersion != nil {
				return false, fmt.Errorf("verify version of kubernetes cluster is running failed with error on minVersion formation %v", errMinVersion)
			}
			v1max, errMaxVersion := version.NewVersion(k8sMaxVersion)
			if errMaxVersion != nil {
				return false, fmt.Errorf("verify version of kubernetes cluster is running failed with error maxVersion formation %v", errMaxVersion)
			}
			v1, errVersion := version.NewVersion(strings.TrimPrefix(serverVersion[1], "v"))
			if errVersion != nil {
				return false, fmt.Errorf("verify version of kubernetes cluster is running failed with error current version formation %v", errVersion)
			}
			if v1.LessThan(v1min) || v1.GreaterThanOrEqual(v1max) {
				return false, fmt.Errorf("K8s Version should be between 1.16 and 1.20, current Version is %s", v1)
			}
			return true, nil
		},
	},
	GOPATH: {
		what:                  "GOPATH",
		AdditionalDescription: "app workspace should be in GOPATH",
		verify: func() (bool, error) {

			gopath := os.Getenv("GOPATH")
			if gopath == "" {
				gopath = build.Default.GOPATH
			}

			workspacePath, err := os.Getwd()
			if err != nil {
				log.Println(err)
			}

			up := ".." + string(os.PathSeparator)
			// path-comparisons using filepath.Abs don't work reliably according to docs (no unique representation).
			if rel, err := filepath.Rel(gopath, workspacePath); err == nil {
				if !strings.HasPrefix(rel, up) && rel != ".." {
					return true, nil
				} else {
					return false, fmt.Errorf("workspace %s is not in GOPATH %s", workspacePath, gopath)
				}
			}
			fmt.Printf("Verifying if workspace is in GOPATH failed with error %s. Ignoring verification.", err)
			return true, nil
		},
	},
}

func PreReqVerify(cmd *cobra.Command, args []string) error {
	for _, util := range preReqs {
		if utils.VerifyAll(cmd) || util.Always {
			if ok, err := util.verify(); ok {
				fmt.Printf("\u2705 %s %s %s\n", util.what, util.Version, util.AdditionalDescription)
			} else {
				fmt.Printf("\u274C %s %s verify failed with err: %v\n", util.what, util.Version, err)
			}
		}
	}
	return nil
}

func printPreReq(req PrerequisteMeta) {
	fmt.Printf("\u2023 %s", req.what)
	if len(req.Version) > 0 {
		fmt.Printf(" (version: %s)", req.Version)
	}
	if len(req.AdditionalDescription) > 0 {
		fmt.Printf(" [ %s ]", req.AdditionalDescription)
	}
	fmt.Println()
}

func PreReqList(cmd *cobra.Command, args []string) error {
	for _, util := range preReqs {
		printPreReq(util)
	}
	return nil
}

func PreReqVerifyOnDemand(reqs []Prerequiste) error {
	for _, current := range reqs {
		util := preReqs[current]
		if ok, err := util.verify(); ok {
			return nil
		} else {
			return utils.GetCustomError(utils.APPLICATION_INIT_PREREQ_FAILED, err).Print().ExitIfFatalOrReturn()
		}
	}
	return nil
}

func PreReqListOnDemand(reqs []Prerequiste) {
	for _, current := range reqs {
		if util, found := preReqs[current]; found {
			printPreReq(util)
		}
	}
}

var PreReqVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "verify all pre-requisites",
	RunE:  PreReqVerify,
}

var PreReqListCmd = &cobra.Command{
	Use:   "list",
	Short: "list all pre-requisites",
	RunE:  PreReqList,
}

var PreReqCmd = &cobra.Command{
	Use:   "prereq",
	Short: "pre-requisites for a successful nexus-sdk experience",
}

func init() {
	PreReqCmd.AddCommand(PreReqListCmd)
	PreReqCmd.AddCommand(PreReqVerifyCmd)
	PreReqCmd.PersistentFlags().BoolVarP(&All, "all", "", false, "For validation check")
}
