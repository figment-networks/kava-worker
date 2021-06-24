package mapper

import (
	"errors"
	"fmt"

	"github.com/figment-networks/indexer-search/structs"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/app"
	"github.com/kava-labs/kava/x/cdp"
	"github.com/tendermint/tendermint/libs/bech32"
)

func CDPCreateCDPToSub(msg sdk.Msg) (se structs.SubsetEvent, err error) {
	m, ok := msg.(cdp.MsgCreateCDP)
	if !ok {
		return se, errors.New("Not a create_cdp type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Sender.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting SenderAddress: %w", err)
	}

	return structs.SubsetEvent{
		Type:   []string{"create_cdp"},
		Module: "cdp",
		Node: map[string][]structs.Account{
			"sender": {{ID: bech32Addr}},
		},
		Amount: map[string]structs.TransactionAmount{
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

func CDPDepositCDPToSub(msg sdk.Msg) (se structs.SubsetEvent, err error) {
	m, ok := msg.(cdp.MsgDeposit)
	if !ok {
		return se, errors.New("Not a deposit_cdp type")
	}

	bech32DepositorAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Depositor.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting DepositorAddress: %w", err)
	}

	bech32OwnerAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Owner.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting OwnerAddress: %w", err)
	}

	return structs.SubsetEvent{
		Type:   []string{"deposit_cdp"},
		Module: "cdp",
		Node: map[string][]structs.Account{
			"depositor": {{ID: bech32DepositorAddr}},
			"owner":     {{ID: bech32OwnerAddr}},
		},
		Amount: map[string]structs.TransactionAmount{
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

func CDPWithdrawCDPToSub(msg sdk.Msg) (se structs.SubsetEvent, err error) {
	m, ok := msg.(cdp.MsgWithdraw)
	if !ok {
		return se, errors.New("Not a withdraw_cdp type")
	}

	bech32DepositorAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Depositor.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting DepositerAddress: %w", err)
	}

	bech32OwnerAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Owner.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting OwnerAddress: %w", err)
	}

	return structs.SubsetEvent{
		Type:   []string{"withdraw_cdp"},
		Module: "cdp",
		Node: map[string][]structs.Account{
			"depositor": {{ID: bech32DepositorAddr}},
			"owner":     {{ID: bech32OwnerAddr}},
		},
		Amount: map[string]structs.TransactionAmount{
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

func CDPDrawCDPToSub(msg sdk.Msg) (se structs.SubsetEvent, err error) {
	m, ok := msg.(cdp.MsgDrawDebt)
	if !ok {
		return se, errors.New("Not a draw_cdp type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Sender.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting SenderAddress: %w", err)
	}

	return structs.SubsetEvent{
		Type:   []string{"draw_cdp"},
		Module: "cdp",
		Node: map[string][]structs.Account{
			"sender": {{ID: bech32Addr}},
		},
		Amount: map[string]structs.TransactionAmount{
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

func CDPRepayCDPToSub(msg sdk.Msg) (se structs.SubsetEvent, err error) {
	m, ok := msg.(cdp.MsgRepayDebt)
	if !ok {
		return se, errors.New("Not a repay_cdp type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Sender.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting SenderAddress: %w", err)
	}

	return structs.SubsetEvent{
		Type:   []string{"repay_cdp"},
		Module: "cdp",
		Node: map[string][]structs.Account{
			"sender": {{ID: bech32Addr}},
		},
		Amount: map[string]structs.TransactionAmount{
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

func CDPLiquidateToSub(msg sdk.Msg) (se structs.SubsetEvent, err error) {
	m, ok := msg.(cdp.MsgLiquidate)
	if !ok {
		return se, errors.New("Not a liquidate type")
	}

	bech32KeeperAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Keeper.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting KeeperAddress: %w", err)
	}

	bech32BorrowerAddr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, m.Borrower.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting BorrowerAddress: %w", err)
	}

	return structs.SubsetEvent{
		Type:   []string{"liquidate"},
		Module: "cdp",
		Node: map[string][]structs.Account{
			"keeper":   {{ID: bech32KeeperAddr}},
			"borrower": {{ID: bech32BorrowerAddr}},
		},
		Additional: map[string][]string{
			"collateral_type": []string{m.CollateralType},
		},
	}, nil
}
