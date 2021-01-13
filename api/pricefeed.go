package api

import (
	"errors"

	shared "github.com/figment-networks/indexer-manager/structs"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/x/pricefeed"
)

func mapPricefeedPostPrice(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	u, ok := msg.(pricefeed.MsgPostPrice)
	if !ok {
		return se, errors.New("Not a post_price type")
	}

	return shared.SubsetEvent{
		Type:   []string{"post_price"},
		Module: "pricefeed",
		Node: map[string][]shared.Account{
			"from": {{ID: u.From.String()}},
		},
		Amount: map[string]shared.TransactionAmount{
			"value": {
				Currency: u.MarketID, // todo "xrp:usd" should this be split?
				Numeric:  u.Price.BigInt(),
				Text:     u.Price.String(),
			},
		},
		Additional: map[string][]string{
			"expiry": []string{u.Expiry.String()},
		},
	}, nil
}
