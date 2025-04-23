package cmd

import (
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// fallbackTxConfig is a wrapper around the client.TxConfig interface that
// provides a custom TxDecoder function.
type fallbackTxConfig struct {
	client.TxConfig
}

// NewClientTxConfig creates a new instance of fallbackTxConfig
func newFallbackTxConfig(txConfig client.TxConfig) client.TxConfig {
	return fallbackTxConfig{txConfig}
}

// TxDecoder returns a custom TxDecoder function that handles decoding failed
// transactions. This is needed to handle connect price oracle's pseudo txs.
func (c fallbackTxConfig) TxDecoder() sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {
		if tx, err := c.TxConfig.TxDecoder()(txBytes); err != nil {
			txBuilder := c.NewTxBuilder()
			txBuilder.SetMemo("decode failed tx")

			return txBuilder.GetTx(), nil
		} else {
			return tx, err
		}
	}
}
