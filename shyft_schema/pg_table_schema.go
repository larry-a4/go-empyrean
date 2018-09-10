package shyftschema

import (
	"fmt"
	"time"
)

// This Package Contains The Queries To Make Tables Required For the PG Database & Structs
// For Scanning and Executing queries to the db database - Many of the structs are used in tests
// So if the schema is changed the structs should be updated accordingly

// MakeTableQuery - returns sql to create db tables
func MakeTableQuery() string {
	return fmt.Sprintf(`%s %s %s %s %s`, blocksTable, txsTable, accountsTable, accountBlocksTable, internalTxsTable)
}

// AccountsTable sql for accounts
// @Shyft NOTE - Benchmark Index for blockDelta

// Account - - struct for reading and writing database data
type Account struct {
	Addr    string `db:"addr"`
	Balance int64  `db:"balance"`
	Nonce   int64  `db:"nonce"`
}

const accountsTable = `
CREATE TABLE IF NOT EXISTS accounts (
  addr text primary key unique,
  balance bigint,
  nonce bigint
);
`

// Block - - struct for reading and writing database data
type Block struct {
	Hash       string    `db:"hash"`
	Coinbase   string    `db:"coinbase"`
	GasUsed    int64     `db:"gasUsed"`
	GasLimit   int64     `db:"gasLimit"`
	TxCount    int64     `db:"txCount"`
	UncleCount int64     `db:"uncleCount"`
	Age        time.Time `db:"age"`
	ParentHash string    `db:"parentHash"`
	UncleHash  string    `db:"uncleHash"`
	Difficulty int64     `db:"difficulty"`
	Size       string    `db:"size"`
	Nonce      int64     `db:"nonce"`
	Rewards    int64     `db:"rewards"`
	Number     int64     `db:"number"`
}

// BlocksTable sql for Blocks
const blocksTable = `
CREATE TABLE IF NOT EXISTS blocks (
  hash text primary key,
  coinbase text,
  gasUsed numeric,
  gasLimit numeric,
  txCount numeric,
  uncleCount numeric,
  age timestamp,
  parentHash text,
  uncleHash text,
  difficulty bigint,
  size text,
  nonce numeric,
  rewards numeric,
  number bigint
);`

// AccountBlock - struct for reading and writing database data
type AccountBlock struct {
	ID        uint64 `db:"id"`
	Acct      string `db:"acct"`
	Blockhash string `db:"blockhash"`
	Delta     int64  `db:"delta"`
}

const accountBlocksTable = `
CREATE TABLE IF NOT EXISTS accountblocks ( 
  id SERIAL PRIMARY KEY,
  acct text NOT NULL REFERENCES accounts(addr) ON DELETE CASCADE ON UPDATE CASCADE DEFERRABLE, 
  blockhash text NOT NULL,
  delta numeric
);
CREATE INDEX IF NOT EXISTS idx_acct_ab ON accountblocks (acct);
CREATE INDEX IF NOT EXISTS idx_block_ab ON accountblocks (blockhash);
`

// TxsTable sql for transactions
const txsTable = `
CREATE TABLE IF NOT EXISTS txs (
  txHash text primary key unique,
  to_addr text,
  from_addr text,
  blockhash text references blocks(hash) ON DELETE CASCADE,
  blocknumber text,
  amount numeric,
  gasprice numeric,
  gas numeric,
  gasLimit numeric,
  txFee numeric,
  nonce numeric,
  txStatus text,
  isContract bool,
  age timestamp,
  data bytea
);
CREATE INDEX IF NOT EXISTS idx_block_txs ON txs (blockhash);
`

// InternalTxsTable sql for transactions
const internalTxsTable = `
CREATE TABLE IF NOT EXISTS internalTxs (
  id SERIAL PRIMARY KEY,
  txHash text references txs(txHash) ON DELETE CASCADE,
  action text,
  to_addr text,
  from_addr text,
  amount text,
  gas numeric,
  gasUsed numeric,
  time text,
  input text,
  output text
);
CREATE INDEX IF NOT EXISTS idx_tx_txs ON txs (txHash);
`
