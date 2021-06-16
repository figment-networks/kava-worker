package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/figment-networks/indexer-manager/structs"
	"github.com/figment-networks/kava-worker/api/types"
)

// BlocksMap map of blocks to control block map
// with extra summary of number of transactions
type BlocksMap struct {
	sync.Mutex
	Blocks map[uint64]structs.Block
	NumTxs uint64
}

// BlockErrorPair to wrap error response
type BlockErrorPair struct {
	Height uint64
	Block  structs.Block
	Err    error
}

// GetBlock fetches most recent block from chain
func (c Client) GetBlock(ctx context.Context, params structs.HeightHash) (block structs.Block, err error) {
	var ok bool
	if params.Height != 0 {
		block, ok = c.Sbc.Get(params.Height)
		if ok {
			return block, nil
		}
	}

	err = c.rateLimiter.Wait(ctx)
	if err != nil {
		return block, err
	}

	sCtx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	req, err := http.NewRequestWithContext(sCtx, http.MethodGet, c.baseURL+"/block", nil)
	if err != nil {
		return block, err
	}

	req.Header.Add("Content-Type", "application/json")
	if c.key != "" {
		req.Header.Add("Authorization", c.key)
	}

	q := req.URL.Query()
	if params.Height > 0 {
		q.Add("height", strconv.FormatUint(params.Height, 10))
	}
	req.URL.RawQuery = q.Encode()

	n := time.Now()
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return block, err
	}
	rawRequestHTTPDuration.WithLabels("/block", resp.Status).Observe(time.Since(n).Seconds())
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	var result *types.GetBlockResponse
	if err = decoder.Decode(&result); err != nil {
		return block, err
	}

	if result.Error.Message != "" {
		return block, fmt.Errorf("[KAVA-API] Error fetching block: %s ", result.Error.Message)
	}
	bTime, err := time.Parse(time.RFC3339Nano, result.Result.Block.Header.Time)
	if err != nil {
		return block, err
	}
	uHeight, err := strconv.ParseUint(result.Result.Block.Header.Height, 10, 64)
	if err != nil {
		return block, err
	}

	numTxs := len(result.Result.Block.Data.Txs)

	block = structs.Block{
		Hash:                 result.Result.BlockID.Hash,
		Height:               uHeight,
		Time:                 bTime,
		ChainID:              result.Result.Block.Header.ChainID,
		NumberOfTransactions: uint64(numTxs),
	}

	c.Sbc.Add(block)
	return block, nil
}
