package mapper

import (
	"errors"
	"fmt"
	"strconv"

	shared "github.com/figment-networks/indexer-manager/structs"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/app"
	"github.com/kava-labs/kava/x/committee"
	"github.com/tendermint/tendermint/libs/bech32"
)

func CommitteeSubmitProposalToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	m, ok := msg.(committee.MsgSubmitProposal)
	if !ok {
		return se, errors.New("Not a commmittee_submit_proposal type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Proposer.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting ProposerAddress: %w", err)
	}

	se = shared.SubsetEvent{
		Type:   []string{"commmittee_submit_proposal"},
		Module: "commmittee",
		Node: map[string][]shared.Account{
			"proposer": {{ID: bech32Addr}},
		},
	}

	se.Additional = map[string][]string{}

	if m.PubProposal.ProposalRoute() != "" {
		se.Additional["proposal_route"] = []string{m.PubProposal.ProposalRoute()}
	}
	if m.PubProposal.ProposalType() != "" {
		se.Additional["proposal_type"] = []string{m.PubProposal.ProposalType()}
	}
	if m.PubProposal.GetDescription() != "" {
		se.Additional["descritpion"] = []string{m.PubProposal.GetDescription()}
	}
	if m.PubProposal.GetTitle() != "" {
		se.Additional["title"] = []string{m.PubProposal.GetTitle()}
	}
	if m.PubProposal.String() != "" {
		se.Additional["content"] = []string{m.PubProposal.String()}
	}

	return se, nil
}

func CommitteeVoteToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	m, ok := msg.(committee.MsgVote)
	if !ok {
		return se, errors.New("Not a committee_vote type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Voter.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting VoterAddress: %w", err)
	}

	return shared.SubsetEvent{
		Type:   []string{"committee_vote"},
		Module: "commmittee",
		Node: map[string][]shared.Account{
			"vote": {{ID: bech32Addr}},
		},
		Additional: map[string][]string{
			"proposal_id": []string{strconv.FormatUint(m.ProposalID, 10)},
		},
	}, nil
}
