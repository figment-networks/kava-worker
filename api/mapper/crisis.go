package mapper

import (
	"errors"
	"fmt"

	"github.com/figment-networks/indexing-engine/structs"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crisis "github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/kava-labs/kava/app"
	"github.com/tendermint/tendermint/libs/bech32"
)

func CrisisVerifyInvariantToSub(msg sdk.Msg) (se structs.SubsetEvent, er error) {
	mvi, ok := msg.(crisis.MsgVerifyInvariant)
	if !ok {
		return se, errors.New("Not a verify_invariant type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, mvi.Sender.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting SenderAddress: %w", err)
	}

	return structs.SubsetEvent{
		Type:   []string{"verify_invariant"},
		Module: "crisis",
		Sender: []structs.EventTransfer{{
			Account: structs.Account{ID: bech32Addr},
		}},
		Additional: map[string][]string{
			"invariant_route":       {mvi.InvariantRoute},
			"invariant_module_name": {mvi.InvariantModuleName},
		},
	}, nil
}
