package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	stypes "github.com/ShyftNetwork/go-empyrean/core/sTypes"
	"github.com/jmoiron/sqlx"
)

///////////
// Getters
//////////
func SGetAllBlocks(sqldb *sqlx.DB) string {
	var arr stypes.BlockRes
	var blockArr string

	rows, err := sqldb.Queryx(`SELECT * FROM blocks ORDER BY number ASC`)
	if err != nil {
		fmt.Println("err")
	}
	defer rows.Close()

	for rows.Next() {
		var hash, coinbase, parentHash, uncleHash, difficulty, size, rewards, num string
		var gasUsed, gasLimit, nonce uint64
		var txCount, uncleCount int
		var age time.Time

		err = rows.Scan(
			&hash, &coinbase, &gasUsed, &gasLimit, &txCount, &uncleCount, &age, &parentHash, &uncleHash, &difficulty, &size, &nonce, &rewards, &num)

		arr.Blocks = append(arr.Blocks, stypes.SBlock{
			Hash:       hash,
			Coinbase:   coinbase,
			GasUsed:    gasUsed,
			GasLimit:   gasLimit,
			TxCount:    txCount,
			UncleCount: uncleCount,
			Age:        age,
			ParentHash: parentHash,
			UncleHash:  uncleHash,
			Difficulty: difficulty,
			Size:       size,
			Nonce:      nonce,
			Rewards:    rewards,
			Number:     num,
		})

		blocks, _ := json.Marshal(arr.Blocks)
		blocksFmt := string(blocks)
		blockArr = blocksFmt
	}
	return blockArr
}

//GetBlock queries to send single block info
//TODO provide blockHash arg passed from handler.go
func SGetBlock(sqldb *sqlx.DB, blockNumber string) string {
	sqlStatement := `SELECT * FROM blocks WHERE number=$1;`
	tx, _ := sqldb.Begin()
	row := sqldb.QueryRow(sqlStatement, blockNumber)
	tx.Commit()
	var hash, coinbase, parentHash, uncleHash, difficulty, size, rewards, num string
	var gasUsed, gasLimit, nonce uint64
	var txCount, uncleCount int
	var age time.Time

	row.Scan(
		&hash, &coinbase, &gasUsed, &gasLimit, &txCount, &uncleCount, &age, &parentHash, &uncleHash, &difficulty, &size, &nonce, &rewards, &num)

	block := stypes.SBlock{
		Hash:       hash,
		Coinbase:   coinbase,
		GasUsed:    gasUsed,
		GasLimit:   gasLimit,
		TxCount:    txCount,
		UncleCount: uncleCount,
		Age:        age,
		ParentHash: parentHash,
		UncleHash:  uncleHash,
		Difficulty: difficulty,
		Size:       size,
		Nonce:      nonce,
		Rewards:    rewards,
		Number:     num,
	}
	json, _ := json.Marshal(block)
	return string(json)
}

func SGetRecentBlock(sqldb *sqlx.DB) string {
	sqlStatement := `SELECT * FROM blocks WHERE number=(SELECT MAX(number) FROM blocks);`
	tx, _ := sqldb.Begin()
	row := sqldb.QueryRow(sqlStatement)
	tx.Commit()
	var hash, coinbase, parentHash, uncleHash, difficulty, size, rewards, num string
	var gasUsed, gasLimit, nonce uint64
	var txCount, uncleCount int
	var age time.Time

	row.Scan(
		&hash, &coinbase, &gasUsed, &gasLimit, &txCount, &uncleCount, &age, &parentHash, &uncleHash, &difficulty, &size, &nonce, &rewards, &num)

	block := stypes.SBlock{
		Hash:       hash,
		Coinbase:   coinbase,
		GasUsed:    gasUsed,
		GasLimit:   gasLimit,
		TxCount:    txCount,
		UncleCount: uncleCount,
		Age:        age,
		ParentHash: parentHash,
		UncleHash:  uncleHash,
		Difficulty: difficulty,
		Size:       size,
		Nonce:      nonce,
		Rewards:    rewards,
		Number:     num,
	}
	json, _ := json.Marshal(block)
	return string(json)
}

