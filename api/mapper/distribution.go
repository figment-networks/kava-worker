package mapper

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/figment-networks/indexing-engine/structs"
	"github.com/figment-networks/kava-worker/api/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/kava-labs/kava/app"
	"github.com/tendermint/tendermint/libs/bech32"
)

var zero big.Int

func DistributionWithdrawValidatorCommissionToSub(msg sdk.Msg, logf types.LogFormat) (se structs.SubsetEvent, err error) {
	wvc, ok := msg.(distribution.MsgWithdrawValidatorCommission)
	if !ok {
		return se, errors.New("Not a withdraw_validator_commission type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(bech32ValPrefix, wvc.ValidatorAddress.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting ValidatorAddress: %w", err)
	}

	se = structs.SubsetEvent{
		Type:   []string{"withdraw_validator_commission"},
		Module: "distribution",
		Node:   map[string][]structs.Account{"validator": {{ID: bech32Addr}}},
		Recipient: []structs.EventTransfer{{
			Account: structs.Account{ID: bech32Addr},
		}},
	}

	err = produceTransfers(&se, "send", "", logf)
	return se, err
}

func DistributionSetWithdrawAddressToSub(msg sdk.Msg) (se structs.SubsetEvent, er error) {
	swa, ok := msg.(distribution.MsgSetWithdrawAddress)
	if !ok {
		return se, errors.New("Not a set_withdraw_address type")
	}

	bech32DelAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, swa.DelegatorAddress.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting DelegatorAddress: %w", err)
	}

	bech32WithdrawAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, swa.WithdrawAddress.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting WithdrawAddress: %w", err)
	}

	return structs.SubsetEvent{
		Type:   []string{"set_withdraw_address"},
		Module: "distribution",
		Node: map[string][]structs.Account{
			"delegator": {{ID: bech32DelAddr}},
			"withdraw":  {{ID: bech32WithdrawAddr}},
		},
	}, nil
}

func DistributionWithdrawDelegatorRewardToSub(msg sdk.Msg, logf types.LogFormat) (se structs.SubsetEvent, err error) {
	wdr, ok := msg.(distribution.MsgWithdrawDelegatorReward)
	if !ok {
		return se, errors.New("Not a withdraw_validator_commission type")
	}
	bech32DelAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, wdr.DelegatorAddress.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting DelegatorAddress: %w", err)
	}
	bech32ValAddr, err := bech32.ConvertAndEncode(bech32ValPrefix, wdr.ValidatorAddress.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting ValidatorAddress: %w", err)
	}

	se = structs.SubsetEvent{
		Type:   []string{"withdraw_delegator_reward"},
		Module: "distribution",
		Node: map[string][]structs.Account{
			"delegator": {{ID: bech32DelAddr}},
			"validator": {{ID: bech32ValAddr}},
		},
		Recipient: []structs.EventTransfer{{
			Account: structs.Account{ID: bech32DelAddr},
		}},
	}

	err = produceTransfers(&se, "reward", "", logf)
	return se, err
}

func DistributionFundCommunityPoolToSub(msg sdk.Msg, logf types.LogFormat) (se structs.SubsetEvent, er error) {
	fcp, ok := msg.(distributiontypes.MsgFundCommunityPool)
	if !ok {
		return se, errors.New("Not a withdraw_validator_commission type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, fcp.Depositor.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting DepositorAddress: %w", err)
	}

	evt, err := distributionProduceEvTx(fcp.Depositor, fcp.Amount)
	se = structs.SubsetEvent{
		Type:   []string{"fund_community_pool"},
		Module: "distribution",
		Node: map[string][]structs.Account{
			"depositor": {{ID: bech32Addr}},
		},
		Sender: []structs.EventTransfer{evt},
	}
	err = produceTransfers(&se, "send", "", logf)
	return se, err
}

func distributionProduceEvTx(account sdk.AccAddress, coins sdk.Coins) (evt structs.EventTransfer, err error) {
	evt = structs.EventTransfer{
		Account: structs.Account{ID: account.String()},
	}
	if len(coins) > 0 {
		evt.Amounts = []structs.TransactionAmount{}
		for _, coin := range coins {
			txa := structs.TransactionAmount{
				Currency: coin.Denom,
				Text:     coin.Amount.String(),
			}

			txa.Numeric.Set(coin.Amount.BigInt())
			evt.Amounts = append(evt.Amounts, txa)
		}
	}

	return evt, nil
}
