package mapper

import (
	"errors"
	"fmt"

	"github.com/figment-networks/indexing-engine/structs"
	"github.com/figment-networks/kava-worker/api/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/kava-labs/kava/app"
	"github.com/tendermint/tendermint/libs/bech32"
)

func BankMultisendToSub(msg sdk.Msg, logf types.LogFormat) (se structs.SubsetEvent, err error) {

	multisend, ok := msg.(bank.MsgMultiSend)
	if !ok {
		return se, errors.New("Not a multisend type")
	}

	se = structs.SubsetEvent{
		Type:   []string{"multisend"},
		Module: "bank",
	}
	for _, i := range multisend.Inputs {
		evt, err := bankProduceEvTx(i.Address, i.Coins)
		if err != nil {
			continue
		}
		se.Sender = append(se.Sender, evt)
	}

	for _, o := range multisend.Outputs {
		evt, err := bankProduceEvTx(o.Address, o.Coins)
		if err != nil {
			continue
		}
		se.Recipient = append(se.Recipient, evt)
	}

	err = produceTransfers(&se, "send", "", logf)
	return se, err
}

func BankSendToSub(msg sdk.Msg, logf types.LogFormat) (se structs.SubsetEvent, err error) {
	send, ok := msg.(bank.MsgSend)
	if !ok {
		return se, errors.New("Not a send type")
	}

	se = structs.SubsetEvent{
		Type:   []string{"send"},
		Module: "bank",
	}

	evt, _ := bankProduceEvTx(send.FromAddress, send.Amount)
	se.Sender = append(se.Sender, evt)

	evt, _ = bankProduceEvTx(send.ToAddress, send.Amount)
	se.Recipient = append(se.Recipient, evt)

	err = produceTransfers(&se, "send", "", logf)
	return se, err
}

func bankProduceEvTx(account sdk.AccAddress, coins sdk.Coins) (evt structs.EventTransfer, err error) {
	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, account.Bytes())
	if err != nil {
		return evt, fmt.Errorf("error converting Address: %w", err)
	}

	evt = structs.EventTransfer{
		Account: structs.Account{ID: bech32Addr},
	}
	if len(coins) > 0 {
		evt.Amounts = []structs.TransactionAmount{}
		for _, coin := range coins {
			evt.Amounts = append(evt.Amounts, structs.TransactionAmount{
				Currency: coin.Denom,
				Numeric:  coin.Amount.BigInt(),
				Text:     coin.Amount.String(),
			})
		}
	}

	return evt, nil
}
