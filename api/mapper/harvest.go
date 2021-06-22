package mapper

/*
import (
	"errors"
	"fmt"

	shared "github.com/figment-networks/indexer-manager/structs"
	"github.com/figment-networks/kava-worker/api/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/app"
	"github.com/tendermint/tendermint/libs/bech32"
)
func HarvestDepositToSub(msg sdk.Msg, logf types.LogFormat) (se shared.SubsetEvent, err error) {
	m, ok := msg.(harvest.MsgDeposit)
	if !ok {
		return se, errors.New("Not a harvest_deposit type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Depositor.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting DepositorAddress: %w", err)
	}

	se = shared.SubsetEvent{
		Type:   []string{"harvest_deposit"},
		Module: "harvest",
		Node: map[string][]shared.Account{
			"depositor": {{ID: bech32Addr}},
		},
		Amount: map[string]shared.TransactionAmount{
			"send": {
				Currency: m.Amount.Denom,
				Numeric:  m.Amount.Amount.BigInt(),
				Text:     m.Amount.String(),
			},
		},
		Additional: map[string][]string{
			"deposit_type": []string{m.DepositType},
		},
	}

	err = produceTransfers(&se, "send", "", logf)
	return se, err
}

func HarvestWithdrawToSub(msg sdk.Msg, logf types.LogFormat) (se shared.SubsetEvent, err error) {
	m, ok := msg.(harvest.MsgWithdraw)
	if !ok {
		return se, errors.New("Not a harvest_withdraw type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Depositor.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting DepositorAddress: %w", err)
	}

	se = shared.SubsetEvent{
		Type:   []string{"harvest_withdraw"},
		Module: "harvest",
		Node: map[string][]shared.Account{
			"depositor": {{ID: bech32Addr}},
		},
		Amount: map[string]shared.TransactionAmount{
			"send": {
				Currency: m.Amount.Denom,
				Numeric:  m.Amount.Amount.BigInt(),
				Text:     m.Amount.String(),
			},
		},
		Additional: map[string][]string{
			"deposit_type": []string{m.DepositType},
		},
	}

	err = produceTransfers(&se, "send", "", logf)
	return se, err
}

func HarvestClaimRewardToSub(msg sdk.Msg, logf types.LogFormat) (se shared.SubsetEvent, err error) {
	m, ok := msg.(harvest.MsgClaimReward)
	if !ok {
		return se, errors.New("Not a claim_harvest_reward type")
	}

	bech32SenderAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Sender.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting SenderAddress: %w", err)
	}

	bech32RecAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Receiver.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting ReceiverAddress: %w", err)
	}

	se = shared.SubsetEvent{
		Type:   []string{"claim_harvest_reward"},
		Module: "harvest",
		Node: map[string][]shared.Account{
			"sender":   {{ID: bech32SenderAddr}},
			"receiver": {{ID: bech32RecAddr}},
		},
		Additional: map[string][]string{
			"deposit_denom":   []string{m.DepositDenom},
			"deposit_type":    []string{m.DepositType},
			"multiplier_name": []string{m.MultiplierName},
		},
	}

	err = produceTransfers(&se, "reward", "", logf)
	return se, err
}
*/
