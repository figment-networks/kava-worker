package client

import (
	"context"
	"encoding/json"
	"time"

	mStructs "github.com/figment-networks/indexer-manager/structs"
	cStructs "github.com/figment-networks/indexer-manager/worker/connectivity/structs"
	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/indexing-engine/structs"
	"go.uber.org/zap"
)

// GetLatest gets latest transactions and blocks.
// It gets latest transaction, then diff it with
func (ic *IndexerClient) GetLatestMark(ctx context.Context, tr cStructs.TaskRequest, stream *cStructs.StreamAccess, client RPC) {
	timer := metrics.NewTimer(getLatestDuration)
	defer timer.ObserveDuration()

	ldr := &mStructs.LatestDataRequest{}
	err := json.Unmarshal(tr.Payload, ldr)
	if err != nil {
		stream.Send(cStructs.TaskResponse{Id: tr.Id, Error: cStructs.TaskError{Msg: "Cannot unmarshal payload"}, Final: true})
	}

	sCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	// (lukanus): Get latest block (height = 0)
	block, err := client.GetBlock(sCtx, structs.HeightHash{})
	if err != nil {
		stream.Send(cStructs.TaskResponse{Id: tr.Id, Error: cStructs.TaskError{Msg: "Error getting block data " + err.Error()}, Final: true})
		return
	}

	tResp := cStructs.TaskResponse{Id: tr.Id, Type: "LatestMark", Order: 0, Final: true}
	tResp.Payload, err = json.Marshal(mStructs.LatestDataResponse{
		LastHash:   block.Hash,
		LastTime:   block.Time,
		LastHeight: block.Height,
		LastEpoch:  block.Epoch,
	})

	if err != nil {
		ic.logger.Error("[KAVA-CLIENT] Error encoding payload data", zap.Error(err))
	}

	if err := stream.Send(tResp); err != nil {
		ic.logger.Error("[KAVA-CLIENT] Error sending end", zap.Error(err))
	}
}
