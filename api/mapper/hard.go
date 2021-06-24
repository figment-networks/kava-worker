package mapper

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/figment-networks/indexer-search/structs"
	"github.com/figment-networks/kava-worker/api/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/app"
	"github.com/kava-labs/kava/x/hard"
	"github.com/tendermint/tendermint/libs/bech32"
)

func HardDepositToSub(msg sdk.Msg, logf types.LogFormat) (se structs.SubsetEvent, err error) {
	m, ok := msg.(hard.MsgDeposit)
	if !ok {
		return se, errors.New("Not a hard_deposit type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Depositor.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting Depositor address: %w", err)
	}

	se = structs.SubsetEvent{
		Type:   []string{"hard_deposit"},
		Module: "hard",
		Node: map[string][]structs.Account{
			"depositor": {{ID: bech32Addr}},
		},
		Amount: hardProduceAmounts(m.Amount),
	}

	err = produceTransfers(&se, "send", "", logf)
	return se, err
}

func HardWithdrawToSub(msg sdk.Msg, logf types.LogFormat) (se structs.SubsetEvent, err error) {
	m, ok := msg.(hard.MsgWithdraw)
	if !ok {
		return se, errors.New("Not a hard_withdraw type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Depositor.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting Depositor Address: %w", err)
	}

	se = structs.SubsetEvent{
		Type:   []string{"hard_withdraw"},
		Module: "hard",
		Node: map[string][]structs.Account{
			"depositor": {{ID: bech32Addr}},
		},
		Amount: hardProduceAmounts(m.Amount),
	}

	err = produceTransfers(&se, "send", "", logf)
	return se, err
}

func HardBorrowToSub(msg sdk.Msg, logf types.LogFormat) (se structs.SubsetEvent, err error) {
	m, ok := msg.(hard.MsgBorrow)
	if !ok {
		return se, errors.New("Not a hard_borrow type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Borrower.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting Borrower address: %w", err)
	}

	se = structs.SubsetEvent{
		Type:   []string{"hard_borrow"},
		Module: "hard",
		Node: map[string][]structs.Account{
			"borrower": {{ID: bech32Addr}},
		},
		Amount: hardProduceAmounts(m.Amount),
	}

	err = produceTransfers(&se, "send", "", logf)
	return se, err
}

func HardRepayToSub(msg sdk.Msg, logf types.LogFormat) (se structs.SubsetEvent, err error) {
	m, ok := msg.(hard.MsgRepay)
	if !ok {
		return se, errors.New("Not a hard_repay type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Sender.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting Sender address: %w", err)
	}

	bech32OwnerAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Owner.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting Owner address: %w", err)
	}

	se = structs.SubsetEvent{
		Type:   []string{"hard_repay"},
		Module: "hard",
		Node: map[string][]structs.Account{
			"sender": {{ID: bech32Addr}},
			"owner":  {{ID: bech32OwnerAddr}},
		},
		Amount: hardProduceAmounts(m.Amount),
	}

	err = produceTransfers(&se, "send", "", logf)
	return se, err
}

func HardLiquidateToSub(msg sdk.Msg, logf types.LogFormat) (se structs.SubsetEvent, err error) {
	m, ok := msg.(hard.MsgLiquidate)
	if !ok {
		return se, errors.New("Not a hard_liquidate type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Keeper.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting Keeper address: %w", err)
	}

	bech32BorrowerAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Borrower.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting Borrower address: %w", err)
	}

	se = structs.SubsetEvent{
		Type:   []string{"hard_liquidate"},
		Module: "hard",
		Node: map[string][]structs.Account{
			"keeper":   {{ID: bech32Addr}},
			"borrower": {{ID: bech32BorrowerAddr}},
		},
	}

	err = produceTransfers(&se, "send", "", logf)
	return se, err
}

func hardProduceAmounts(coins sdk.Coins) map[string]structs.TransactionAmount {

	if len(coins) > 0 {
		txAm := make(map[string]structs.TransactionAmount)
		for i, coin := range coins {
			txAm[strconv.Itoa(i)] = structs.TransactionAmount{
				Currency: coin.Denom,
				Numeric:  coin.Amount.BigInt(),
				Text:     coin.Amount.String(),
			}
		}
		return txAm
	}

	return nil
}
