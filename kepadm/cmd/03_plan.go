package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/calebamiles/keps/pkg/settings"
	"github.com/calebamiles/keps/pkg/workflow"
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "beginning planning the implementation of a KEP",
	Long: `
Plan helps KEP authors to prepare for implementation by creating
places to draft documentation targeted at:
  - developers (similar to the more traditional "design proposal")
  - teachers, such as Technical Writers or Trainers
  - operators, for software based enhancements

We believe that KEP authors should begin to work as early as possible with the cross
functional team responsible for the sucessful landing of an enhancement. It is not
expected that the documentation produced during the plan phase will be publication
quality, however, a good faith effort to prepare draft quality documentation is
expected and will be checked by the KEP reviewers and approvers.`,
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
		err = workflow.Plan(runtimeSettings)
		if err != nil {
			return err
		}

		fmt.Println("successfully started planning for KEP!")
		return nil
	},
}
