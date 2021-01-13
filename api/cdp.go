package api

import (
	"errors"

	shared "github.com/figment-networks/indexer-manager/structs"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/x/cdp"
)

func mapCDPCreateCDPToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	m, ok := msg.(cdp.MsgCreateCDP)
	if !ok {
		return se, errors.New("Not a create_cdp type")
	}

	return shared.SubsetEvent{
		Type:   []string{"create_cdp"},
		Module: "cdp",
		Node: map[string][]shared.Account{
			"sender": {{ID: m.Sender.String()}},
		},
		Amount: map[string]shared.TransactionAmount{
			"collateral": {
				Currency: m.Collateral.Denom,
				Numeric:  m.Collateral.Amount.BigInt(),
				Text:     m.Collateral.String(),
			},
			"principal": {
				Currency: m.Principal.Denom,
				Numeric:  m.Principal.Amount.BigInt(),
				Text:     m.Principal.String(),
			},
		},
		Additional: map[string][]string{
			"collateral_type": []string{m.CollateralType},
		},
	}, nil
}

func mapCDPDepositCDPToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	m, ok := msg.(cdp.MsgDeposit)
	if !ok {
		return se, errors.New("Not a deposit_cdp type")
	}

	return shared.SubsetEvent{
		Type:   []string{"deposit_cdp"},
		Module: "cdp",
		Node: map[string][]shared.Account{
			"depositor": {{ID: m.Depositor.String()}},
			"owner":     {{ID: m.Owner.String()}},
		},
		Amount: map[string]shared.TransactionAmount{
			"collateral": {
				Currency: m.Collateral.Denom,
				Numeric:  m.Collateral.Amount.BigInt(),
				Text:     m.Collateral.String(),
			},
		},
		Additional: map[string][]string{
			"collateral_type": []string{m.CollateralType},
		},
	}, nil
}

func mapCDPWithdrawCDPToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	m, ok := msg.(cdp.MsgWithdraw)
	if !ok {
		return se, errors.New("Not a withdraw_cdp type")
	}

	return shared.SubsetEvent{
		Type:   []string{"withdraw_cdp"},
		Module: "cdp",
		Node: map[string][]shared.Account{
			"depositor": {{ID: m.Depositor.String()}},
			"owner":     {{ID: m.Owner.String()}},
		},
		Amount: map[string]shared.TransactionAmount{
			"collateral": {
				Currency: m.Collateral.Denom,
				Numeric:  m.Collateral.Amount.BigInt(),
				Text:     m.Collateral.String(),
			},
		},
		Additional: map[string][]string{
			"collateral_type": []string{m.CollateralType},
		},
	}, nil
}

func mapCDPDrawCDPToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	m, ok := msg.(cdp.MsgDrawDebt)
	if !ok {
		return se, errors.New("Not a draw_cdp type")
	}

	return shared.SubsetEvent{
		Type:   []string{"draw_cdp"},
		Module: "cdp",
		Node: map[string][]shared.Account{
			"sender": {{ID: m.Sender.String()}},
		},
		Amount: map[string]shared.TransactionAmount{
			"principal": {
				Currency: m.Principal.Denom,
				Numeric:  m.Principal.Amount.BigInt(),
				Text:     m.Principal.String(),
			},
		},
		Additional: map[string][]string{
			"collateral_type": []string{m.CollateralType},
		},
	}, nil
}

func mapCDPRepayCDPToSub(msg sdk.Msg) (se shared.SubsetEvent, err error) {
	m, ok := msg.(cdp.MsgRepayDebt)
	if !ok {
		return se, errors.New("Not a repay_cdp type")
	}

	return shared.SubsetEvent{
		Type:   []string{"repay_cdp"},
		Module: "cdp",
		Node: map[string][]shared.Account{
			"sender": {{ID: m.Sender.String()}},
		},
		Amount: map[string]shared.TransactionAmount{
			"payment": {
				Currency: m.Payment.Denom,
				Numeric:  m.Payment.Amount.BigInt(),
				Text:     m.Payment.String(),
			},
		},
		Additional: map[string][]string{
			"collateral_type": []string{m.CollateralType},
		},
	}, nil
}
