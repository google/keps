package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/calebamiles/keps/pkg/settings"
	"github.com/calebamiles/keps/pkg/workflow"
)

// approveCmd represents the approve command
var approveCmd = &cobra.Command{
	Use:   "approve",
	Short: "approve the KEP design documentation; and mark a KEP as ready for implementation",
	Long: `
Approve the implementation approach for a KEP on behalf of a SIG.

Approve communicates that the SIG trusts the KEP authors to begin the implementation
of the described enhancement. At this point Kubernetes project resources such as a
mailing list, or git repository should be allocated to assist the KEP authors towards
implementation.`,
	Args: cobra.ExactArgs(1), // accept just one argument, location of KEP
	RunE: func(cmd *cobra.Command, args []string) error {
		targetPath := args[0] // we have a validator ensuring we will have exactly one positional

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
		err = workflow.Approve(runtimeSettings)
		if err != nil {
			return err
		}

		fmt.Println("sucessfully marked KEP as approved!")
		return nil
	},
}
