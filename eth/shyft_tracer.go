package eth

import (
	"context"

	"github.com/ShyftNetwork/go-empyrean/common"
	"github.com/ShyftNetwork/go-empyrean/params"
)

var Global_config *Config

// ShyftTracer sets up params required to tracer transaction through debugAPI
type ShyftTracer struct {
	ChainConfig *params.ChainConfig
	TraceConfig *TraceConfig
	Eth         *Ethereum
}

// SetGlobalConfig gives access to eth.Config
func SetGlobalConfig(c *Config) {
	Global_config = c
}

// TraceTransaction invokes api debug.Transaction in order to trace a transaction for internal txs
func (st ShyftTracer) TraceTransaction(txhash common.Hash, blockhash common.Hash) (interface{}, error) {
	var ctx context.Context
	privateAPI := NewPrivateDebugAPI(st.ChainConfig, st.Eth)

	return privateAPI.STraceTransaction(ctx, txhash, blockhash, st.TraceConfig)
}
