package client

import (
	"bytes"
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/figment-networks/indexing-engine/metrics"

	"github.com/figment-networks/indexer-manager/structs"
	"github.com/google/uuid"
	"go.uber.org/zap"

	cStructs "github.com/figment-networks/indexer-manager/worker/connectivity/structs"
	"github.com/figment-networks/indexer-manager/worker/process/ranged"
	"github.com/figment-networks/indexer-manager/worker/store"
	"github.com/figment-networks/kava-worker/api"
)

const page = 100
const blockchainEndpointLimit = 20

var (
	getTransactionDuration *metrics.GroupObserver
	getLatestDuration      *metrics.GroupObserver
	getBlockDuration       *metrics.GroupObserver
)

type OutputSender interface {
	Send(cStructs.TaskResponse) error
}

type RPC interface {
	GetBlock(ctx context.Context, params structs.HeightHash) (block structs.Block, err error)
	SearchTx(ctx context.Context, r structs.HeightHash, block structs.Block, perPage uint64) (txs []structs.Transaction, err error)
}

type LCD interface {
	GetReward(ctx context.Context, params structs.HeightAccount) (resp structs.GetRewardResponse, err error)
}

// IndexerClient is implementation of a client (main worker code)
type IndexerClient struct {
	rpcCli RPC
	lcdCli LCD

	logger  *zap.Logger
	streams map[uuid.UUID]*cStructs.StreamAccess
	sLock   sync.Mutex

	Reqester            *ranged.RangeRequester
	storeClient         store.HeightStoreCaller
	maximumHeightsToGet uint64
}

// NewIndexerClient is IndexerClient constructor
func NewIndexerClient(ctx context.Context, logger *zap.Logger, rpcCli RPC, lcdCli LCD, storeClient store.HeightStoreCaller, maximumHeightsToGet uint64) *IndexerClient {
	getTransactionDuration = endpointDuration.WithLabels("getTransactions")
	getLatestDuration = endpointDuration.WithLabels("getLatest")
	getBlockDuration = endpointDuration.WithLabels("getBlock")
	api.InitMetrics()

	ic := &IndexerClient{
		logger:              logger,
		rpcCli:              rpcCli,
		lcdCli:              lcdCli,
		storeClient:         storeClient,
		maximumHeightsToGet: maximumHeightsToGet,
		streams:             make(map[uuid.UUID]*cStructs.StreamAccess),
	}

	ic.Reqester = ranged.NewRangeRequester(ic, 20)
	return ic
}

// CloseStream removes stream from worker/client
func (ic *IndexerClient) CloseStream(ctx context.Context, streamID uuid.UUID) error {
	ic.sLock.Lock()
	defer ic.sLock.Unlock()

	ic.logger.Debug("[KAVA-CLIENT] Close Stream", zap.Stringer("streamID", streamID))
	delete(ic.streams, streamID)

	return nil
}

// RegisterStream adds new listeners to the streams - currently fixed number per stream
func (ic *IndexerClient) RegisterStream(ctx context.Context, stream *cStructs.StreamAccess) error {
	ic.logger.Debug("[KAVA-CLIENT] Register Stream", zap.Stringer("streamID", stream.StreamID))
	newStreamsMetric.WithLabels().Inc()

	ic.sLock.Lock()
	defer ic.sLock.Unlock()
	ic.streams[stream.StreamID] = stream

	// Limit workers not to create new goroutines over and over again
	for i := 0; i < 20; i++ {
		go ic.Run(ctx, stream)
	}

	return nil
}

// Run listens on the stream events (new tasks)
func (ic *IndexerClient) Run(ctx context.Context, stream *cStructs.StreamAccess) {
	for {
		select {
		case <-ctx.Done():
			ic.sLock.Lock()
			delete(ic.streams, stream.StreamID)
			ic.sLock.Unlock()
			return
		case <-stream.Finish:
			return
		case taskRequest := <-stream.RequestListener:
			receivedRequestsMetric.WithLabels(taskRequest.Type).Inc()
			nCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
			switch taskRequest.Type {
			case structs.ReqIDGetTransactions:
				ic.GetTransactions(nCtx, taskRequest, stream, ic.rpcCli)
			case structs.ReqIDGetLatestMark:
				ic.GetLatestMark(nCtx, taskRequest, stream, ic.rpcCli)
			case structs.ReqIDGetReward:
				ic.GetReward(nCtx, taskRequest, stream, ic.lcdCli)
			default:
				stream.Send(cStructs.TaskResponse{
					Id:    taskRequest.Id,
					Error: cStructs.TaskError{Msg: "There is no such handler " + taskRequest.Type},
					Final: true,
				})
			}
			cancel()
		}
	}
}

