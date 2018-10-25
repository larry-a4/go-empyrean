package ethdb

import (
	"fmt"
	"time"
)

// This Package Contains The Queries To Make Tables Required For the PG Database & Structs
// For Scanning and Executing queries to the db database - Many of the structs are used in tests
// So if the schema is changed the structs should be updated accordingly

//FOR TABLE COLUMN NAMES PLEASE DO NOT USE CAMELCASE

// MakeTableQuery - returns sql to create db tables
func MakeTableQuery() string {
	return fmt.Sprintf(`%s %s %s %s %s`, blocksTable, txsTable, accountsTable, accountBlocksTable, internalTxsTable)
}

// AccountsTable sql for accounts
// @Shyft NOTE - Benchmark Index for blockDelta

// Account - - struct for reading and writing database data
type Account struct {
	Addr    string `db:"addr"`
	Balance string `db:"balance"`
	Nonce   uint64 `db:"nonce"`
}

const accountsTable = `
CREATE TABLE IF NOT EXISTS accounts (
  addr text primary key unique,
  balance numeric(78,0),
  nonce numeric
);
`

// Block - - struct for reading and writing database data
type Block struct {
	Hash       string    `db:"hash"`
	Coinbase   string    `db:"coinbase"`
	GasUsed    int64     `db:"gasused"`
	GasLimit   int64     `db:"gaslimit"`
	TxCount    int64     `db:"txcount"`
	UncleCount int64     `db:"unclecount"`
	Age        time.Time `db:"age"`
	ParentHash string    `db:"parenthash"`
	UncleHash  string    `db:"unclehash"`
	Difficulty string     `db:"difficulty"`
	Size       string    `db:"size"`
	Nonce      int64     `db:"nonce"`
	Rewards    string     `db:"rewards"`
	Number     string     `db:"number"`
}

// BlocksTable sql for Blocks
const blocksTable = `
CREATE TABLE IF NOT EXISTS blocks (
  hash text primary key,
  coinbase text,
  gasused numeric(78,0),
  gaslimit numeric(78,0),
  txcount numeric,
  unclecount numeric,
  age timestamp,
  parenthash text,
  unclehash text,
  difficulty bigint,
  size text,
  nonce numeric(78,0),
  rewards numeric(78,0),
  number numeric(78,0)
);`

// AccountBlock - struct for reading and writing database data
type AccountBlock struct {
	Acct      string `db:"acct"`
	Blockhash string `db:"blockhash"`
	Delta     int64  `db:"delta"`
	TxCount   int64  `db:"txcount"`
}

type AccountBlockArray struct {
	AccountBlocks []AccountBlock
}

const accountBlocksTable = `
CREATE TABLE IF NOT EXISTS accountblocks ( 
  acct text NOT NULL REFERENCES accounts(addr) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED, 
  blockhash text NOT NULL,
  delta numeric,
  txcount bigint,
  primary key(acct, blockhash)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_acct_bhash ON accountblocks (acct, blockhash);
CREATE INDEX IF NOT EXISTS idx_acct_ab ON accountblocks (acct);
CREATE INDEX IF NOT EXISTS idx_block_ab ON accountblocks (blockhash);
`

// PgTransaction struct for scanning db transactions from table txs
type PgTransaction struct {
	TxHash      string    `db:"txhash"`
	To          string    `db:"to_addr"`
	From        string    `db:"from_addr"`
	Blockhash   string    `db:"blockhash"`
	Blocknumber string    `db:"blocknumber"`
	Amount      string    `db:"amount"`
	Gasprice    uint64    `db:"gasprice"`
	Gas         uint64    `db:"gas"`
	GasLimit    uint64    `db:"gaslimit"`
	TxFee       string    `db:"txfee"`
	Nonce       uint64    `db:"nonce"`
	TxStatus    string    `db:"txstatus"`
	IsContract  bool      `db:"iscontract"`
	Age         time.Time `db:"age"`
	Data        []byte    `db:"data"`
}

// TxsTable sql for transactions
const txsTable = `
CREATE TABLE IF NOT EXISTS txs (
  txhash text primary key unique,
  to_addr text,
  from_addr text,
  blockhash text,
  blocknumber text,
  amount text,
  gasprice numeric,
  gas numeric,
  gaslimit numeric,
  txfee numeric,
  nonce numeric,
  txstatus text,
  iscontract bool,
  age timestamp,
  data bytea
);
CREATE INDEX IF NOT EXISTS idx_block_txs ON txs (blockhash);
`

