package api

import (
	"errors"
	"fmt"

	shared "github.com/figment-networks/indexer-manager/structs"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crisis "github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/kava-labs/kava/app"
	"github.com/tendermint/tendermint/libs/bech32"
)

func mapCrisisVerifyInvariantToSub(msg sdk.Msg) (se shared.SubsetEvent, er error) {
	mvi, ok := msg.(crisis.MsgVerifyInvariant)
	if !ok {
		return se, errors.New("Not a verify_invariant type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, mvi.Sender.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting SenderAddress: %w", err)
	}

	return shared.SubsetEvent{
		Type:   []string{"verify_invariant"},
		Module: "crisis",
		Sender: []shared.EventTransfer{{
			Account: shared.Account{ID: bech32Addr},
		}},
		Additional: map[string][]string{
			"invariant_route":       {mvi.InvariantRoute},
			"invariant_module_name": {mvi.InvariantModuleName},
		},
	}, nil
}
