package api

import (
	"errors"
	"fmt"
	"strconv"

	shared "github.com/figment-networks/indexer-manager/structs"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/app"
	"github.com/kava-labs/kava/x/issuance"
	"github.com/tendermint/tendermint/libs/bech32"
)

func mapIssuanceIssueTokensToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	m, ok := msg.(issuance.MsgIssueTokens)
	if !ok {
		return se, errors.New("Not a issue_tokens type")
	}

	bech32SenderAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Sender.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting SenderAddress: %w", err)
	}

	bech32RecAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Receiver.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting ReceiverAddress: %w", err)
	}

	return shared.SubsetEvent{
		Type:   []string{"issue_tokens"},
		Module: "issuance",
		Node: map[string][]shared.Account{
			"sender":   {{ID: bech32SenderAddr}},
			"receiver": {{ID: bech32RecAddr}},
		},
		Amount: map[string]shared.TransactionAmount{
			"send": {
				Currency: m.Tokens.Denom,
				Numeric:  m.Tokens.Amount.BigInt(),
				Text:     m.Tokens.String(),
			},
		},
	}, nil
}

func mapIssuanceRedeemTokensToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	m, ok := msg.(issuance.MsgRedeemTokens)
	if !ok {
		return se, errors.New("Not a redeem_tokens type")
	}

	bech32SenderAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Sender.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting SenderAddress: %w", err)
	}

	return shared.SubsetEvent{
		Type:   []string{"redeem_tokens"},
		Module: "issuance",
		Node: map[string][]shared.Account{
			"sender": {{ID: bech32SenderAddr}},
		},
		Amount: map[string]shared.TransactionAmount{
			"send": {
				Currency: m.Tokens.Denom,
				Numeric:  m.Tokens.Amount.BigInt(),
				Text:     m.Tokens.String(),
			},
		},
	}, nil
}

func mapIssuanceBlockAddressToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	m, ok := msg.(issuance.MsgBlockAddress)
	if !ok {
		return se, errors.New("Not a block_address type")
	}

	bech32SenderAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Sender.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting SenderAddress: %w", err)
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Address.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting Address: %w", err)
	}

	return shared.SubsetEvent{
		Type:   []string{"block_address"},
		Module: "issuance",
		Node: map[string][]shared.Account{
			"sender":          {{ID: bech32SenderAddr}},
			"blocked_address": {{ID: bech32Addr}},
		},
		Additional: map[string][]string{
			"denom": []string{m.Denom},
		},
	}, nil
}

func mapIssuanceUnblockAddressToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	m, ok := msg.(issuance.MsgBlockAddress)
	if !ok {
		return se, errors.New("Not a unblock_address type")
	}

	bech32SenderAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Sender.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting SenderAddress: %w", err)
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Address.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting Address: %w", err)
	}

	return shared.SubsetEvent{
		Type:   []string{"unblock_address"},
		Module: "issuance",
		Node: map[string][]shared.Account{
			"sender":  {{ID: bech32SenderAddr}},
			"address": {{ID: bech32Addr}},
		},
		Additional: map[string][]string{
			"denom": []string{m.Denom},
		},
	}, nil
}

func mapIssuanceMsgSetPauseStatusToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	m, ok := msg.(issuance.MsgSetPauseStatus)
	if !ok {
		return se, errors.New("Not a change_pause_status type")
	}

	bech32SenderAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Sender.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting SenderAddress: %w", err)
	}

	return shared.SubsetEvent{
		Type:   []string{"change_pause_status"},
		Module: "issuance",
		Node: map[string][]shared.Account{
			"sender": {{ID: bech32SenderAddr}},
		},
		Additional: map[string][]string{
			"denom":  []string{m.Denom},
			"status": []string{strconv.FormatBool(m.Status)},
		},
	}, nil
}
