package mapper

import (
	"errors"
	"fmt"

	"github.com/figment-networks/indexing-engine/structs"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/tendermint/tendermint/libs/bech32"
)

func SlashingUnjailToSub(msg sdk.Msg) (se structs.SubsetEvent, er error) {
	unjail, ok := msg.(slashing.MsgUnjail)
	if !ok {
		return se, errors.New("Not an unjail type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(bech32ValPrefix, unjail.ValidatorAddr.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting ValidatorAddr: %w", err)
	}

	return structs.SubsetEvent{
		Type:   []string{"unjail"},
		Module: "slashing",
		Node:   map[string][]structs.Account{"validator": {{ID: bech32Addr}}},
	}, nil
}
