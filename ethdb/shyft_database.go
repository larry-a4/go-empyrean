package ethdb

import (
	"database/sql"
	"flag"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/ShyftNetwork/go-empyrean/core/sTypes"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var blockExplorerDb *sqlx.DB

const (
	defaultTestDb  = "shyftdbtest"
	defaultDb      = "shyftdb"
	connStrTest    = "user=postgres password=docker sslmode=disable"
	connStrDocker  = "user=postgres host=pg password=docker sslmode=disable"
	connStrDefault = "user=postgres host=localhost sslmode=disable"
)

var TestDbInstances []string

type SPGDatabase struct {
	db *sqlx.DB // PostgresDB instance
}

// NewLDBDatabase returns a PostgresDB wrapped object.
func NewShyftDatabase() (*SPGDatabase, error) {
	if blockExplorerDb == nil {
		_, err := InitDB(false)
		if err != nil {
			return nil, err
		}
	}
	conn := blockExplorerDb
	conn.Ping()
	return &SPGDatabase{
		db: conn,
	}, nil
}

func NewTestInstanceShyftDatabase() (*SPGDatabase, error) {
	if blockExplorerDb == nil {
		_, err := InitDB(true)
		if err != nil {
			return nil, err
		}
	}
	conn := blockExplorerDb
	conn.Ping()
	return &SPGDatabase{
		db: conn,
	}, nil
}

func ReturnShyftDatabase() (*SPGDatabase, error) {
	return NewShyftDatabase()
}

// InitDB - initializes a Postgresql Database for use by the Blockexplorer
func InitDB(flag bool) (*sqlx.DB, error) {
	// To set the environment you can run the program with an ENV variable DBENV.
	// DBENV defaults to local for purposes of running the correct local
	// database connection parameters but will use docker connection parameters if DBENV=docker
	if flag {
		DbTestName := AssignTestDbInstanceName()
		DeletePgDb(DbTestName)
		exist, _ := DbExists(DbTestName)
		if !exist {
			CreatePgDb(DbTestName)
		}
	} else {
		// Check for existence of Database
		exist, _ := DbExists(DbName())
		if !exist {
			// create the db
			CreatePgDb(DbName())
		}
	}
	// connect to the designated db & create tables if necessary
	blockExplorerDb = Connect(ShyftConnectStr())
	blockExplorerDb.MustExec(MakeTableQuery())
	return blockExplorerDb, nil
}

func stripNumber(str string) int {
	re := regexp.MustCompile(`\w*_([0-9]+)$`)
	match := re.FindStringSubmatch(str)
	d, err := strconv.Atoi(match[1])
	if err != nil {
		return -1
	}
	return d
}

func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// AssignTestDbInstanceName - returns a Db Name for testing
func AssignTestDbInstanceName() string {
	// returns name of current test db
	dbNameAssigned := "_1"
	TestDbInstances = append(TestDbInstances, dbNameAssigned)
	var dbNumbersUsed []int
	for _, x := range TestDbInstances {
		dbNumbersUsed = append(dbNumbersUsed, stripNumber(x))
	}

	dbNum := false
	dbInt := 1
	for !dbNum {
		if intInSlice(dbInt, dbNumbersUsed) {
			dbInt++
		} else {
			dbNum = true
		}
	}
	dbNameAssigned = defaultTestDb + "_" + strconv.Itoa(dbInt)
	TestDbInstances = append(TestDbInstances, dbNameAssigned)
	return dbNameAssigned
}

func (db *SPGDatabase) AccountExists(addr string) (string, string, error) {
	var addressBalance, accountNonce string
	sqlExistsStatement := `SELECT balance, nonce from accounts WHERE addr = ($1)`
	err := db.db.QueryRow(sqlExistsStatement, strings.ToLower(addr)).Scan(&addressBalance, &accountNonce)
	switch {
	case err == sql.ErrNoRows:
		return addressBalance, accountNonce, err
	case err != nil:
		panic(err)
	default:
		return addressBalance, accountNonce, err
	}
}

//BlockExists checks if block exists in Postgres Db
//Refactor as a transaction
func (db *SPGDatabase) BlockExists(hash string) bool {
	var res bool
	sqlExistsStatement := `SELECT exists(select hash from blocks WHERE hash= ($1));`
	err := db.db.QueryRow(sqlExistsStatement, strings.ToLower(hash)).Scan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err)
		}
	}
	return res
}

