package api

import (
	"errors"
	"strconv"

	shared "github.com/figment-networks/indexer-manager/structs"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/x/issuance"
)

func mapIssuanceIssueTokensToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	m, ok := msg.(issuance.MsgIssueTokens)
	if !ok {
		return se, errors.New("Not a issue_tokens type")
	}

	return shared.SubsetEvent{
		Type:   []string{"issue_tokens"},
		Module: "issuance",
		Node: map[string][]shared.Account{
			"sender":   {{ID: m.Sender.String()}},
			"receiver": {{ID: m.Receiver.String()}},
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

	return shared.SubsetEvent{
		Type:   []string{"redeem_tokens"},
		Module: "issuance",
		Node: map[string][]shared.Account{
			"sender": {{ID: m.Sender.String()}},
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

	return shared.SubsetEvent{
		Type:   []string{"block_address"},
		Module: "issuance",
		Node: map[string][]shared.Account{
			"sender":          {{ID: m.Sender.String()}},
			"blocked_address": {{ID: m.Address.String()}},
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

	return shared.SubsetEvent{
		Type:   []string{"unblock_address"},
		Module: "issuance",
		Node: map[string][]shared.Account{
			"sender":  {{ID: m.Sender.String()}},
			"address": {{ID: m.Address.String()}},
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

	return shared.SubsetEvent{
		Type:   []string{"change_pause_status"},
		Module: "issuance",
		Node: map[string][]shared.Account{
			"sender": {{ID: m.Sender.String()}},
		},
		Additional: map[string][]string{
			"denom":  []string{m.Denom},
			"status": []string{strconv.FormatBool(m.Status)},
		},
	}, nil
}
