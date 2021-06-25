package mapper

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/figment-networks/indexing-engine/structs"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/app"
	"github.com/kava-labs/kava/x/issuance"
	"github.com/tendermint/tendermint/libs/bech32"
)

func IssuanceIssueTokensToSub(msg sdk.Msg) (se structs.SubsetEvent, err error) {
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

	return structs.SubsetEvent{
		Type:   []string{"issue_tokens"},
		Module: "issuance",
		Node: map[string][]structs.Account{
			"sender":   {{ID: bech32SenderAddr}},
			"receiver": {{ID: bech32RecAddr}},
		},
		Amount: map[string]structs.TransactionAmount{
			"send": {
				Currency: m.Tokens.Denom,
				Numeric:  m.Tokens.Amount.BigInt(),
				Text:     m.Tokens.String(),
			},
		},
	}, nil
}

func IssuanceRedeemTokensToSub(msg sdk.Msg) (se structs.SubsetEvent, err error) {
	m, ok := msg.(issuance.MsgRedeemTokens)
	if !ok {
		return se, errors.New("Not a redeem_tokens type")
	}

	bech32SenderAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Sender.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting SenderAddress: %w", err)
	}

	return structs.SubsetEvent{
		Type:   []string{"redeem_tokens"},
		Module: "issuance",
		Node: map[string][]structs.Account{
			"sender": {{ID: bech32SenderAddr}},
		},
		Amount: map[string]structs.TransactionAmount{
			"send": {
				Currency: m.Tokens.Denom,
				Numeric:  m.Tokens.Amount.BigInt(),
				Text:     m.Tokens.String(),
			},
		},
	}, nil
}

func IssuanceBlockAddressToSub(msg sdk.Msg) (se structs.SubsetEvent, err error) {
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

	return structs.SubsetEvent{
		Type:   []string{"block_address"},
		Module: "issuance",
		Node: map[string][]structs.Account{
			"sender":          {{ID: bech32SenderAddr}},
			"blocked_address": {{ID: bech32Addr}},
		},
		Additional: map[string][]string{
			"denom": []string{m.Denom},
		},
	}, nil
}

func IssuanceUnblockAddressToSub(msg sdk.Msg) (se structs.SubsetEvent, err error) {
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

	return structs.SubsetEvent{
		Type:   []string{"unblock_address"},
		Module: "issuance",
		Node: map[string][]structs.Account{
			"sender":  {{ID: bech32SenderAddr}},
			"address": {{ID: bech32Addr}},
		},
		Additional: map[string][]string{
			"denom": []string{m.Denom},
		},
	}, nil
}

func IssuanceMsgSetPauseStatusToSub(msg sdk.Msg) (se structs.SubsetEvent, err error) {
	m, ok := msg.(issuance.MsgSetPauseStatus)
	if !ok {
		return se, errors.New("Not a change_pause_status type")
	}

	bech32SenderAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Sender.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting SenderAddress: %w", err)
	}

	return structs.SubsetEvent{
		Type:   []string{"change_pause_status"},
		Module: "issuance",
		Node: map[string][]structs.Account{
			"sender": {{ID: bech32SenderAddr}},
		},
		Additional: map[string][]string{
			"denom":  []string{m.Denom},
			"status": []string{strconv.FormatBool(m.Status)},
		},
	}, nil
}
