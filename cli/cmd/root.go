package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "keps-cli",
	Short: "a basic CLI for interacting with KEPs",
	Long: `

The KEPs CLI helps KEP authors, editors, reviewers, and approvers
perform CRUD type operations against the existing KEP content. Most interactions
with KEPs are expected to use either the CLI or the KEP library code which this
CLI exercises. The happy path is as follows:

1. [author] kep init <path-to-location-where-kep-should-live>
2. [author] kep propose <path-to-created-kep>
3. [SIG]    kep accept <path-to-created-kep>
4. [author] kep plan <path-to-created-kep>
5. [SIG]    kep approve <path-to-created-kep>`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// load all the available commands here rather than in init() functions sprinkled everywhere
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(proposeCmd)
	rootCmd.AddCommand(acceptCmd)
	rootCmd.AddCommand(planCmd)
	rootCmd.AddCommand(approveCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