// InternalTxsTable sql for transactions
const internalTxsTable = `
CREATE TABLE IF NOT EXISTS internalTxs (
  id SERIAL PRIMARY KEY,
  txhash text references txs(txhash) ON DELETE CASCADE,
  blockhash text references blocks(hash) ON DELETE CASCADE,
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
CREATE INDEX IF NOT EXISTS idx_tx_txs ON txs (txhash);
`

// AccountRollBack - finds the relevant account and reverses the balance and nonce
const AccountRollback = `
UPDATE accounts
  SET balance = ((SELECT balance FROM accounts WHERE addr = $1) - $2),
      nonce = ((SELECT nonce FROM accounts WHERE addr = $1) - $3)
WHERE addr = $1;
`

// TransactionRollback - deletes all transactions whose blockhash is contained in the array of blockheaders
const TransactionRollback = `
DELETE FROM txs WHERE blockhash = ANY($1);
`

// InternalTransactionRollback - deletes all transactions whose blockhash is contained in the array of blockheaders
const InternalTransactionRollback = `
DELETE FROM internaltxs WHERE blockhash = ANY($1);
`

// BlockRollback - deletes all blocks whose hash is contained in the array of blockheaders
const BlockRollback = `
DELETE FROM blocks WHERE hash = ANY($1);
`

// FindOrCreateAcctStmnt - query to create account if it doesnt exist - and return it if it does
// Parameters are addr = $1 balance = $2 nonce = $3
const FindOrCreateAcctStmnt = `
INSERT INTO accounts(addr, balance, nonce) VALUES($1, $2, $3)
ON CONFLICT ON CONSTRAINT accounts_pkey DO NOTHING;
`

//FindOrCreateAcctBlockStmnt - query to find or create an accountblock record returning
const FindOrCreateAcctBlockStmnt = `
INSERT INTO accountblocks(acct, blockhash, delta, txcount) VALUES($1, $2, $3, 1)
ON CONFLICT (acct, blockhash)
DO
  UPDATE
    SET delta = ((SELECT delta FROM accountblocks WHERE accountblocks.acct = $1 AND accountblocks.blockhash = $2) + $3),
        txcount = ((SELECT txcount FROM accountblocks WHERE accountblocks.acct = $1 AND accountblocks.blockhash = $2) + 1)
WHERE accountblocks.acct = $1 AND accountblocks.blockhash = $2;  
`

//UpdateBalanceNonce - query to update the balance and nonce for a transaction
// Parameters are addr = $1 amount = $2
const UpdateBalanceNonce = `
INSERT INTO accounts (addr, balance, nonce) VALUES($1, $2, 1) 
ON CONFLICT ON CONSTRAINT accounts_pkey DO 
UPDATE 
	SET balance = ((SELECT balance from accounts where accounts.addr = $1) + $2), 
		nonce = ((SELECT nonce from accounts where accounts.addr = $1) + 1) 
WHERE accounts.addr = $1
`

//Creates an internal transaction records in the internal tx table
const CreateInternalTxTableStmnt = `
INSERT INTO internaltxs(action, txhash, blockHash, from_addr, to_addr, amount, gas, gasUsed, time, input, output)
VALUES(($1), ($2), ($3), ($4), ($5), ($6), ($7), ($8), ($9), ($10), ($11));
`

//Creates a transaction record in the tx table
const CreateTxTableStmnt = `
INSERT INTO txs(txhash, from_addr, to_addr, blockhash, blockNumber, amount, gasprice, gas, gasLimit, txfee, nonce, isContract, txStatus, age, data)
VALUES(($1), ($2), ($3), ($4), ($5), ($6), ($7), ($8), ($9), ($10), ($11), ($12), ($13), ($14), ($15)) ON CONFLICT ON CONSTRAINT txs_pkey DO NOTHING;
`

//FindOrCreateAcctBlockStmntForInternals - query to find or create an accountblock record returning
const FindOrCreateAcctBlockStmntForInternals = `
INSERT INTO accountblocks(acct, blockhash, delta, txcount) VALUES($1, (()), $3, 1)
ON CONFLICT ON CONSTRAINT accountblocks_pkey
DO
  UPDATE
    SET delta = ((SELECT delta FROM accountblocks WHERE accountblocks.acct = $1 AND accountblocks.blockhash = $2) + $3),
        txcount = ((SELECT txcount FROM accountblocks WHERE accountblocks.acct = $1 AND accountblocks.blockhash = $2) + 1)
WHERE accountblocks.acct = $1 AND accountblocks.blockhash = $2;  
`