//IsContract checks if toAddr is from a contract in Postgres Db
func (db *SPGDatabase) IsContract(addr string) bool {
	var isContract bool
	sqlExistsStatement := `SELECT isContract from txs WHERE to_addr=($1);`
	err := db.db.QueryRowx(sqlExistsStatement, strings.ToLower(addr)).Scan(&isContract)
	switch {
	case err == sql.ErrNoRows:
		return isContract
	default:
		return isContract
	}
}

// Transact - A wrapper around pg - transaction to allow a panic after a rollback
func (db *SPGDatabase) Transact(txFunc func(*sqlx.Tx) error) (err error) {
	tx, err := db.db.Beginx()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()
	err = txFunc(tx)
	return err
}

//CreateAccount writes new account to Postgres Db
func (db *SPGDatabase) CreateAccount(addr string, balance string, nonce string) error {
	addr = strings.ToLower(addr)
	return db.Transact(func(tx *sqlx.Tx) error {

		accountStmnt := FindOrCreateAcctStmnt

		if _, err := tx.Exec(accountStmnt, addr, balance, nonce); err != nil {
			return err
		}
		return nil
	})
}

//updateMinerAccount updates account in Postgres Db
func (db *SPGDatabase) UpdateMinerAccount(addr string, blockHash string, reward *big.Int) error {
	rewardInt := reward.Int64()
	addr = strings.ToLower(addr)

	return db.Transact(func(tx *sqlx.Tx) error {
		// Updates and/or Creates Account for Miner if it doesnt exist
		_, err := tx.Exec(UpdateToBalanceNonce, addr, rewardInt)
		if err != nil {
			panic(err)
		}
		if rewardInt != 0 {
			_, err = tx.Exec(FindOrCreateAcctBlockStmnt, addr, blockHash, rewardInt)
			if err != nil {
				panic(err)
			}
		}
		return nil
	})
}

//InsertBlock writes block to Postgres Db
func (db *SPGDatabase) InsertBlock(blockData stypes.SBlock) {
	sqlStatement := `INSERT INTO blocks(hash, coinbase, number, gasUsed, gasLimit, txCount, uncleCount, age, parentHash, uncleHash, difficulty, size, rewards, nonce) VALUES(($1), ($2), ($3), ($4), ($5), ($6), ($7), ($8), ($9), ($10), ($11), ($12),($13), ($14)) RETURNING number;`
	qerr := db.db.QueryRow(sqlStatement, strings.ToLower(blockData.Hash), strings.ToLower(blockData.Coinbase), blockData.Number, blockData.GasUsed, blockData.GasLimit, blockData.TxCount, blockData.UncleCount, blockData.Age, blockData.ParentHash, blockData.UncleHash, blockData.Difficulty, blockData.Size, blockData.Rewards, strconv.FormatUint(blockData.Nonce, 10)).Scan(&blockData.Number)
	if qerr != nil {
		panic(qerr)
	}
}

