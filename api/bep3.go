package api

import (
	"errors"
	"strconv"

	shared "github.com/figment-networks/indexer-manager/structs"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/x/bep3"
)

func mapBep3CreateAtomicSwapToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	m, ok := msg.(bep3.MsgCreateAtomicSwap)
	if !ok {
		return se, errors.New("Not a createAtomicSwap type")
	}

	se = shared.SubsetEvent{
		Type:   []string{"create_atomic_swap"},
		Module: "bep3",
		Node: map[string][]shared.Account{
			"from": {{ID: m.From.String()}},
			"to":   {{ID: m.To.String()}},
		},
		Additional: map[string][]string{
			"recipient_other_chain": []string{m.RecipientOtherChain},
			"sender_other_chain":    []string{m.SenderOtherChain},
			"random_number_hash":    []string{m.RandomNumberHash.String()},
			"timestamp":             []string{strconv.FormatInt(m.Timestamp, 10)},
			"height_span":           []string{strconv.FormatUint(m.HeightSpan, 10)},
		},
	}

	amts := make([]shared.TransactionAmount, len(m.Amount))
	for i, amt := range m.Amount {
		amts[i] = shared.TransactionAmount{
			Currency: amt.Denom,
			Numeric:  amt.Amount.BigInt(),
			Text:     amt.String(),
		}
	}

	se.Transfers = map[string][]shared.EventTransfer{
		"send": []shared.EventTransfer{
			{Amounts: amts},
		},
	}

	return se, nil
}

func mapBep3ClaimAtomicSwapToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	m, ok := msg.(bep3.MsgClaimAtomicSwap)
	if !ok {
		return se, errors.New("Not a claimAtomicSwap type")
	}
	// todo get amount (from logs?)

	return shared.SubsetEvent{
		Type:   []string{"claim_atomic_swap"},
		Module: "bep3",
		Node: map[string][]shared.Account{
			"from": {{ID: m.From.String()}},
		},
		Additional: map[string][]string{
			"swap_id":       []string{m.SwapID.String()},
			"random_number": []string{m.RandomNumber.String()},
		},
	}, nil
}

func mapBep3RefundAtomicSwapToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	m, ok := msg.(bep3.MsgRefundAtomicSwap)
	if !ok {
		return se, errors.New("Not a refundAtomicSwap type")
	}

	return shared.SubsetEvent{
		Type:   []string{"refund_atomic_swap"},
		Module: "bep3",
		Node: map[string][]shared.Account{
			"from": {{ID: m.From.String()}},
		},
		Amount: map[string]shared.TransactionAmount{},
		Additional: map[string][]string{
			"swap_id": []string{m.SwapID.String()},
		},
	}, nil
}
