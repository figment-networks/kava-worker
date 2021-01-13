package api

import (
	"errors"
	"fmt"

	shared "github.com/figment-networks/indexer-manager/structs"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/app"
	"github.com/kava-labs/kava/x/pricefeed"
	"github.com/tendermint/tendermint/libs/bech32"
)

func mapPricefeedPostPrice(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	m, ok := msg.(pricefeed.MsgPostPrice)
	if !ok {
		return se, errors.New("Not a post_price type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.From.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting Address: %w", err)
	}

	return shared.SubsetEvent{
		Type:   []string{"post_price"},
		Module: "pricefeed",
		Node: map[string][]shared.Account{
			"from": {{ID: bech32Addr}},
		},
		Amount: map[string]shared.TransactionAmount{
			"value": {
				Currency: m.MarketID, // todo "xrp:usd" should this be split?
				Numeric:  m.Price.BigInt(),
				Text:     m.Price.String(),
			},
		},
		Additional: map[string][]string{
			"expiry": []string{m.Expiry.String()},
		},
	}, nil
}
