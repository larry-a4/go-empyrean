package core

import (
	"database/sql"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ShyftNetwork/go-empyrean/shyft_schema"

	"fmt"

	"github.com/ShyftNetwork/go-empyrean/common"
	Rewards "github.com/ShyftNetwork/go-empyrean/consensus/ethash"
	"github.com/ShyftNetwork/go-empyrean/core/sTypes"
	"github.com/ShyftNetwork/go-empyrean/core/types"
	"github.com/ShyftNetwork/go-empyrean/shyfttracerinterface"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

//IShyftTracer Used to initialize ShyftTracer
var IShyftTracer shyfttracerinterface.IShyftTracer

//SetIShyftTracer sets tracer type
func SetIShyftTracer(st shyfttracerinterface.IShyftTracer) {
	IShyftTracer = st
}

//SWriteBlock writes to block info to sql db
func SWriteBlock(block *types.Block, receipts []*types.Receipt) error {
	//Get miner rewards
	rewards := swriteMinerRewards(block)
	//Format block time to be stored
	i, err := strconv.ParseInt(block.Time().String(), 10, 64)
	if err != nil {
		panic(err)
	}
	age := time.Unix(i, 0)

	blockData := stypes.SBlock{
		Hash:       block.Header().Hash().Hex(),
		Coinbase:   block.Header().Coinbase.String(),
		Number:     block.Header().Number.String(),
		GasUsed:    block.Header().GasUsed,
		GasLimit:   block.Header().GasLimit,
		TxCount:    block.Transactions().Len(),
		UncleCount: len(block.Uncles()),
		ParentHash: block.ParentHash().String(),
		UncleHash:  block.UncleHash().String(),
		Difficulty: block.Difficulty().String(),
		Size:       block.Size().String(),
		Nonce:      block.Nonce(),
		Rewards:    rewards,
		Age:        age,
	}

	//Inserts block data into DB
	InsertBlock(blockData)

	if block.Transactions().Len() > 0 {
		for _, tx := range block.Transactions() {
			swriteTransactions(tx, block.Header().Hash(), blockData.Number, receipts, age, blockData.GasLimit)
		}
	}
	return nil
}

//swriteTransactions writes to sqldb, a SHYFT postgres instance
func swriteTransactions(tx *types.Transaction, blockHash common.Hash, blockNumber string, receipts []*types.Receipt, age time.Time, gasLimit uint64) error {
	var isContract bool
	var statusFromReciept, toAddr string
	var contractAddressFromReciept common.Address
	if tx.To() == nil {
		for _, receipt := range receipts {
			statusReciept := (*types.ReceiptForStorage)(receipt).Status
			contractAddressFromReciept = (*types.ReceiptForStorage)(receipt).ContractAddress
			switch {
			case statusReciept == 0:
				statusFromReciept = "FAIL"
			case statusReciept == 1:
				statusFromReciept = "SUCCESS"
			}
		}
		isContract = true
		tempAddr := &contractAddressFromReciept
		toAddr = tempAddr.String()
	} else {
		isContract = false
		for _, receipt := range receipts {
			statusReciept := (*types.ReceiptForStorage)(receipt).Status
			switch {
			case statusReciept == 0:
				statusFromReciept = "FAIL"
			case statusReciept == 1:
				statusFromReciept = "SUCCESS"
			}
		}
		toAddr = tx.To().String()
	}

	txData := stypes.ShyftTxEntryPretty{
		TxHash:      tx.Hash().Hex(),
		From:        tx.From().Hex(),
		To:          toAddr,
		BlockHash:   blockHash.Hex(),
		BlockNumber: blockNumber,
		Amount:      tx.Value().String(),
		Cost:        tx.Cost().String(),
		GasPrice:    tx.GasPrice().Uint64(),
		GasLimit:    gasLimit,
		Gas:         tx.Gas(),
		Nonce:       tx.Nonce(),
		Age:         age,
		Data:        tx.Data(),
		Status:      statusFromReciept,
		IsContract:  isContract,
	}
	isContractCheck := IsContract(txData.To)
	if isContractCheck == true {
		InsertTx(txData)
		//Runs necessary functions for tracing internal transactions through tracers.go
		IShyftTracer.GetTracerToRun(tx.Hash(), blockHash)
	} else {
		//Inserts Tx into DB
		InsertTx(txData)
	}
	return nil
}

// @NOTE: This function is extremely complex and requires heavy testing and knowdlege of edge cases:
// uncle blocks, account balance updates based on reorgs, diverges that get dropped.
// Reason for this is because the accounts are not deterministic like the block and tx hashes.
// @TODO: Calculate reorg - Determine how we include in acctblocks
func swriteMinerRewards(block *types.Block) string {
	minerAddr := block.Coinbase().String()
	shyftConduitAddress := Rewards.ShyftNetworkConduitAddress.String()
	// Calculate the total gas used in the block
	totalGas := new(big.Int)
	for _, tx := range block.Transactions() {
		totalGas.Add(totalGas, new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(tx.Gas())))
	}

	blockHash := block.Hash().Hex()
	totalMinerReward := totalGas.Add(totalGas, Rewards.ShyftMinerBlockReward)

	// References:
	// https://ethereum.stackexchange.com/questions/27172/different-uncles-reward
	// line 551 in consensus.go (go-empyrean/consensus/ethash/consensus.go)
	// Some weird constants to avoid constant memory allocs for them.
	var big8 = big.NewInt(8)
	var uncleRewards []*big.Int
	var uncleAddrs []string

	// uncleReward is overwritten after each iteration
	// Based on calculation in consensus.go accumulateRewards()
	uncleReward := new(big.Int)
	for _, uncle := range block.Uncles() {
		uncleReward.Add(uncle.Number, big8)
		uncleReward.Sub(uncleReward, block.Number())
		uncleReward.Mul(uncleReward, Rewards.ShyftMinerBlockReward)
		uncleReward.Div(uncleReward, big8)
		uncleRewards = append(uncleRewards, uncleReward)
		uncleAddrs = append(uncleAddrs, uncle.Coinbase.String())
	}
	updateMinerAccount(minerAddr, blockHash, totalMinerReward)
	updateMinerAccount(shyftConduitAddress, blockHash, Rewards.ShyftNetworkBlockReward)
	var uncRewards = new(big.Int)
	for i := 0; i < len(uncleAddrs); i++ {
		_ = uncleRewards[i]
		updateMinerAccount(uncleAddrs[i], blockHash, uncleRewards[i])
	}

	fullRewardValue := new(big.Int)
	fullRewardValue.Add(totalMinerReward, Rewards.ShyftNetworkBlockReward)
	fullRewardValue.Add(fullRewardValue, uncRewards)

	return fullRewardValue.String()
}

