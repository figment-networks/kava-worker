package mapper

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/figment-networks/indexer-search/structs"

	sdk "github.com/cosmos/cosmos-sdk/types"
	evidence "github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/kava-labs/kava/app"
	"github.com/tendermint/tendermint/libs/bech32"
)

func EvidenceSubmitEvidenceToSub(msg sdk.Msg) (se structs.SubsetEvent, er error) {
	mse, ok := msg.(evidence.MsgSubmitEvidence)
	if !ok {
		return se, errors.New("Not a submit_evidence type")
	}

	bech32Addr, err := bech32.ConvertAndEncode(app.Bech32MainPrefix, mse.Submitter.Bytes())
	if err != nil {
		return se, fmt.Errorf("error converting SubmitterAddress: %w", err)
	}

	return structs.SubsetEvent{
		Type:   []string{"submit_evidence"},
		Module: "evidence",
		Node:   map[string][]structs.Account{"submitter": {{ID: bech32Addr}}},
		Additional: map[string][]string{
			"evidence_consensus":       {mse.Evidence.GetConsensusAddress().String()},
			"evidence_height":          {strconv.FormatInt(mse.Evidence.GetHeight(), 10)},
			"evidence_total_power":     {strconv.FormatInt(mse.Evidence.GetTotalPower(), 10)},
			"evidence_validator_power": {strconv.FormatInt(mse.Evidence.GetValidatorPower(), 10)},
		},
	}, nil
}