//func SGetRecentBlockHash() string {
//	sqldb, _ := DBConnection()
//	sqlStatement := `SELECT hash FROM blocks WHERE number=(SELECT MAX(number) FROM blocks);`
//	tx, _ := sqldb.Begin()
//	row := sqldb.QueryRow(sqlStatement)
//	tx.Commit()
//	var hash string
//
//	row.Scan(
//		&hash)
//
//	blockhash := stypes.BlockHash{
//		Hash: hash,
//	}
//	json, _ := json.Marshal(blockhash)
//	return string(json)
//}

func SGetAllTransactionsFromBlock(sqldb *sqlx.DB, blockNumber string) string {
	var arr stypes.TxRes
	var txx string
	sqlStatement := `SELECT * FROM txs WHERE blocknumber=$1`
	tx, _ := sqldb.Begin()
	rows, err := sqldb.Query(sqlStatement, blockNumber)
	tx.Commit()
	if err != nil {
		fmt.Println("err")
	}
	defer rows.Close()
	for rows.Next() {
		var txhash, to_addr, from_addr, txfee, blockhash, blocknumber, amount, status string
		var gasprice, gas, gasLimit, nonce uint64
		var isContract bool
		var age time.Time
		var data []byte

		err = rows.Scan(
			&txhash, &to_addr, &from_addr, &blockhash, &blocknumber, &amount, &gasprice, &gas, &gasLimit, &txfee, &nonce, &status, &isContract, &age, &data,
		)

		arr.TxEntry = append(arr.TxEntry, stypes.ShyftTxEntryPretty{
			TxHash:      txhash,
			To:          to_addr,
			From:        from_addr,
			BlockHash:   blockhash,
			BlockNumber: blocknumber,
			Amount:      amount,
			GasPrice:    gasprice,
			Gas:         gas,
			GasLimit:    gasLimit,
			Cost:        txfee,
			Nonce:       nonce,
			Status:      status,
			IsContract:  isContract,
			Age:         age,
			Data:        data,
		})

		txData, _ := json.Marshal(arr.TxEntry)
		newtx := string(txData)
		txx = newtx
	}
	return txx
}

func SGetAllBlocksMinedByAddress(sqldb *sqlx.DB, coinbase string) string {
	var arr stypes.BlockRes
	var blockArr string
	sqlStatement := `SELECT * FROM blocks WHERE coinbase=$1`
	tx, _ := sqldb.Begin()
	rows, err := sqldb.Query(sqlStatement, coinbase)
	tx.Commit()
	if err != nil {
		fmt.Println("err")
	}
	defer rows.Close()

	for rows.Next() {
		var hash, coinbase, parentHash, uncleHash, difficulty, size, rewards, num string
		var gasUsed, gasLimit, nonce uint64
		var txCount, uncleCount int
		var age time.Time

		err = rows.Scan(
			&hash, &coinbase, &gasUsed, &gasLimit, &txCount, &uncleCount, &age, &parentHash, &uncleHash, &difficulty, &size, &nonce, &rewards, &num)

		arr.Blocks = append(arr.Blocks, stypes.SBlock{
			Hash:       hash,
			Coinbase:   coinbase,
			GasUsed:    gasUsed,
			GasLimit:   gasLimit,
			TxCount:    txCount,
			UncleCount: uncleCount,
			Age:        age,
			ParentHash: parentHash,
			UncleHash:  uncleHash,
			Difficulty: difficulty,
			Size:       size,
			Nonce:      nonce,
			Rewards:    rewards,
			Number:     num,
		})

		blocks, _ := json.Marshal(arr.Blocks)
		blocksFmt := string(blocks)
		blockArr = blocksFmt
	}
	return blockArr
}

