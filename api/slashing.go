package api

import (
	"errors"
	"fmt"

	shared "github.com/figment-networks/indexer-manager/structs"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/kava-labs/kava/app"
	"github.com/tendermint/tendermint/libs/bech32"
)

func mapSlashingUnjailToSub(msg sdk.Msg) (se shared.SubsetEvent, er error) {
	unjail, ok := msg.(slashing.MsgUnjail)
	if !ok {
		return se, errors.New("Not an unjail type")
	}

	//todo
	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, unjail.ValidatorAddr.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting ValidatorAddr: %w", err)
	}

	return shared.SubsetEvent{
		Type:   []string{"unjail"},
		Module: "slashing",
		Node:   map[string][]shared.Account{"validator": {{ID: bech32Addr}}},
	}, nil
}
