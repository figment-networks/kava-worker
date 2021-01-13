package api

import (
	"errors"
	"strconv"

	shared "github.com/figment-networks/indexer-manager/structs"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/x/committee"
)

func mapCommitteeSubmitProposalToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	m, ok := msg.(committee.MsgSubmitProposal)
	if !ok {
		return se, errors.New("Not a commmittee_submit_proposal type")
	}

	return shared.SubsetEvent{
		Type:   []string{"commmittee_submit_proposal"},
		Module: "commmittee",
		Node: map[string][]shared.Account{
			"proposer": {{ID: m.Proposer.String()}},
		},
		Additional: map[string][]string{
			"title":          []string{m.PubProposal.GetTitle()},
			"description":    []string{m.PubProposal.GetDescription()},
			"proposal_route": []string{m.PubProposal.ProposalRoute()},
			"proposal_type":  []string{m.PubProposal.ProposalType()},
		},
	}, nil
}

func mapCommitteeVoteToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	m, ok := msg.(committee.MsgVote)
	if !ok {
		return se, errors.New("Not a committee_vote type")
	}

	return shared.SubsetEvent{
		Type:   []string{"committee_vote"},
		Module: "commmittee",
		Node: map[string][]shared.Account{
			"vote": {{ID: m.Voter.String()}},
		},
		Additional: map[string][]string{
			"proposal_id": []string{strconv.FormatUint(m.ProposalID, 10)},
		},
	}, nil
}
