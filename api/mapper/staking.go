package mapper

import (
	"errors"
	"fmt"
	"math/big"

	shared "github.com/figment-networks/indexer-manager/structs"
	"github.com/figment-networks/kava-worker/api/types"
	"github.com/figment-networks/kava-worker/api/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/kava-labs/kava/app"
	"github.com/tendermint/tendermint/libs/bech32"
)

const unbondedTokensPoolAddr = "kava1tygms3xhhs3yv487phx3dw4a95jn7t7lawprey"

func StakingUndelegateToSub(msg sdk.Msg, logf types.LogFormat) (se shared.SubsetEvent, err error) {
	u, ok := msg.(staking.MsgUndelegate)
	if !ok {
		return se, errors.New("Not a begin_unbonding type")
	}

	bech32DelAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, u.DelegatorAddress.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting DelegatorAddress: %w", err)
	}

	bech32ValAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, u.ValidatorAddress.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting ValidatorAddress: %w", err)
	}

	se = shared.SubsetEvent{
		Type:   []string{"begin_unbonding"},
		Module: "staking",
		Node: map[string][]shared.Account{
			"delegator": {{ID: bech32DelAddr}},
			"validator": {{ID: bech32ValAddr}},
		},
		Amount: map[string]shared.TransactionAmount{
			"undelegate": {
				Currency: u.Amount.Denom,
				Numeric:  u.Amount.Amount.BigInt(),
				Text:     u.Amount.String(),
			},
		},
	}

	var withdrawAddr string
	rewards := []shared.TransactionAmount{}
	for _, ev := range logf.Events {
		if ev.Type != "transfer" {
			continue
		}

		var latestRecipient string
		for _, attr := range ev.Attributes {
			if len(attr.Recipient) > 0 {
				latestRecipient = attr.Recipient[0]
			}
			if latestRecipient == unbondedTokensPoolAddr {
				continue
			}
			withdrawAddr = latestRecipient

			for _, amount := range attr.Amount {
				attrAmt := shared.TransactionAmount{Numeric: &big.Int{}}
				sliced := util.GetCurrency(amount)
				var (
					c       *big.Int
					exp     int32
					coinErr error
				)
				if len(sliced) == 3 {
					attrAmt.Currency = sliced[2]
					c, exp, coinErr = util.GetCoin(sliced[1])
				} else {
					c, exp, coinErr = util.GetCoin(amount)
				}
				if coinErr != nil {
					return se, fmt.Errorf("[COSMOS-API] Error parsing amount '%s': %s ", amount, coinErr)
				}
				attrAmt.Text = amount
				attrAmt.Numeric.Set(c)
				attrAmt.Exp = exp
				if attrAmt.Numeric.Cmp(&zero) != 0 {
					rewards = append(rewards, attrAmt)
				}
			}
		}
	}

	if len(rewards) == 0 {
		return se, nil
	}
	se.Transfers = map[string][]shared.EventTransfer{
		"reward": []shared.EventTransfer{{
			Amounts: rewards,
			Account: shared.Account{ID: withdrawAddr},
		}},
	}

	return se, nil
}

func StakingDelegateToSub(msg sdk.Msg, logf types.LogFormat) (se shared.SubsetEvent, err error) {
	d, ok := msg.(staking.MsgDelegate)
	if !ok {
		return se, errors.New("Not a delegate type")
	}

	bech32DelAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, d.DelegatorAddress.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting DelegatorAddress: %w", err)
	}

	bech32ValAddr, err := bech32.ConvertAndEncode(bech32ValPrefix, d.ValidatorAddress.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting ValidatorAddress: %w", err)
	}

	se = shared.SubsetEvent{
		Type:   []string{"delegate"},
		Module: "staking",
		Node: map[string][]shared.Account{
			"delegator": {{ID: bech32DelAddr}},
			"validator": {{ID: bech32ValAddr}},
		},
		Amount: map[string]shared.TransactionAmount{
			"delegate": {
				Currency: d.Amount.Denom,
				Numeric:  d.Amount.Amount.BigInt(),
				Text:     d.Amount.String(),
			},
		},
	}

	err = produceTransfers(&se, "reward", logf)
	return se, err
}

