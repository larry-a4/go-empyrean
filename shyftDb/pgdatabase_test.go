package shyftdb_test

import (
	"database/sql"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ShyftNetwork/go-empyrean/common"
	"github.com/ShyftNetwork/go-empyrean/core"
	stypes "github.com/ShyftNetwork/go-empyrean/core/sTypes"
	"github.com/ShyftNetwork/go-empyrean/core/types"
	"github.com/ShyftNetwork/go-empyrean/crypto"

	"github.com/ShyftNetwork/go-empyrean/shyft_schema"
	"github.com/jmoiron/sqlx"
)

type ShyftTracer struct{}

const (
	testAddress = "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"
)

var tx, _ = types.NewTransaction(
	3,
	common.HexToAddress("b94f5374fce5edbc8e2a8697c15331677e6ebf0b"),
	big.NewInt(10),
	2000,
	big.NewInt(1),
	common.FromHex("5544"),
).WithSignature(
	types.HomesteadSigner{},
	common.Hex2Bytes("98ff921201554726367d2be8c804a7ff89ccf285ebc57dff8ae4c44b9c19ac4a8887321be575c8095f789dd4c743dfe42c1820f9231f98a962b210e3ac2452a301"),
)

func TestDbCreationExistence(t *testing.T) {
	core.DeletePgDb(core.DbName())
	db, err := core.InitDB()
	if err != nil || err == sql.ErrNoRows {
		fmt.Println(err)
	}
	t.Run("Creates the PG DB if it Doesnt Exist", func(t *testing.T) {
		_, err = core.DbExists(core.DbName())
		if err != nil || err == sql.ErrNoRows {
			t.Errorf("Error in Database Connection - DB doesn't Exist - %s", err)
		}
	})
	t.Run("Creates the Tables Required from the Migration Schema", func(t *testing.T) {
		db, err := core.InitDB()
		if err != nil || err == sql.ErrNoRows {
			fmt.Println(err)
		}
		tableNameQuery := `select table_name from information_schema.tables where table_schema = 'public' AND table_type = 'BASE TABLE' order by table_name ASC;`
		db = core.Connect(core.ShyftConnectStr())
		rows, err := db.Query(tableNameQuery)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		var tablenames string
		var table string
		notLast := rows.Next()
		for notLast {
			//... rows.Scan
			err = rows.Scan(&table)
			if err != nil {
				panic(err)
			}
			notLast = rows.Next()
			if notLast {
				tablenames += table + ", "
			} else {
				tablenames += table
			}
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}
		want := "accountblocks, accounts, blocks, internaltxs, txs"
		if tablenames != want {
			t.Errorf("Test Failed as wanted: %s  - got: %s", want, tablenames)
		}
	})
	core.DeletePgDb(core.DbName())
	db, err = core.InitDB()
	if err != nil || err == sql.ErrNoRows {
		fmt.Println(err)
	}
	db.Close()
}

func deleteAllTables(db *sqlx.DB) {
	db.MustExec("DELETE FROM accounts;")
	db.MustExec("DELETE FROM accountblocks;")
	db.MustExec("DELETE FROM blocks;")
	db.MustExec("DELETE FROM txs;")
	db.MustExec("DELETE FROM internalTxs")
}

func TestCreateAccount(t *testing.T) {
	t.Run("CreateAccount - creates an account in the PG db ", func(t *testing.T) {
		db, err := core.InitDB()
		deleteAllTables(db)
		addr := "0x7ef5a6135f1fd6a02593eedc869c6d41d934aef8"
		balance, _ := new(big.Int).SetString("3500000000", 10)
		accountNonce := strconv.Itoa(int(1))
		err = core.CreateAccount(addr, balance.String(), accountNonce)
		if err != nil {
			fmt.Println(err)
			return
		}
		newDbAccounts := []shyftschema.Account{}
		err = db.Select(&newDbAccounts, "SELECT * FROM accounts WHERE addr = $1", addr)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("account written: %+v\n", newDbAccounts[0])
		if len(newDbAccounts) > 1 {
			t.Errorf("Got %v Accounts Created: Expected 1", len(newDbAccounts))
		}
		stringBalance := strconv.FormatInt(newDbAccounts[0].Balance, 10)
		if newDbAccounts[0].Addr != addr || stringBalance != "3500000000" || accountNonce != "1" {
			t.Errorf("Account: Got %v Accounts Created: Expected addr: %s balance: %d nonce %s", newDbAccounts, addr, balance, accountNonce)
		}
	})
}