//GetAllTransactions getter fn for API
func SGetAllTransactions(sqldb *sqlx.DB) string {
	var arr stypes.TxRes
	var txx string
	tx, _ := sqldb.Begin()
	rows, err := sqldb.Query(`SELECT * FROM txs`)
	tx.Commit()
	if err != nil {
		fmt.Println("err")
	}
	defer rows.Close()
	for rows.Next() {
		var txhash, to_addr, from_addr, txfee, blockhash, blocknumber, amount, status string
		var gasprice, gas, gasLimit, nonce uint64
		var isContract bool
		var age time.Time
		var data []byte

		err = rows.Scan(
			&txhash, &to_addr, &from_addr, &blockhash, &blocknumber, &amount, &gasprice, &gas, &gasLimit, &txfee, &nonce, &status, &isContract, &age, &data,
		)

		arr.TxEntry = append(arr.TxEntry, stypes.ShyftTxEntryPretty{
			TxHash:      txhash,
			To:          to_addr,
			From:        from_addr,
			BlockHash:   blockhash,
			BlockNumber: blocknumber,
			Amount:      amount,
			GasPrice:    gasprice,
			Gas:         gas,
			GasLimit:    gasLimit,
			Cost:        txfee,
			Nonce:       nonce,
			Status:      status,
			IsContract:  isContract,
			Age:         age,
			Data:        data,
		})

		txData, _ := json.Marshal(arr.TxEntry)
		newtx := string(txData)
		txx = newtx
	}
	return txx
}

//GetTransaction fn returns single tx
func SGetTransaction(sqldb *sqlx.DB, txHash string) string {
	sqlStatement := `SELECT * FROM txs WHERE txhash=$1;`
	tx, _ := sqldb.Begin()
	row := sqldb.QueryRow(sqlStatement, txHash)
	tx.Commit()
	var txhash, to_addr, from_addr, txfee, blockhash, blocknumber, amount, status string
	var gasprice, gas, gasLimit, nonce uint64
	var isContract bool
	var age time.Time
	var data []byte

	row.Scan(
		&txhash, &to_addr, &from_addr, &blockhash, &blocknumber, &amount, &gasprice, &gas, &gasLimit, &txfee, &nonce, &status, &isContract, &age, &data)

	txData := stypes.ShyftTxEntryPretty{
		TxHash:      txhash,
		To:          to_addr,
		From:        from_addr,
		BlockHash:   blockhash,
		BlockNumber: blocknumber,
		Amount:      amount,
		GasPrice:    gasprice,
		Gas:         gas,
		GasLimit:    gasLimit,
		Cost:        txfee,
		Nonce:       nonce,
		Status:      status,
		IsContract:  isContract,
		Age:         age,
		Data:        data,
	}
	json, _ := json.Marshal(txData)

	return string(json)
}

func InnerSGetAccount(sqldb *sqlx.DB, address string) (stypes.SAccounts, bool) {
	sqlStatement := `SELECT * FROM accounts WHERE addr=$1;`
	var addr, balance, nonce string
	tx, _ := sqldb.Begin()
	err := sqldb.QueryRow(sqlStatement, address).Scan(&addr, &balance, &nonce)
	tx.Commit()
	if err == sql.ErrNoRows {
		return stypes.SAccounts{}, false
	} else {
		account := stypes.SAccounts{
			Addr:         addr,
			Balance:      balance,
			AccountNonce: nonce,
		}
		return account, true
	}
}

//GetAccount returns account balances
func SGetAccount(sqldb *sqlx.DB, address string) string {
	var account, _ = InnerSGetAccount(sqldb, address)
	json, _ := json.Marshal(account)
	return string(json)
}

//GetAllAccounts returns all accounts and balances
func SGetAllAccounts(sqldb *sqlx.DB) string {
	var array stypes.AccountRes
	var accountsArr, nonce string
	tx, _ := sqldb.Begin()
	accs, err := sqldb.Query(`
		SELECT
			addr,
			balance,
			nonce
		FROM accounts`)
	tx.Commit()
	if err != nil {
		fmt.Println(err)
	}

	defer accs.Close()

	for accs.Next() {
		var addr, balance string
		err = accs.Scan(
			&addr, &balance, &nonce,
		)

		array.AllAccounts = append(array.AllAccounts, stypes.SAccounts{
			Addr:         addr,
			Balance:      balance,
			AccountNonce: nonce,
		})

		accounts, _ := json.Marshal(array.AllAccounts)
		accountsFmt := string(accounts)
		accountsArr = accountsFmt
	}
	return accountsArr
}

