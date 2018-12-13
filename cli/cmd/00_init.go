package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/calebamiles/keps/pkg/settings"
	"github.com/calebamiles/keps/pkg/workflow"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "create a new KEP directory at location",
	Long: `
Creates a KEP at the given location. The KEP title is inferred from the
last path element. For example:

- Kubernetes Wide KEP:
	kep init content/kubernetes-enhancement-proposal-process

  will create a KEP titled "Kubernetes Enhancement Proposal Process" at
  content/kubernetes-wide/kubernetes-enhancement-proposal-process

- SIG Wide KEP:
	kep init content/api-machinery/server-side-apply

  will create a KEP titled "Server Side Apply" at
  content/api-machinery/sig-wide/server-side-apply

- Subproject Level KEP:
	kep init content/sig-node/kubelet/dynamic-kubelet-config

  will create a KEP titled "Dynamic Kubelet Config" at
  content/sig-node/kubelet/dynamic-kubelet-config

init will create a KEP directory containing templated files for the required KEP sections.`,
	Args: cobra.ExactArgs(1), // accept just one argument, target location for content
	RunE: func(cmd *cobra.Command, args []string) error {
		targetPath := args[0] // we have a validator ensuring we will have exactly one positional arg

		contentRoot, err := settings.FindContentRoot()
		if err != nil {
			return err
		}

		// save it now to avoid the expensive look everywhere under $HOME next time
		err = settings.SaveContentRoot(contentRoot)
		if err != nil {
			return err
		}

		principal, err := settings.FindPrincipal()
		if err != nil {
			return err
		}

		runtimeSettings := settings.NewRuntime(contentRoot, targetPath, principal)
		kepContentDir, err := workflow.Init(runtimeSettings)
		if err != nil {
			return err
		}

		fmt.Printf("created KEP template at: %s\n", kepContentDir)
		return nil
	},
}

// TODO maybe kill off these init() functions
// the best worst place for them might be inside of cmd.Execute() which lives inside of root.go
func init() {
	// rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