//InsertTx writes tx to Postgres Db
func (db *SPGDatabase) InsertTx(txData stypes.ShyftTxEntryPretty) error {
	acctAddrs := [2]string{strings.ToLower(txData.To), strings.ToLower(txData.From)}
	return db.Transact(func(tx *sqlx.Tx) error {
		txHash := strings.ToLower(txData.TxHash)
		// @SHYFT NOTE - AFTER CHAIN RESTART IT APPEARS CURRENTLY TRANSACTION JOURNAL OR REPEATED SEND TRANSACTION IS NOT FUNCTIONING
		txExistsStmnt := fmt.Sprintf(`select exists(SELECT txhash FROM txs WHERE txhash = '%s');`, txHash)
		var exists bool
		err := db.db.QueryRow(txExistsStmnt).Scan(&exists)
		if err != nil {
			panic(err)
		}
		if !exists {
			toAcctCredit := new(big.Int)
			toAcctCredit, _ = toAcctCredit.SetString(txData.Amount, 10)
			var one = big.NewInt(-1)
			fromAcctDebit := new(big.Int).Mul(toAcctCredit, one)
			// Add Transaction Table entry
			_, err = tx.Exec(CreateTxTableStmnt, txHash, strings.ToLower(txData.From),
				strings.ToLower(txData.To), strings.ToLower(txData.BlockHash), txData.BlockNumber, txData.Amount,
				txData.GasPrice, txData.Gas, txData.GasLimit, txData.Cost, txData.Nonce, txData.IsContract,
				txData.Status, txData.Age, txData.Data)
			if err != nil {
				fmt.Println("CREATE TX TABLE ISSUE")
				panic(err)
			}
			// Update account balances and account Nonces
			// Updates/Creates Account for To
				_, err = tx.Exec(UpdateToBalanceNonce, acctAddrs[0], toAcctCredit.String())
				if err != nil {
					fmt.Println("UPDATE BALANCE NONCE ISSUE")
					panic(err)
				}
				//Update/Create TO accountblock
				_, err = tx.Exec(FindOrCreateAcctBlockStmnt, acctAddrs[0], txData.BlockHash, toAcctCredit.String())
				if err != nil {
					panic(err)
				}
			if acctAddrs[1] != "genesis" {
				// Updates/Creates Account for From
				_, err = tx.Exec(UpdateBalanceNonce, acctAddrs[1], fromAcctDebit.String())
				if err != nil {
					panic(err)
				}
				//Update/Create FROM accountblock
				_, err = tx.Exec(FindOrCreateAcctBlockStmnt, acctAddrs[1], txData.BlockHash, fromAcctDebit.String())
				if err != nil {
					panic(err)
				}
			}
		}
		return nil
	})
}

//InsertInternals - Inserts transactions to pg internaltxs and updates/creates accounts/accountblocks tables
//accordingly
func (db *SPGDatabase) InsertInternals(i stypes.InteralWrite) error {
	acctAddrs := [2]string{strings.ToLower(i.To), strings.ToLower(i.From)}

	return db.Transact(func(tx *sqlx.Tx) error {

		toAcctCredit, _ := strconv.Atoi(i.Value)
		fromAcctDebit := -1 * toAcctCredit
		// Update account balances and account Nonces
		// Updates/Creates Account for To
		_, err := tx.Exec(UpdateToBalanceNonce, acctAddrs[0], toAcctCredit)
		if err != nil {
			panic(err)
		}
		// Updates/Creates Account for From
		_, err = tx.Exec(UpdateBalanceNonce, acctAddrs[1], fromAcctDebit)
		if err != nil {
			panic(err)
		}
		// // Add Internal Transaction Table entry
		_, err = tx.Exec(CreateInternalTxTableStmnt, i.Action, strings.ToLower(i.Hash), strings.ToLower(i.BlockHash), strings.ToLower(i.From), strings.ToLower(i.To), i.Value, i.Gas, i.GasUsed, i.Time, i.Input, i.Output)
		if err != nil {
			panic(err)
		}
		if i.Value != "0" {
			//Update/Create TO accountblock
			_, err = tx.Exec(FindOrCreateAcctBlockStmnt, acctAddrs[0], i.BlockHash, toAcctCredit)
			if err != nil {
				panic(err)
			}
			//Update/Create FROM accountblock
			_, err = tx.Exec(FindOrCreateAcctBlockStmnt, acctAddrs[1], i.BlockHash, fromAcctDebit)
			if err != nil {
				panic(err)
			}
		}
		return nil
	})
	return nil
}

