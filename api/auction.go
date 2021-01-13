package api

import (
	"errors"
	"strconv"

	shared "github.com/figment-networks/indexer-manager/structs"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/x/auction"
)

func mapAuctionPlaceBidToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	m, ok := msg.(auction.MsgPlaceBid)
	if !ok {
		return se, errors.New("Not a place_bid type")
	}

	se = shared.SubsetEvent{
		Type:   []string{"place_bid"},
		Module: "auction",
		Node: map[string][]shared.Account{
			"bidder": {{ID: m.Bidder.String()}},
		},
		Amount: map[string]shared.TransactionAmount{
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

	return se, nil
}
