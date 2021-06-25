package mapper

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/figment-networks/indexing-engine/structs"
	"github.com/figment-networks/kava-worker/api/types"
	"github.com/figment-networks/kava-worker/api/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/app"
)

var bech32ValPrefix string = app.Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixOperator

func produceTransfers(se *structs.SubsetEvent, transferType, skipAddr string, logf types.LogFormat) (err error) {
	var evts []structs.EventTransfer

	for _, ev := range logf.Events {
		if ev.Type != "transfer" {
			continue
		}

		var latestRecipient string
		for _, attr := range ev.Attributes {
			if len(attr.Recipient) > 0 {
				latestRecipient = attr.Recipient[0]
			}

			if latestRecipient == skipAddr {
				continue
			}

			amts := []structs.TransactionAmount{}

			if len(attr.Amount) == 0 {
				continue
			}

			for _, amount := range strings.Split(attr.Amount[0], ",") {
				attrAmt := structs.TransactionAmount{Numeric: &big.Int{}}
				sliced := util.GetCurrency(amount)
				var (
					c       *big.Int
					exp     int32
					coinErr error
				)
				if len(sliced) == 3 {
					attrAmt.Currency = sliced[2]
					c, exp, coinErr = util.GetCoin(sliced[1])
				} else {
					c, exp, coinErr = util.GetCoin(amount)
				}
				if coinErr != nil {
					return fmt.Errorf("[KAVA-API] Error parsing amount '%s': %s ", amount, coinErr)
				}

				attrAmt.Text = amount
				attrAmt.Exp = exp
				attrAmt.Numeric.Set(c)

				amts = append(amts, attrAmt)
			}
			evts = append(evts, structs.EventTransfer{
				Amounts: amts,
				Account: structs.Account{ID: latestRecipient},
			})
		}
	}

	if len(evts) <= 0 {
		return
	}

	if se.Transfers[transferType] == nil {
		se.Transfers = make(map[string][]structs.EventTransfer)
	}
	se.Transfers[transferType] = evts

	return
}