///////////////////////
//DB Utility functions
//////////////////////

//AccountExists checks if account exists in Postgres Db
// Refactor - Transaction
func AccountExists(addr string) (string, string, error) {
	sqldb, _ := DBConnection()
	var addressBalance, accountNonce string
	sqlExistsStatement := `SELECT balance, nonce from accounts WHERE addr = ($1)`

	err := sqldb.QueryRow(sqlExistsStatement, strings.ToLower(addr)).Scan(&addressBalance, &accountNonce)
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
func BlockExists(hash string) bool {
	var res bool
	sqlExistsStatement := `SELECT exists(select hash from blocks WHERE hash= ($1));`
	sqldb, _ := DBConnection()
	err := sqldb.QueryRow(sqlExistsStatement, strings.ToLower(hash)).Scan(&res)
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
func IsContract(addr string) bool {
	sqldb, _ := DBConnection()
	var isContract bool
	sqlExistsStatement := `SELECT isContract from txs WHERE to_addr=($1);`
	err := sqldb.QueryRowx(sqlExistsStatement, strings.ToLower(addr)).Scan(&isContract)
	switch {
	case err == sql.ErrNoRows:
		return isContract
	default:
		return isContract
	}
}

// Transact - A wrapper around pg - transaction to allow a panic after a rollback
func Transact(db *sqlx.DB, txFunc func(*sqlx.Tx) error) (err error) {
	tx, err := db.Beginx()
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
func CreateAccount(addr string, balance string, nonce string) error {
	sqldb, _ := DBConnection()
	addr = strings.ToLower(addr)
	//bal := new(big.Int)
	//numericBalance, _ := bal.SetString(balance, 10)

	//non := new(big.Int)
	//intNonce, _ := non.SetString(nonce, 10)
	return Transact(sqldb, func(tx *sqlx.Tx) error {

		accountStmnt := shyftschema.FindOrCreateAcctStmnt

		if _, err := tx.Exec(accountStmnt, addr, balance, nonce); err != nil {
			return err
		}
		return nil
	})
}

//updateMinerAccount updates account in Postgres Db
func updateMinerAccount(addr string, blockHash string, reward *big.Int) error {
	sqldb, _ := DBConnection()
	rewardInt := reward.Int64()
	addr = strings.ToLower(addr)

	return Transact(sqldb, func(tx *sqlx.Tx) error {
		// Updates and/or Creates Account for Miner if it doesnt exist
		_, err := tx.Exec(shyftschema.UpdateBalanceNonce, addr, rewardInt)
		if err != nil {
			panic(err)
		}
		if rewardInt != 0 {
			_, err = tx.Exec(shyftschema.FindOrCreateAcctBlockStmnt, addr, blockHash, rewardInt)
			if err != nil {
				panic(err)
			}
		}
		return nil
	})
}

//InsertBlock writes block to Postgres Db
func InsertBlock(blockData stypes.SBlock) {
	sqldb, _ := DBConnection()
	sqlStatement := `INSERT INTO blocks(hash, coinbase, number, gasUsed, gasLimit, txCount, uncleCount, age, parentHash, uncleHash, difficulty, size, rewards, nonce) VALUES(($1), ($2), ($3), ($4), ($5), ($6), ($7), ($8), ($9), ($10), ($11), ($12),($13), ($14)) RETURNING number;`
	qerr := sqldb.QueryRow(sqlStatement, strings.ToLower(blockData.Hash), strings.ToLower(blockData.Coinbase), blockData.Number, blockData.GasUsed, blockData.GasLimit, blockData.TxCount, blockData.UncleCount, blockData.Age, blockData.ParentHash, blockData.UncleHash, blockData.Difficulty, blockData.Size, blockData.Rewards, strconv.FormatUint(blockData.Nonce, 10)).Scan(&blockData.Number)
	if qerr != nil {
		fmt.Println("INSERT BLOCK ISSUE Executing This Command")
		fmt.Println(fmt.Sprintf(sqlStatement, strings.ToLower(blockData.Hash), strings.ToLower(blockData.Coinbase), blockData.Number, blockData.GasUsed, blockData.GasLimit, blockData.TxCount, blockData.UncleCount, blockData.Age, blockData.ParentHash, blockData.UncleHash, blockData.Difficulty, blockData.Size, blockData.Rewards, blockData.Nonce))
		panic(qerr)
	}
}

//InsertTx writes tx to Postgres Db
func InsertTx(txData stypes.ShyftTxEntryPretty) error {
	acctAddrs := [2]string{strings.ToLower(txData.To), strings.ToLower(txData.From)}
	sqldb, _ := DBConnection()

	return Transact(sqldb, func(tx *sqlx.Tx) error {
		txHash := strings.ToLower(txData.TxHash)
		// @SHYFT NOTE - AFTER CHAIN RESTART IT APPEARS CURRENTLY TRANSACTION JOURNAL OR REPEATED SEND TRANSACTION IS NOT FUNCTIONING
		// NEED TO CHECK sendTransactions JS TO CONFIRM WHAT IS HAPPENING ALSO BLOCK HEIGHT IS NOT SET CORRECTLY FOR BLOCKS AFTER CHAIN
		// RESTART
		// RELATED BLOCKS
		txExistsStmnt := fmt.Sprintf(`select exists(SELECT txhash FROM txs WHERE txhash = '%s');`, txHash)
		var exists bool
		err := sqldb.QueryRow(txExistsStmnt).Scan(&exists)
		if err != nil {
			panic(err)
		}
		if !exists {
			toAcctCredit := new(big.Int)
			toAcctCredit, _ = toAcctCredit.SetString(txData.Amount, 10)
			var one = big.NewInt(-1)
			fromAcctDebit := new(big.Int).Mul(toAcctCredit, one)
			// Add Transaction Table entry
			_, err = tx.Exec(shyftschema.CreateTxTableStmnt, txHash, strings.ToLower(txData.From),
				strings.ToLower(txData.To), strings.ToLower(txData.BlockHash), txData.BlockNumber, txData.Amount,
				txData.GasPrice, txData.Gas, txData.GasLimit, txData.Cost, txData.Nonce, txData.IsContract,
				txData.Status, txData.Age, txData.Data)
			if err != nil {
				fmt.Println("CREATE TX TABLE ISSUE")
				panic(err)
			}
			// Update account balances and account Nonces
			// Updates/Creates Account for To
			_, err = tx.Exec(shyftschema.UpdateBalanceNonce, acctAddrs[0], toAcctCredit.String())
			if err != nil {
				fmt.Println("UPDATE BALANCE NONCE ISSUE")
				panic(err)
			}
			//Update/Create TO accountblock
			_, err = tx.Exec(shyftschema.FindOrCreateAcctBlockStmnt, acctAddrs[0], txData.BlockHash, toAcctCredit.String())
			if err != nil {
				panic(err)
			}
			if acctAddrs[1] != "genesis" {
				// Updates/Creates Account for From
				_, err = tx.Exec(shyftschema.UpdateBalanceNonce, acctAddrs[1], fromAcctDebit.String())
				if err != nil {
					panic(err)
				}
				//Update/Create FROM accountblock
				_, err = tx.Exec(shyftschema.FindOrCreateAcctBlockStmnt, acctAddrs[1], txData.BlockHash, fromAcctDebit.String())
				if err != nil {
					panic(err)
				}
			}
		}
		//
		return nil
	})
	return nil
}

//InsertInternals - Inserts transactions to pg internaltxs and updates/creates accounts/accountblocks tables
//accordingly
func InsertInternals(i stypes.InteralWrite) error {
	acctAddrs := [2]string{strings.ToLower(i.To), strings.ToLower(i.From)}
	sqldb, _ := DBConnection()

	return Transact(sqldb, func(tx *sqlx.Tx) error {

		toAcctCredit, _ := strconv.Atoi(i.Value)
		fromAcctDebit := -1 * toAcctCredit
		// Update account balances and account Nonces
		// Updates/Creates Account for To
		_, err := tx.Exec(shyftschema.UpdateBalanceNonce, acctAddrs[0], toAcctCredit)
		if err != nil {
			panic(err)
		}
		// Updates/Creates Account for From
		_, err = tx.Exec(shyftschema.UpdateBalanceNonce, acctAddrs[1], fromAcctDebit)
		if err != nil {
			panic(err)
		}
		// // Add Internal Transaction Table entry
		_, err = tx.Exec(shyftschema.CreateInternalTxTableStmnt, i.Action, strings.ToLower(i.Hash), strings.ToLower(i.BlockHash), strings.ToLower(i.From), strings.ToLower(i.To), i.Value, i.Gas, i.GasUsed, i.Time, i.Input, i.Output)
		if err != nil {
			panic(err)
		}
		if i.Value != "0" {
			//Update/Create TO accountblock
			_, err = tx.Exec(shyftschema.FindOrCreateAcctBlockStmnt, acctAddrs[0], i.BlockHash, toAcctCredit)
			if err != nil {
				panic(err)
			}
			//Update/Create FROM accountblock
			_, err = tx.Exec(shyftschema.FindOrCreateAcctBlockStmnt, acctAddrs[1], i.BlockHash, fromAcctDebit)
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
func RollbackPgDb(blockheaders []string) error {
	sqldb, _ := DBConnection()

	return Transact(sqldb, func(tx *sqlx.Tx) error {
		acctBlockStmnt := `SELECT * FROM accountblocks WHERE accountblocks.blockhash = ANY($1)`
		accountBlocks := []shyftschema.AccountBlock{}

		// Get all accountblocks containing the blockhash
		err := tx.Select(&accountBlocks, acctBlockStmnt, pq.Array(blockheaders))
		if err != nil {
			panic(err)
		}
		// Rollback account balances
		for _, acctBlock := range accountBlocks {
			// Get delta and txCount from accountblocks and adjust account balance and account nonce accordingly
			_, err = tx.Exec(shyftschema.AccountRollback, acctBlock.Acct, int(acctBlock.Delta), int(acctBlock.TxCount))
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
		_, err = tx.Exec(shyftschema.TransactionRollback, pq.Array(blockheaders))
		if err != nil {
			panic(err)
		}
		// Delete all internal transactions containing the blockhash
		_, err = tx.Exec(shyftschema.InternalTransactionRollback, pq.Array(blockheaders))
		if err != nil {
			panic(err)
		}

		// Delete all blocks whose hash is within the blockheader array
		_, err = tx.Exec(shyftschema.BlockRollback, pq.Array(blockheaders))
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