func TestInsertTx(t *testing.T) {
	// Set up a  test transaction
	key, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	signer := types.NewEIP155Signer(big.NewInt(2147483647))

	//Nonce, To Address,Value, GasLimit, Gasprice, data
	tx1 := types.NewTransaction(1, common.BytesToAddress([]byte{0x11}), big.NewInt(5), 1111, big.NewInt(11111), []byte{0x11, 0x11, 0x11})
	tx, _ := types.SignTx(tx1, signer, key)

	blockHash := "0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056"
	txData := stypes.ShyftTxEntryPretty{
		TxHash:      tx.Hash().Hex(),
		From:        tx.From().Hex(),
		To:          tx.To().String(),
		BlockHash:   blockHash,
		BlockNumber: strconv.Itoa(21234),
		Amount:      tx.Value().String(),
		Cost:        tx.Cost().Uint64(),
		GasPrice:    tx.GasPrice().Uint64(),
		GasLimit:    uint64(18000),
		Gas:         tx.Gas(),
		Nonce:       tx.Nonce(),
		Age:         time.Now(),
		Data:        tx.Data(),
		Status:      "SUCCESS",
		IsContract:  false,
	}
	fmt.Printf("Transaction To Be Inserted:\n  %+v \n", txData)
	t.Run("InsertTx - No Account exists inserts a transaction to the database and updates/creates accounts accordingly", func(t *testing.T) {
		db, _ := core.InitDB()
		deleteAllTables(db)

		core.InsertTx(txData)
		dbTransactions := []shyftschema.PgTransaction{}
		err := db.Select(&dbTransactions, "SELECT * FROM txs")
		if err != nil {
			panic(err)
		}
		pgdb := dbTransactions[0]
		txInput := txData
		if len(dbTransactions) != 1 {
			t.Errorf("Got %v db transactions created \nExpected 1", len(dbTransactions))
		}
		if pgdb.TxHash != txInput.TxHash && pgdb.Blockhash != txData.BlockHash && pgdb.To != txData.To && pgdb.From != txData.From && pgdb.Amount != txData.Amount {
			t.Errorf("Got %+v \nExpected %+v", dbTransactions[0], txData)
		}
		newDbAccounts := []shyftschema.Account{}
		err = db.Select(&newDbAccounts, "SELECT * FROM accounts")
		if err != nil {
			panic(err)
		}
		if len(newDbAccounts) != 2 {
			t.Errorf("Got %v db transactions created -  Expected 2", len(newDbAccounts))
		}
		toAcct := newDbAccounts[0]
		fromAcct := newDbAccounts[1]
		if toAcct.Addr != txData.To && new(big.Int).SetInt64(toAcct.Balance) != tx.Value() && toAcct.Nonce != 1 {
			t.Errorf("Got %+v \nExpected %s %s %d", toAcct, txData.To, txData.Amount, 1)
		}
		fromAcctBal, _ := strconv.Atoi(txData.Amount)
		fromBalInt := -1 * fromAcctBal
		product := new(big.Int)
		product.Mul(new(big.Int).SetInt64(-1), tx.Value())
		if fromAcct.Addr != txData.To && new(big.Int).SetInt64(fromAcct.Balance) != product &&
			fromAcct.Nonce != 1 {
			t.Errorf("Got %+v \nExpected %s %d %d", fromAcct, txData.From, fromBalInt, 1)
		}
		newDbAccountBlocks := []shyftschema.AccountBlock{}
		err = db.Select(&newDbAccountBlocks, "SELECT * FROM accountblocks")
		if err != nil {
			panic(err)
		}
		if len(newDbAccountBlocks) != 2 {
			t.Errorf("Got %v db accountblocks created -  Expected 2", len(dbTransactions))
		}
		toAcctBl := newDbAccountBlocks[0]
		fromAcctBl := newDbAccountBlocks[1]
		if toAcctBl.Acct != txData.To && toAcctBl.Blockhash != txData.BlockHash &&
			strconv.Itoa(int(toAcctBl.Delta)) != txData.Amount && int(toAcctBl.TxCount) != 1 {
			t.Errorf("Got %+v \nExpected %s %s %s", toAcctBl, txData.To, txData.BlockHash, txData.Amount)
		}
		if fromAcctBl.Acct != txData.To && fromAcctBl.Blockhash != txData.BlockHash &&
			strconv.Itoa(int(fromAcctBl.Delta)*-1) != txData.Amount && int(toAcctBl.TxCount) != 1 {
			t.Errorf("Got %+v \nExpected %s %s %s", toAcctBl, txData.To, txData.BlockHash, txData.Amount)
		}
	})
	//TODO: Add tests for:
	//         Multiple Transactions re AccountBlock Generation
	//         Genesis Block - correct setting of pg tables
	//         Rollback
}

func TestGenesisBlockCreation(t *testing.T) {
	// db, _ := core.InitDB()
	// deleteAllTables(db)
	// edb, _ := eth.NewShyftTestLDB()
	// shyftTracer := new(eth.ShyftTracer)
	// core.SetIShyftTracer(shyftTracer)

	// ethConf := &eth.Config{
	// 	Genesis:   core.DeveloperGenesisBlock(15, common.Address{}),
	// 	Etherbase: common.HexToAddress(testAddress),
	// 	Ethash: ethash.Config{
	// 		PowMode: ethash.ModeTest,
	// 	},
	// }

	// eth.SetGlobalConfig(ethConf)

	// eth.InitTracerEnv()
	// // key, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	// // signer := types.NewEIP155Signer(big.NewInt(2147483647))
	// t.Run("SetupGenesisBlock - populates the pg accounts, transactions, and accountblocks appropriately", func(t *testing.T) {
	// 	deleteAllTables(db)
	// 	core.SetupGenesisBlock(edb, ethConf.Genesis)
	// 	newDbAccounts := []shyftschema.Account{}
	// 	err := db.Select(&newDbAccounts, "SELECT * FROM accounts")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if len(newDbAccounts) != 2 {
	// 		t.Errorf("Got %v db transactions created -  Expected 2", len(newDbAccounts))
	// 	}
	// })
}