//RollbackPgDb - rollsback the PG database by:
// deleting blocks designated by the passed in Blockheaders
// deleting all transactions contained in the foregoing Blockheaders
// reversing each account balance by the delta included in the account blocks table
// reversing the account nonce values by the transaction count included in the accountblocks table
func (db *SPGDatabase) RollbackPgDb(blockheaders []string) error {

	return db.Transact(func(tx *sqlx.Tx) error {
		acctBlockStmnt := `SELECT * FROM accountblocks WHERE accountblocks.blockhash = ANY($1)`
		accountBlocks := []AccountBlock{}

		// Get all accountblocks containing the blockhash
		err := tx.Select(&accountBlocks, acctBlockStmnt, pq.Array(blockheaders))
		if err != nil {
			panic(err)
		}
		// Rollback account balances
		for _, acctBlock := range accountBlocks {
			// Get delta and txCount from accountblocks and adjust account balance and account nonce accordingly
			_, err = tx.Exec(AccountRollback, acctBlock.Acct, int(acctBlock.Delta), int(acctBlock.TxCount))
			if err != nil {
				panic(err)
			}
		}
		// Prune all 0 Balance accounts
		delZeroBalances := `DELETE FROM accounts WHERE balance = '0'`
		_, err = tx.Exec(delZeroBalances)
		if err != nil {
			panic(err)
		}
		// Delete all transactions containing the blockhash
		_, err = tx.Exec(TransactionRollback, pq.Array(blockheaders))
		if err != nil {
			panic(err)
		}
		// Delete all internal transactions containing the blockhash
		_, err = tx.Exec(InternalTransactionRollback, pq.Array(blockheaders))
		if err != nil {
			panic(err)
		}

		// Delete all blocks whose hash is within the blockheader array
		_, err = tx.Exec(BlockRollback, pq.Array(blockheaders))
		if err != nil {
			panic(err)
		}

		// Delete all accountblocks whose blockhash is included in the blockheader array
		_, err = tx.Exec("DELETE from accountblocks WHERE blockhash = ANY($1)", pq.Array(blockheaders))
		if err != nil {
			panic(err)
		}
		return nil
	})
	return nil
}

// ShyftConnectStr - Returns the Connection String With The appropriate database
func ShyftConnectStr() string {
	return fmt.Sprintf("%s%s%s", ConnectionStr(), " dbname=", DbName())
}

// Connect - return a connection to a postgres database wi
func Connect(connectURL string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", connectURL)
	if err != nil {
		fmt.Println("ERROR OPENING DB, NOT INITIALIZING")
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

// DbName - gets the correct db name based on the environment
func DbName() string {
	if flag.Lookup("test.v") == nil {
		return defaultDb
	} else {
		//dbTestName := AssignTestDbInstanceName()
		return defaultTestDb
	}
}

// CreatePgDb - Creates a DB
func CreatePgDb(dbName string) {
	conn := Connect(ConnectionStr())
	sqlCmd := fmt.Sprintf(`CREATE DATABASE %s;`, dbName)
	_, err := conn.Exec(sqlCmd)
	if err != nil {
		panic(err)
	}
	conn.Close()
}

// DeletePgDb - Deletes the designated DB
func DeletePgDb(dbName string) {
	conn := Connect(ConnectionStr())
	q1 := fmt.Sprintf(`SELECT pg_terminate_backend(pid)FROM pg_stat_activity WHERE datname = '%s';`, dbName)
	_, err1 := conn.Exec(q1)
	if err1 != nil || err1 == sql.ErrNoRows {
		panic(err1)
	}
	q2 := fmt.Sprintf(`DROP DATABASE IF EXISTS %s;`, dbName)
	_, err2 := conn.Exec(q2)
	if err2 != nil || err2 == sql.ErrNoRows {
		panic(err2)
	}
	conn.Close()
}

// ConnectionStr - return a Connection to the PG admin database
func ConnectionStr() string {
	dbEnv := os.Getenv("DBENV")
	if flag.Lookup("test.v") == nil {
		switch dbEnv {
		default:
			return connStrDefault
		case "docker":
			return connStrDocker
		}
	} else {
		return connStrTest
	}
}

// DbExists - Checks whether the named database exists returns true or false
func DbExists(dbname string) (bool, error) {
	sqldb := Connect(ConnectionStr())
	var exists bool
	sqlStatement := fmt.Sprintf(`select exists(SELECT datname FROM pg_catalog.pg_database WHERE lower(datname) = '%s');`, strings.ToLower(dbname))
	error := sqldb.QueryRow(sqlStatement).Scan(&exists)
	switch {
	case error == sql.ErrNoRows:
		sqldb.Close()
		fmt.Println("DB Exists ", error)
		return false, error
	case error != nil:
		return false, error
		panic(error)
	default:
		sqldb.Close()
		return exists, error
	}
}

// TruncateTables - Is primarily user to clear the pg database between unit tests
func (db *SPGDatabase) TruncateTables() {
	tx, _ := db.db.Begin()
	sqlStatement := `TRUNCATE TABLE txs, accounts, blocks, internaltxs RESTART IDENTITY CASCADE;`
	_, err := tx.Exec(sqlStatement)
	tx.Commit()
	if err != nil {
		panic(err)
	}
}
