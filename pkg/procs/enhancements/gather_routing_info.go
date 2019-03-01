package enhancements

import (
	"fmt"
	"strings"

	"github.com/calebamiles/keps/pkg/changes/changeset"
	"github.com/calebamiles/keps/pkg/changes/routing"
	"github.com/calebamiles/keps/pkg/keps"
	"github.com/calebamiles/keps/pkg/keps/events"
	"github.com/calebamiles/keps/pkg/orgs"
	"github.com/calebamiles/keps/pkg/settings"
)

func GatherProposalInfoFrom(runtime settings.Runtime, org orgs.Instance, kep keps.Instance) (routing.Info, error) {
	kepSummary, err := kep.Summary()
	if err != nil {
		return nil, err
	}

	title := changeset.Title(strings.Title(fmt.Sprintf("propose %s for sponsorship by %s", kep.Title(), kep.OwningSIG())))
	fullDescription := changeset.FullDescription(kepSummary)
	shortSummary := changeset.ShortSummary(strings.Title(fmt.Sprintf("propose %s for sponsorship by %s", kep.Title(), kep.OwningSIG())))

	var changeReceipt = func() string {
		return kep.GetLifecyclePR(events.Proposal)
	}

	description, err := changeset.Describe(title, fullDescription, shortSummary, changeReceipt)
	if err != nil {
		return nil, err
	}

	routingInfo, err := routing.NewInfoFrom(runtime, kep, org, description)
	if err != nil {
		return nil, err
	}

	return routingInfo, nil
}
