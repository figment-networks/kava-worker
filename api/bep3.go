package api

import (
	"errors"
	"fmt"
	"strconv"

	shared "github.com/figment-networks/indexer-manager/structs"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/app"
	"github.com/kava-labs/kava/x/bep3"
	"github.com/tendermint/tendermint/libs/bech32"
)

func mapBep3CreateAtomicSwapToSub(msg sdk.Msg, logf LogFormat) (se shared.SubsetEvent, err error) {
	m, ok := msg.(bep3.MsgCreateAtomicSwap)
	if !ok {
		return se, errors.New("Not a createAtomicSwap type")
	}

	bech32FromAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.From.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting FromAddress: %w", err)
	}
	bech32ToAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.To.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting ToAddress: %w", err)
	}

	se = shared.SubsetEvent{
		Type:   []string{"create_atomic_swap"},
		Module: "bep3",
		Node: map[string][]shared.Account{
			"from": {{ID: bech32FromAddr}},
			"to":   {{ID: bech32ToAddr}},
		},
		Additional: map[string][]string{
			"recipient_other_chain": []string{m.RecipientOtherChain},
			"sender_other_chain":    []string{m.SenderOtherChain},
			"random_number_hash":    []string{m.RandomNumberHash.String()},
			"timestamp":             []string{strconv.FormatInt(m.Timestamp, 10)},
			"height_span":           []string{strconv.FormatUint(m.HeightSpan, 10)},
		},
	}

	txAmount := map[string]shared.TransactionAmount{}

	for i, coin := range m.Amount {
		am := shared.TransactionAmount{
			Currency: coin.Denom,
			Numeric:  coin.Amount.BigInt(),
			Text:     coin.Amount.String(),
		}

		key := "send"
		if i > 0 {
			key += "_" + strconv.Itoa(i)
		}

		txAmount[key] = am
	}

	se.Amount = txAmount
	err = produceTransfers(&se, "send", logf)
	return se, nil
}

func mapBep3ClaimAtomicSwapToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	m, ok := msg.(bep3.MsgClaimAtomicSwap)
	if !ok {
		return se, errors.New("Not a claimAtomicSwap type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.From.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting FromAddress: %w", err)
	}

	return shared.SubsetEvent{
		Type:   []string{"claim_atomic_swap"},
		Module: "bep3",
		Node: map[string][]shared.Account{
			"from": {{ID: bech32Addr}},
		},
		Additional: map[string][]string{
			"swap_id":       []string{m.SwapID.String()},
			"random_number": []string{m.RandomNumber.String()},
		},
	}, nil
}

func mapBep3RefundAtomicSwapToSub(msg sdk.Msg, logf LogFormat) (se shared.SubsetEvent, err error) {
	m, ok := msg.(bep3.MsgRefundAtomicSwap)
	if !ok {
		return se, errors.New("Not a refundAtomicSwap type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.From.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting FromAddress: %w", err)
	}

	se = shared.SubsetEvent{
		Type:   []string{"refund_atomic_swap"},
		Module: "bep3",
		Node: map[string][]shared.Account{
			"from": {{ID: bech32Addr}},
		},
		Additional: map[string][]string{
			"swap_id": []string{m.SwapID.String()},
		},
	}

	err = produceTransfers(&se, "send", logf)
	return se, err
}
