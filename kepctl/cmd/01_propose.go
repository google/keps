package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/calebamiles/keps/pkg/changes/changeset"
	"github.com/calebamiles/keps/pkg/changes/inplace"
	"github.com/calebamiles/keps/pkg/changes/routing"
	"github.com/calebamiles/keps/pkg/keps"
	"github.com/calebamiles/keps/pkg/keps/events"
	"github.com/calebamiles/keps/pkg/orgs"
	"github.com/calebamiles/keps/pkg/settings"
	"github.com/calebamiles/keps/pkg/workflow"
)

// proposeCmd represents the propose command
var proposeCmd = &cobra.Command{
	Use:   "propose",
	Short: "propose a new KEP for SIG review",
	Long:  DescribePropose,
	Args:  cobra.ExactArgs(1), // accept just one argument, location of KEP
	RunE: func(cmd *cobra.Command, args []string) error {

		targetPath := args[0] // we have a validator ensuring we will have exactly one positional

		runtime, kep, org, err := Setup(targetPath)
		if err != nil {
			return err
		}

		sucessMsg, err := ExecutePropose(runtime, org, kep)
		if err != nil {
			return err
		}

		fmt.Println(successMsg)
		return nil
	},
}

const DescribePropose = `
Propose a new KEP for sponsorship by a SIG. The KEP process is
an iterative process for managing "Change at the Speed of Trust". Proposing an
enhancement for SIG reviewal frames a problem or gap for SIG consideration and
allows a SIG to decide whether SIG resources should be dedicated to shepherd the
enhancement to completion. The propose command serves as a preflight checklist
to ensure SIGs spend their limited bandwidth reviewing high quality proposals.
`

func ExecutePropose(runtime settings.Runtime, org orgs.Instance, kep keps.Instance) (string, error) {
	prUrl, err := workflow.Propose(runtime, org, kep)
	if err != nil {
		return "", err
	}

	successMsg := fmt.Sprintf("KEP has been proposed!\n continue the discussion at:\n%s", prUrl)
	return successMsg, nil
}
