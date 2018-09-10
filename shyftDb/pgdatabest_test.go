package shyftdb_test

import (
	"database/sql"
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/ShyftNetwork/go-empyrean/core"
	"github.com/ShyftNetwork/go-empyrean/shyft_schema"
	"github.com/jmoiron/sqlx"
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

func DeleteAllTables(db *sqlx.DB) {
	db.MustExec("DELETE FROM accounts;")
	db.MustExec("DELETE FROM accountblocks;")
	db.MustExec("DELETE FROM blocks;")
	db.MustExec("DELETE FROM txs;")
	db.MustExec("DELETE FROM internalTxs")
}

func TestDbApi(t *testing.T) {
	t.Run("CreateAccount - creates an account and an accountblock record in the PG db ", func(t *testing.T) {
		db, err := core.InitDB()
		DeleteAllTables(db)
		addr := "0x7ef5a6135f1fd6a02593eedc869c6d41d934aef8"
		balance, _ := new(big.Int).SetString("3500000000", 10)
		accountNonce := strconv.Itoa(int(1))
		blockHash := "0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056"
		core.CreateAccount(addr, balance.String(), accountNonce, blockHash)
		newDbAccounts := []shyftschema.Account{}
		err = db.Select(&newDbAccounts, "SELECT * FROM accounts WHERE addr = $1", addr)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("account written: %+v\n", newDbAccounts[0])
		newDbAccountBlocks := []shyftschema.AccountBlock{}
		err = db.Select(&newDbAccountBlocks, "SELECT * FROM accountblocks")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("accountBlock written: %+v\n", newDbAccountBlocks[0])

		if len(newDbAccounts) > 1 {
			t.Errorf("Got %v Accounts Created: Expected 1", len(newDbAccounts))
		}
		if len(newDbAccountBlocks) > 1 {
			t.Errorf("Got %v Accounts Created: Expected 1", len(newDbAccounts))
		}
		stringBalance := strconv.FormatInt(newDbAccounts[0].Balance, 10)
		stringDelta := strconv.FormatInt(newDbAccountBlocks[0].Delta, 10)
		if newDbAccounts[0].Addr != addr || stringBalance != "3500000000" || accountNonce != "1" {
			t.Errorf("Account: Got %v Accounts Created: Expected addr: %s balance: %d nonce %s", newDbAccounts, addr, balance, accountNonce)
		}
		if newDbAccountBlocks[0].Acct != addr || stringDelta != "3500000000" || blockHash != newDbAccountBlocks[0].Blockhash {
			t.Errorf("AccountBlocks: Got %v Accounts Created: Expected acct: %s blockHash: %s delta %d", newDbAccountBlocks, addr, blockHash, balance)

		}
	})
	t.Run("UpdateAccount - updates an account and updates the blockHash delta accordingly", func(t *testing.T) {
		db, err := core.InitDB()
		DeleteAllTables(db)
		addr := "0x7ef5a6135f1fd6a02593eedc869c6d41d934aef8"
		balance, _ := new(big.Int).SetString("3500000000", 10)
		accountNonce := strconv.Itoa(int(1))
		blockHash := "0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056"
		core.CreateAccount(addr, balance.String(), accountNonce, blockHash)
		newDbAccounts := []shyftschema.Account{}
		err = db.Select(&newDbAccounts, "SELECT * FROM accounts WHERE addr = $1", addr)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("account written: %+v\n", newDbAccounts[0])
		newDbAccountBlocks := []shyftschema.AccountBlock{}
		err = db.Select(&newDbAccountBlocks, "SELECT * FROM accountblocks")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("accountBlock written: %+v\n", newDbAccountBlocks[0])

		if len(newDbAccounts) > 1 {
			t.Errorf("Got %v Accounts Created: Expected 1", len(newDbAccounts))
		}
		if len(newDbAccountBlocks) > 1 {
			t.Errorf("Got %v Accounts Created: Expected 1", len(newDbAccounts))
		}
		stringBalance := strconv.FormatInt(newDbAccounts[0].Balance, 10)
		stringDelta := strconv.FormatInt(newDbAccountBlocks[0].Delta, 10)
		if newDbAccounts[0].Addr != addr || stringBalance != "3500000000" || accountNonce != "1" {
			t.Errorf("Account: Got %v Accounts Created: Expected addr: %s balance: %d nonce %s", newDbAccounts, addr, balance, accountNonce)
		}
		if newDbAccountBlocks[0].Acct != addr || stringDelta != "3500000000" || blockHash != newDbAccountBlocks[0].Blockhash {
			t.Errorf("AccountBlocks: Got %v Accounts Created: Expected acct: %s blockHash: %s delta %d", newDbAccountBlocks, addr, blockHash, balance)

		}
	})
}
