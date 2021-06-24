package mapper

import (
	"errors"
	"fmt"

	"github.com/figment-networks/indexer-search/structs"
	"github.com/figment-networks/kava-worker/api/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/app"
	"github.com/kava-labs/kava/x/incentive"
	"github.com/tendermint/tendermint/libs/bech32"
)

func IncentiveClaimUSDXMintingRewardToSub(msg sdk.Msg, logf types.LogFormat) (se structs.SubsetEvent, err error) {
	m, ok := msg.(incentive.MsgClaimUSDXMintingReward)
	if !ok {
		return se, errors.New("Not a claim_usdx_minting_reward type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Sender.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting Address: %w", err)
	}

	se = structs.SubsetEvent{
		Type:   []string{"claim_usdx_minting_reward"},
		Module: "incentive",
		Node: map[string][]structs.Account{
			"sender": {{ID: bech32Addr}},
		},
		Additional: map[string][]string{
			"multiplier_name": {m.MultiplierName},
		},
	}

	err = produceTransfers(&se, "reward", "", logf)
	return se, err
}

func IncentiveClaimHardRewardToSub(msg sdk.Msg, logf types.LogFormat) (se structs.SubsetEvent, err error) {
	m, ok := msg.(incentive.MsgClaimHardReward)
	if !ok {
		return se, errors.New("Not a claim_hard_reward type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Sender.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting Address: %w", err)
	}

	se = structs.SubsetEvent{
		Type:   []string{"claim_hard_reward"},
		Module: "incentive",
		Node: map[string][]structs.Account{
			"sender": {{ID: bech32Addr}},
		},
		Additional: map[string][]string{
			"multiplier_name": {m.MultiplierName},
		},
	}

	err = produceTransfers(&se, "reward", "", logf)
	return se, err
}
