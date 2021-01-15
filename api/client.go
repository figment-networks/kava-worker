package api

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/kava-labs/kava/app"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// Client is a Tendermint RPC client for cosmos using figmentnetworks datahub
type Client struct {
	baseURL    string
	key        string
	httpClient *http.Client
	cdc        *codec.Codec
	logger     *zap.Logger

	rateLimiter *rate.Limiter
	Sbc         *SimpleBlockCache
	CallMap     sync.Map
}

// NewClient returns a new client for a given endpoint
func NewClient(url, key string, logger *zap.Logger, c *http.Client, reqPerSecLimit int) *Client {
	fmt.Println("[NewClient]")

	if c == nil {
		c = &http.Client{
			Timeout: time.Second * 10,
		}
	}

	rateLimiter := rate.NewLimiter(rate.Limit(reqPerSecLimit), reqPerSecLimit)

	cli := &Client{
		logger:      logger,
		baseURL:     url, //tendermint rpc url
		key:         key,
		httpClient:  c,
		rateLimiter: rateLimiter,
		cdc:         app.MakeCodec(),
		Sbc:         NewSimpleBlockCache(400),
	}
	return cli
}

// InitMetrics initialise metrics
func InitMetrics() {
	convertionDurationObserver = conversionDuration.WithLabels("conversion")
	transactionConversionDuration = conversionDuration.WithLabels("transaction")
	blockCacheEfficiencyHit = blockCacheEfficiency.WithLabels("hit")
	blockCacheEfficiencyMissed = blockCacheEfficiency.WithLabels("missed")
	numberOfItemsTransactions = numberOfItems.WithLabels("transactions")
}
