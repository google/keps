package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/calebamiles/keps/pkg/settings"
	"github.com/calebamiles/keps/pkg/workflow"
)

// proposeCmd represents the propose command
var proposeCmd = &cobra.Command{
	Use:   "propose",
	Short: "propose a new KEP for SIG review",
	Long: `
Propose a new KEP for sponsorship by a SIG. The KEP process is
an iterative process for managing "Change at the Speed of Trust". Proposing an
enhancement for SIG reviewal frames a problem or gap for SIG consideration and
allows a SIG to decide whether SIG resources should be dedicated to shepherd the
enhancement to completion. The propose command serves as a preflight checklist
to ensure SIGs spend their limited bandwidth reviewing high quality proposals.`,
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
		err = workflow.Propose(runtimeSettings)
		if err != nil {
			return err
		}

		fmt.Println("KEP is ready for proposal!")
		return nil
	},
}

// TODO maybe kill off these init() functions
// the best worst place for them might be inside of cmd.Execute() which lives inside of root.go
func init() {
	//rootCmd.AddCommand(proposeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// proposeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// proposeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
