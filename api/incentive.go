package api

import (
	"errors"
	"fmt"

	shared "github.com/figment-networks/indexer-manager/structs"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/app"
	"github.com/kava-labs/kava/x/incentive"
	"github.com/tendermint/tendermint/libs/bech32"
)

func mapIncentiveClaimRewardToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	m, ok := msg.(incentive.MsgClaimReward)
	if !ok {
		return se, errors.New("Not a claim_reward type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Sender.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting Address: %w", err)
	}

	return shared.SubsetEvent{
		Type:   []string{"claim_reward"},
		Module: "incentive",
		Node: map[string][]shared.Account{
			"sender": {{ID: bech32Addr}},
		}, // todo amount?
		Additional: map[string][]string{
			"multiplier_name": []string{m.MultiplierName},
			"collateral_type": []string{m.CollateralType},
		},
	}, nil
}
