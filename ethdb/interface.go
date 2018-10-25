// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package ethdb

import (
	"math/big"
	"github.com/ShyftNetwork/go-empyrean/core/sTypes"
)

// Code using batches should try to add this much data to the batch.
// The value was determined empirically.
const IdealBatchSize = 100 * 1024

// Putter wraps the database write operation supported by both batches and regular databases.
type Putter interface {
	Put(key []byte, value []byte) error
}

// Database wraps all database operations. All methods are safe for concurrent use.
type Database interface {
	Putter
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	Delete(key []byte) error
	Close()
	NewBatch() Batch
}

// Batch is a write-only database that commits changes to its host database
// when Write is called. Batch cannot be used concurrently.
type Batch interface {
	Putter
	ValueSize() int // amount of data in the batch
	Write() error
	// Reset resets the batch for reuse
	Reset()
}

type SDatabase interface {
	AccountExists(addr string) (string, string, error)
	BlockExists(hash string) bool
	IsContract(addr string) bool
	CreateAccount(addr string, balance string, nonce string) error
	UpdateMinerAccount(addr string, blockHash string, reward *big.Int) error
	InsertBlock(blockData stypes.SBlock)
	InsertTx(txData stypes.ShyftTxEntryPretty) error
	InsertInternals(i stypes.InteralWrite) error
	RollbackPgDb(blockheaders []string) error
	TruncateTables()
}

type SGetters interface {
	SGetAllBlocks() (string, error)
	SGetBlock(blockNumber string) string
	SGetRecentBlock() string
	SGetAllTransactionsFromBlock(blockNumber string) string
	SGetAllBlocksMinedByAddress(coinbase string) string
	SGetAllTransactions() string
	SGetTransaction(txHash string) string
	SGetAllAccounts() string
	SGetAccountTxs(address string) string
	SGetAllInternalTransactions() string
	SGetInternalTransaction(txHash string) string
}
