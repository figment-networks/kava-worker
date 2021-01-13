package api

import (
	"errors"

	shared "github.com/figment-networks/indexer-manager/structs"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/x/harvest"
)

func mapHarvestDepositToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	u, ok := msg.(harvest.MsgDeposit)
	if !ok {
		return se, errors.New("Not a harvest_deposit type")
	}

	return shared.SubsetEvent{
		Type:   []string{"harvest_deposit"},
		Module: "harvest",
		Node: map[string][]shared.Account{
			"depositor": {{ID: u.Depositor.String()}},
		},
		Amount: map[string]shared.TransactionAmount{
			"send": {
				Currency: u.Amount.Denom,
				Numeric:  u.Amount.Amount.BigInt(),
				Text:     u.Amount.String(),
			},
		},
		Additional: map[string][]string{
			"deposit_type": []string{u.DepositType},
		},
	}, nil
}

func mapHarvestWithdrawToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	u, ok := msg.(harvest.MsgWithdraw)
	if !ok {
		return se, errors.New("Not a harvest_withdraw type")
	}

	se = shared.SubsetEvent{
		Type:   []string{"harvest_deposit"},
		Module: "harvest",
		Node: map[string][]shared.Account{
			"depositor": {{ID: u.Depositor.String()}},
		},
		Amount: map[string]shared.TransactionAmount{
			"send": {
				Currency: u.Amount.Denom,
				Numeric:  u.Amount.Amount.BigInt(),
				Text:     u.Amount.String(),
			},
		},
		Additional: map[string][]string{
			"deposit_type": []string{u.DepositType},
		},
	}

	return se, nil
}

func mapHarvestClaimRewardToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	u, ok := msg.(harvest.MsgClaimReward)
	if !ok {
		return se, errors.New("Not a claim_harvest_reward type")
	}

	se = shared.SubsetEvent{
		Type:   []string{"claim_harvest_reward"},
		Module: "harvest",
		Node: map[string][]shared.Account{
			"sender":   {{ID: u.Sender.String()}},
			"receiver": {{ID: u.Receiver.String()}},
		},
		Amount: map[string]shared.TransactionAmount{
			"reward": {
				Currency: u.DepositDenom,
				// Numeric:  u.Amount.Amount.BigInt(), // todo get amount from logs?
				// Text:     u.Amount.String(),
			},
		},
		Additional: map[string][]string{
			"deposit_type":    []string{u.DepositType},
			"multiplier_name": []string{u.MultiplierName},
		},
	}

	return se, nil
}
