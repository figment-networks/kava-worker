package mapper

import (
	"errors"
	"fmt"
	"strconv"

	shared "github.com/figment-networks/indexer-manager/structs"
	"github.com/figment-networks/kava-worker/api/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/kava-labs/kava/app"
	"github.com/tendermint/tendermint/libs/bech32"
)

func GovDepositToSub(msg sdk.Msg, logf types.LogFormat) (se shared.SubsetEvent, err error) {
	dep, ok := msg.(gov.MsgDeposit)
	if !ok {
		return se, errors.New("Not a deposit type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, dep.Depositor.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting DepositorAddress: %w", err)
	}

	se = shared.SubsetEvent{
		Type:       []string{"deposit"},
		Module:     "gov",
		Node:       map[string][]shared.Account{"depositor": {{ID: bech32Addr}}},
		Additional: map[string][]string{"proposalID": {strconv.FormatUint(dep.ProposalID, 10)}},
	}

	sender := shared.EventTransfer{Account: shared.Account{ID: bech32Addr}}
	txAmount := map[string]shared.TransactionAmount{}

	for i, coin := range dep.Amount {
		am := shared.TransactionAmount{
			Currency: coin.Denom,
			Numeric:  coin.Amount.BigInt(),
			Text:     coin.Amount.String(),
		}

		sender.Amounts = append(sender.Amounts, am)
		key := "deposit"
		if i > 0 {
			key += "_" + strconv.Itoa(i)
		}

		txAmount[key] = am
	}

	se.Sender = []shared.EventTransfer{sender}
	se.Amount = txAmount

	err = produceTransfers(&se, "send", logf)
	return se, err
}

func GovVoteToSub(msg sdk.Msg) (se shared.SubsetEvent, er error) {
	vote, ok := msg.(gov.MsgVote)
	if !ok {
		return se, errors.New("Not a vote type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, vote.Voter.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting VoterAddress: %w", err)
	}

	return shared.SubsetEvent{
		Type:   []string{"vote"},
		Module: "gov",
		Node:   map[string][]shared.Account{"voter": {{ID: bech32Addr}}},
		Additional: map[string][]string{
			"proposalID": {strconv.FormatUint(vote.ProposalID, 10)},
			"option":     {vote.Option.String()},
		},
	}, nil
}

func GovSubmitProposalToSub(msg sdk.Msg, logf types.LogFormat) (se shared.SubsetEvent, err error) {
	sp, ok := msg.(gov.MsgSubmitProposal)
	if !ok {
		return se, errors.New("Not a submit_proposal type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, sp.Proposer.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting ProposerAddress: %w", err)
	}

	se = shared.SubsetEvent{
		Type:   []string{"submit_proposal"},
		Module: "gov",
		Node:   map[string][]shared.Account{"proposer": {{ID: bech32Addr}}},
	}

	sender := shared.EventTransfer{Account: shared.Account{ID: bech32Addr}}
	txAmount := map[string]shared.TransactionAmount{}

	for i, coin := range sp.InitialDeposit {
		am := shared.TransactionAmount{
			Currency: coin.Denom,
			Numeric:  coin.Amount.BigInt(),
			Text:     coin.Amount.String(),
		}

		sender.Amounts = append(sender.Amounts, am)
		key := "initial_deposit"
		if i > 0 {
			key += "_" + strconv.Itoa(i)
		}

		txAmount[key] = am
	}
	se.Sender = []shared.EventTransfer{sender}
	se.Amount = txAmount

	se.Additional = map[string][]string{}

	if sp.Content.ProposalRoute() != "" {
		se.Additional["proposal_route"] = []string{sp.Content.ProposalRoute()}
	}
	if sp.Content.ProposalType() != "" {
		se.Additional["proposal_type"] = []string{sp.Content.ProposalType()}
	}
	if sp.Content.GetDescription() != "" {
		se.Additional["descritpion"] = []string{sp.Content.GetDescription()}
	}
	if sp.Content.GetTitle() != "" {
		se.Additional["title"] = []string{sp.Content.GetTitle()}
	}
	if sp.Content.String() != "" {
		se.Additional["content"] = []string{sp.Content.String()}
	}

	err = produceTransfers(&se, "send", logf)
	return se, err
}
