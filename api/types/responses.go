package types

// TxResponse is result of querying for a tx
type TxResponse struct {
	Hash   string  `json:"hash"`
	Height string  `json:"height"`
	Index  float64 `json:"index"`

	TxResult ResponseDeliverTx `json:"tx_result"`
	// TxData is base64 encoded transaction data
	TxData string `json:"tx"`

	All int64
}

// ResponseDeliverTx result
type ResponseDeliverTx struct {
	Log       string  `json:"log"`
	Code      float64 `json:"code"`
	Codespace string  `json:"codespace"`
	GasWanted string  `json:"gasWanted"`
	GasUsed   string  `json:"gasUsed"`
	Tags      []TxTag `json:"tags"`
}

// TxTag is tag from chain
type TxTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ResultBlock is result of fetching block
type ResultBlock struct {
	Block   Block   `json:"block"`
	BlockID BlockID `json:"block_id"`
}

// ResultBlockchain is result of fetching block
type ResultBlockchain struct {
	LastHeight string      `json:"last_height"`
	BlockMetas []BlockMeta `json:"block_metas"`
}

// BlockMeta info
type BlockMeta struct {
	BlockID BlockID     `json:"block_id"`
	Header  BlockHeader `json:"header"`
	NumTxs  string      `json:"num_txs"`
}

// BlockID info
type BlockID struct {
	Hash string `json:"hash"`
}

// Block is kava block data
type Block struct {
	Header BlockHeader `json:"header"`
	Data   BlockData   `json:"data"`
}

// BlockHeader structures
type BlockHeader struct {
	Height  string `json:"height"`
	ChainID string `json:"chain_id"`
	Time    string `json:"time"`
}

// BlockData structures
type BlockData struct {
	Txs []string `json:"txs"`
}

// Error is api error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

// ResultTxSearch of searching for txs
type ResultTxSearch struct {
	Txs        []TxResponse `json:"txs"`
	TotalCount string       `json:"total_count"`
}

// GetTxSearchResponse cosmos response for search
type GetTxSearchResponse struct {
	// ID     string         `json:"id"`
	RPC    string         `json:"jsonrpc"`
	Result ResultTxSearch `json:"result"`
	Error  Error          `json:"error"`
}

// GetBlockResponse cosmos response from block
type GetBlockResponse struct {
	// ID     string      `json:"id"`
	RPC    string      `json:"jsonrpc"`
	Result ResultBlock `json:"result"`
	Error  Error       `json:"error"`
}

// GetBlockchainResponse cosmos response from blockchain
type GetBlockchainResponse struct {
	//ID     string           `json:"id"`
	RPC    string           `json:"jsonrpc"`
	Result ResultBlockchain `json:"result"`
	Error  Error            `json:"error"`
}