// GetBlock gets block
func (ic *IndexerClient) GetBlock(ctx context.Context, tr cStructs.TaskRequest, stream *cStructs.StreamAccess, client RPC) {
	timer := metrics.NewTimer(getBlockDuration)
	defer timer.ObserveDuration()

	hr := &structs.HeightHash{}
	err := json.Unmarshal(tr.Payload, hr)
	if err != nil {
		stream.Send(cStructs.TaskResponse{
			Id:    tr.Id,
			Error: cStructs.TaskError{Msg: "Cannot unmarshal payload"},
			Final: true,
		})
		return
	}

	block, err := client.GetBlock(ctx, *hr)
	if err != nil {
		ic.logger.Error("Error getting block", zap.Error(err))
		stream.Send(cStructs.TaskResponse{
			Id:    tr.Id,
			Error: cStructs.TaskError{Msg: "Error getting block data " + err.Error()},
			Final: true,
		})
		return
	}

	out := make(chan cStructs.OutResp, 1)
	out <- cStructs.OutResp{
		ID:      tr.Id,
		Type:    "Block",
		Payload: block,
	}
	close(out)

	sendResp(ctx, tr.Id, out, ic.logger, stream, nil)
}

// GetReward gets reward
func (ic *IndexerClient) GetReward(ctx context.Context, tr cStructs.TaskRequest, stream *cStructs.StreamAccess, client LCD) {
	timer := metrics.NewTimer(getBlockDuration)
	defer timer.ObserveDuration()

	ha := &structs.HeightAccount{}
	err := json.Unmarshal(tr.Payload, ha)
	if err != nil {
		stream.Send(cStructs.TaskResponse{
			Id:    tr.Id,
			Error: cStructs.TaskError{Msg: "Cannot unmarshal payload"},
			Final: true,
		})
		return
	}

	reward, err := client.GetReward(ctx, *ha)
	if err != nil {
		ic.logger.Error("Error getting reward", zap.Error(err))
		stream.Send(cStructs.TaskResponse{
			Id:    tr.Id,
			Error: cStructs.TaskError{Msg: "Error getting reward data " + err.Error()},
			Final: true,
		})
		return
	}

	out := make(chan cStructs.OutResp, 1)
	out <- cStructs.OutResp{
		ID:      tr.Id,
		Type:    "Reward",
		Payload: reward,
	}
	close(out)

	sendResp(ctx, tr.Id, out, ic.logger, stream, nil)
}

// sendResp constructs protocol response and send it out to transport
func sendResp(ctx context.Context, id uuid.UUID, in <-chan cStructs.OutResp, logger *zap.Logger, stream *cStructs.StreamAccess, fin chan bool) {
	b := &bytes.Buffer{}
	enc := json.NewEncoder(b)
	order := uint64(0)

	var contextDone bool

SendLoop:
	for {
		select {
		case <-ctx.Done():
			contextDone = true
			break SendLoop
		case t, ok := <-in:
			if !ok && t.Type == "" {
				break SendLoop
			}
			b.Reset()

			err := enc.Encode(t.Payload)
			if err != nil {
				logger.Error("[KAVA-CLIENT] Error encoding payload data", zap.Error(err))
			}

			tr := cStructs.TaskResponse{
				Id:      id,
				Type:    t.Type,
				Order:   order,
				Payload: make([]byte, b.Len()),
			}

			b.Read(tr.Payload)
			order++
			err = stream.Send(tr)
			if err != nil {
				logger.Error("[KAVA-CLIENT] Error sending data", zap.Error(err))
			}
			sendResponseMetric.WithLabels(t.Type, "yes").Inc()
		}
	}

	err := stream.Send(cStructs.TaskResponse{
		Id:    id,
		Type:  "END",
		Order: order,
		Final: true,
	})

	if err != nil {
		logger.Error("[KAVA-CLIENT] Error sending end", zap.Error(err))
	}

	if fin != nil {
		if !contextDone {
			fin <- true
		}
		close(fin)
	}

}
