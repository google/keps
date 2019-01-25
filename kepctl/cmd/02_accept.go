package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/calebamiles/keps/pkg/settings"
	"github.com/calebamiles/keps/pkg/workflow"
)

// acceptCmd represents the accept command
var acceptCmd = &cobra.Command{
	Use:   "accept",
	Short: "mark a KEP as accepted by a SIG, allowing for further refinement before implementation",
	Long: `
Accept the motivation for a KEP on behalf of a SIG. Accept communicates that the
SIG agrees with the addressing the need described in the KEP; and the SIG is going
to dedicate resources to shepherding the enhancement through its lifecycle.

Accept will likely be run by a SIG Chair, TL, or their delegate according to the
OWNERS file for the SIG. Accept will mark the KEP as "provisional" to signal
that the KEP author(s) should be working towards the inclusion of draft development,
teaching, and operations guides for approval before implementation`,
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
		err = workflow.Accept(runtimeSettings)
		if err != nil {
			return err
		}

		fmt.Println("successfully marked KEP as accepted!")
		return nil
	},
}