func StakingBeginRedelegateToSub(msg sdk.Msg, logf types.LogFormat) (se shared.SubsetEvent, err error) {
	br, ok := msg.(staking.MsgBeginRedelegate)
	if !ok {
		return se, errors.New("Not a begin_redelegate type")
	}

	bech32DelAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, br.DelegatorAddress.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting DelegatorAddress: %w", err)
	}

	bech32ValDstAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, br.ValidatorDstAddress.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting ValidatorDstAddress: %w", err)
	}

	bech32ValSrcAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, br.ValidatorSrcAddress.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting ValidatorSrcAddress: %w", err)
	}

	se = shared.SubsetEvent{
		Type:   []string{"begin_redelegate"},
		Module: "staking",
		Node: map[string][]shared.Account{
			"delegator":             {{ID: bech32DelAddr}},
			"validator_destination": {{ID: bech32ValDstAddr}},
			"validator_source":      {{ID: bech32ValSrcAddr}},
		},
		Amount: map[string]shared.TransactionAmount{
			"delegate": {
				Currency: br.Amount.Denom,
				Numeric:  br.Amount.Amount.BigInt(),
				Text:     br.Amount.String(),
			},
		},
	}

	err = produceTransfers(&se, "reward", logf)
	return se, err
}

func StakingCreateValidatorToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	ev, ok := msg.(staking.MsgCreateValidator)
	if !ok {
		return se, errors.New("Not a create_validator type")
	}

	bech32DelAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, ev.DelegatorAddress.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting DelegatorAddress: %w", err)
	}
	bech32ValAddr, err := bech32.ConvertAndEncode(bech32ValPrefix, ev.ValidatorAddress.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting ValidatorAddress: %w", err)
	}

	return shared.SubsetEvent{
		Type:   []string{"create_validator"},
		Module: "distribution",
		Node: map[string][]shared.Account{
			"delegator": {{ID: bech32DelAddr}},
			"validator": {
				{
					ID: bech32ValAddr,
					Details: &shared.AccountDetails{
						Name:        ev.Description.Moniker,
						Description: ev.Description.Details,
						Contact:     ev.Description.SecurityContact,
						Website:     ev.Description.Website,
					},
				},
			},
		},
		Amount: map[string]shared.TransactionAmount{
			"self_delegation": {
				Currency: ev.Value.Denom,
				Numeric:  ev.Value.Amount.BigInt(),
				Text:     ev.Value.String(),
			},
			"self_delegation_min": {
				Text:    ev.MinSelfDelegation.String(),
				Numeric: ev.MinSelfDelegation.BigInt(),
			},
			"commission_rate": {
				Text:    ev.Commission.Rate.String(),
				Numeric: ev.Commission.Rate.Int,
			},
			"commission_max_rate": {
				Text:    ev.Commission.MaxRate.String(),
				Numeric: ev.Commission.MaxRate.Int,
			},
			"commission_max_change_rate": {
				Text:    ev.Commission.MaxChangeRate.String(),
				Numeric: ev.Commission.MaxChangeRate.Int,
			}},
	}, err
}

func StakingEditValidatorToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	ev, ok := msg.(staking.MsgEditValidator)
	if !ok {
		return se, errors.New("Not a edit_validator type")
	}

	bech32ValAddr, err := bech32.ConvertAndEncode(bech32ValPrefix, ev.ValidatorAddress.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting ValidatorAddress: %w", err)
	}

	sev := shared.SubsetEvent{
		Type:   []string{"edit_validator"},
		Module: "distribution",
		Node: map[string][]shared.Account{
			"validator": {
				{
					ID: bech32ValAddr,
					Details: &shared.AccountDetails{
						Name:        ev.Description.Moniker,
						Description: ev.Description.Details,
						Contact:     ev.Description.SecurityContact,
						Website:     ev.Description.Website,
					},
				},
			},
		},
	}

	if ev.MinSelfDelegation != nil || ev.CommissionRate != nil {
		sev.Amount = map[string]shared.TransactionAmount{}
		if ev.MinSelfDelegation != nil {
			sev.Amount["self_delegation_min"] = shared.TransactionAmount{
				Text:    ev.MinSelfDelegation.String(),
				Numeric: ev.MinSelfDelegation.BigInt(),
			}
		}

		if ev.CommissionRate != nil {
			sev.Amount["commission_rate"] = shared.TransactionAmount{
				Text:    ev.CommissionRate.String(),
				Numeric: ev.CommissionRate.Int,
			}
		}
	}
	return sev, err
}