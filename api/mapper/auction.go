package mapper

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/figment-networks/indexing-engine/structs"
	"github.com/figment-networks/kava-worker/api/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/app"
	"github.com/kava-labs/kava/x/auction"
	"github.com/tendermint/tendermint/libs/bech32"
)

func AuctionPlaceBidToSub(msg sdk.Msg, logf types.LogFormat) (se structs.SubsetEvent, err error) {
	m, ok := msg.(auction.MsgPlaceBid)
	if !ok {
		return se, errors.New("Not a place_bid type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Bidder.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting Address: %w", err)
	}

	se = structs.SubsetEvent{
		Type:   []string{"place_bid"},
		Module: "auction",
		Node: map[string][]structs.Account{
			"bidder": {{ID: bech32Addr}},
		},
		Amount: map[string]structs.TransactionAmount{
			"bid": {
				Currency: m.Amount.Denom,
				Numeric:  m.Amount.Amount.BigInt(),
				Text:     m.Amount.String(),
			},
		},
		Additional: map[string][]string{
			"auction_id": []string{strconv.FormatUint(m.AuctionID, 10)},
		},
	}

	err = produceTransfers(&se, "send", "", logf)
	return se, nil
}