//GetAccount returns account balances
func SGetAccountTxs(sqldb *sqlx.DB, address string) string {
	var arr stypes.TxRes
	var txx string
	sqlStatement := `SELECT * FROM txs WHERE to_addr=$1 OR from_addr=$1;`
	tx, _ := sqldb.Begin()
	rows, err := sqldb.Query(sqlStatement, address)
	tx.Commit()
	if err != nil {
		fmt.Println("err", err)
	}
	defer rows.Close()
	for rows.Next() {
		var txhash, to_addr, from_addr, txfee, blockhash, blocknumber, amount, status string
		var gasprice, gas, gasLimit, nonce uint64
		var isContract bool
		var age time.Time
		var data []byte

		err = rows.Scan(
			&txhash, &to_addr, &from_addr, &blockhash, &blocknumber, &amount, &gasprice, &gas, &gasLimit, &txfee, &nonce, &status, &isContract, &age, &data,
		)

		arr.TxEntry = append(arr.TxEntry, stypes.ShyftTxEntryPretty{
			TxHash:      txhash,
			To:          to_addr,
			From:        from_addr,
			BlockHash:   blockhash,
			BlockNumber: blocknumber,
			Amount:      amount,
			GasPrice:    gasprice,
			Gas:         gas,
			GasLimit:    gasLimit,
			Cost:        txfee,
			Nonce:       nonce,
			Status:      status,
			IsContract:  isContract,
			Age:         age,
			Data:        data,
		})

		txData, _ := json.Marshal(arr.TxEntry)
		newtx := string(txData)
		txx = newtx
	}
	return txx
}

//GetAllInternalTransactions getter fn for API
func SGetAllInternalTransactions(sqldb *sqlx.DB) string {
	var arr stypes.InternalArray
	var internaltx string
	tx, _ := sqldb.Begin()
	rows, err := sqldb.Query(`SELECT * FROM internaltxs`)
	tx.Commit()
	if err != nil {
		fmt.Println("err")
	}
	defer rows.Close()
	for rows.Next() {
		var txhash, blockhash, action, to_addr, from_addr, amount, input, output string
		var gas, gasUsed uint64
		var id int
		var age string

		err = rows.Scan(
			&id, &txhash, &blockhash, &action, &to_addr, &from_addr, &amount, &gas, &gasUsed, &age, &input, &output,
		)

		arr.InternalEntry = append(arr.InternalEntry, stypes.InteralWrite{
			ID:      id,
			Hash:    txhash,
			BlockHash: blockhash,
			Action:  action,
			To:      to_addr,
			From:    from_addr,
			Value:   amount,
			Gas:     gas,
			GasUsed: gasUsed,
			Time:    age,
			Input:   input,
			Output:  output,
		})

		txData, _ := json.Marshal(arr.InternalEntry)
		newtx := string(txData)
		internaltx = newtx
	}
	return internaltx
}

//GetInternalTransaction fn returns single tx
func SGetInternalTransaction(sqldb *sqlx.DB, txHash string) string {
	var arr stypes.InternalArray
	var internaltx string

	sqlStatement := `SELECT * FROM internaltxs WHERE txhash=$1;`
	tx, _ := sqldb.Begin()
	rows, err := sqldb.Query(sqlStatement, txHash)
	tx.Commit()
	if err != nil {
		fmt.Println("err")
	}
	defer rows.Close()

	for rows.Next() {
		var txhash, blockhash, action, to_addr, from_addr, amount, input, output string
		var id int
		var gas, gasUsed uint64
		var age string

		err = rows.Scan(
			&id, &txhash, &blockhash, &action, &to_addr, &from_addr, &amount, &gas, &gasUsed, &age, &input, &output,
		)

		arr.InternalEntry = append(arr.InternalEntry, stypes.InteralWrite{
			ID:      id,
			Hash:    txhash,
			BlockHash: blockhash,
			Action:  action,
			To:      to_addr,
			From:    from_addr,
			Value:   amount,
			Gas:     gas,
			GasUsed: gasUsed,
			Time:    age,
			Input:   input,
			Output:  output,
		})

		txData, _ := json.Marshal(arr.InternalEntry)
		newtx := string(txData)
		internaltx = newtx
	}
	return internaltx
}
