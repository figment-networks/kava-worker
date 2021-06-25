package mapper

import (
	"errors"
	"fmt"

	"github.com/figment-networks/indexing-engine/structs"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/app"
	"github.com/kava-labs/kava/x/pricefeed"
	"github.com/tendermint/tendermint/libs/bech32"
)

func PricefeedPostPrice(msg sdk.Msg) (se structs.SubsetEvent, err error) {
	m, ok := msg.(pricefeed.MsgPostPrice)
	if !ok {
		return se, errors.New("Not a post_price type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.From.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting Address: %w", err)
	}

	return structs.SubsetEvent{
		Type:   []string{"post_price"},
		Module: "pricefeed",
		Node: map[string][]structs.Account{
			"from": {{ID: bech32Addr}},
		},
		Amount: map[string]structs.TransactionAmount{
			"value": {
				Currency: m.MarketID, // todo "xrp:usd" should this be split?
				Numeric:  m.Price.BigInt(),
				Text:     m.Price.String(),
			},
		},
		Additional: map[string][]string{
			"expiry": {m.Expiry.String()},
		},
	}, nil
}
