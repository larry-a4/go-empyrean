package core

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/ShyftNetwork/go-empyrean/core/types"
	"github.com/ShyftNetwork/go-empyrean/core/sTypes"
	"encoding/json"
	"strings"
	"github.com/ShyftNetwork/go-empyrean/ethdb"
	"fmt"
)

type ShyftTracer struct{}



var SQLDB *sqlx.DB

//var tx, _ = types.NewTransaction(3, common.HexToAddress("b94f5374fce5edbc8e2a8697c15331677e6ebf0b"), big.NewInt(10), 2000, big.NewInt(1), common.FromHex("5544"), ).WithSignature(types.HomesteadSigner{}, common.Hex2Bytes("98ff921201554726367d2be8c804a7ff89ccf285ebc57dff8ae4c44b9c19ac4a8887321be575c8095f789dd4c743dfe42c1820f9231f98a962b210e3ac2452a301"), )

//func init() {
//	runtime.LockOSThread()
//}

//func TestMain(m *testing.M) {
//	go func() {
//		os.Exit(m.Run())
//	}()
//	sqldb, err := InitDB()
//	if err != nil {
//		panic(err)
//	}
//	SQLDB = sqldb
//}

func SetUpTracerEnv() {
	//eth.NewShyftTestLDB()
	//shyftTracer := new(eth.ShyftTracer)
	//SetIShyftTracer(shyftTracer)
	//
	//ethConf := &eth.Config{
	//	Genesis:   DeveloperGenesisBlock(15, common.Address{}),
	//	Etherbase: common.HexToAddress(testAddress),
	//	Ethash: ethash.Config{
	//		PowMode: ethash.ModeTest,
	//	},
	//}
	//
	//eth.SetGlobalConfig(ethConf)
	//eth.InitTracerEnv()
}

func TestBlock(t *testing.T) {
	//SET UP FOR TEST FUNCTIONS
	//SetUpTracerEnv()

	shyftdb, err := ethdb.NewShyftDatabase()
	if err != nil {
		fmt.Println(err)
	}
	shyftdb.TruncateTables()

	fromAddr := "0x71562b71999873db5b286df957af199ec94617f7"
	shyftdb.CreateAccount(fromAddr, "201", "1")

	t.Run("TestBlockToReturnBlock", func(t *testing.T) {
		for _, bl := range CreateTestBlocks() {
			// Write and verify the block in the database
			if err := SWriteBlock(shyftdb, bl, CreateTestReceipts()); err != nil {
				t.Fatalf("Failed to write block into database: %v", err)
			}
		}
		blocks := CreateTestBlocks()

		entry, _ := ethdb.SGetBlock(blocks[0].Number().String())
		byt := []byte(entry)
		var data stypes.SBlock
		json.Unmarshal(byt, &data)

		//TODO Difficulty, rewards, age
		if blocks[0].Hash().String() != data.Hash {
			t.Fatalf("Block Hash [%v]: Block hash not found", blocks[0].Hash().String())
		}
		if blocks[0].Coinbase().String() != data.Coinbase {
			t.Fatalf("Block coinbase [%v]: Block coinbase not found", blocks[0].Coinbase().String())
		}
		if blocks[0].Number().String() != data.Number {
			t.Fatalf("Block number [%v]: Block number not found", blocks[0].Number().String())
		}
		if blocks[0].GasUsed() != data.GasUsed {
			t.Fatalf("Gas Used [%v]: Gas used not found", blocks[0].GasUsed())
		}
		if blocks[0].GasLimit() != data.GasLimit {
			t.Fatalf("Gas Limit [%v]: Gas limit not found", blocks[0].GasLimit())
		}
		if blocks[0].Transactions().Len() != data.TxCount {
			t.Fatalf("Tx Count [%v]: Tx Count not found", blocks[0].Transactions().Len())
		}
		if len(blocks[0].Uncles()) != data.UncleCount {
			t.Fatalf("Uncle count [%v]: Uncle count not found", len(blocks[0].Uncles()))
		}
		if blocks[0].ParentHash().String() != data.ParentHash {
			t.Fatalf("Parent hash [%v]: Parent hash not found", blocks[0].ParentHash().String())
		}
		if blocks[0].UncleHash().String() != data.UncleHash {
			t.Fatalf("Uncle hash [%v]: Uncle hash not found", blocks[0].UncleHash().String())
		}
		if blocks[0].Size().String() != data.Size {
			t.Fatalf("Size [%v]: Size not found", blocks[0].Size().String())
		}
		if blocks[0].Nonce() != data.Nonce {
			t.Fatalf("Block nonce [%v]: Block nonce not found", blocks[0].Nonce())
		}

		if getAllBlocks, _ := ethdb.SGetAllBlocks(); len(getAllBlocks) == 0 {
			t.Fatalf("GetAllBlocks [%v]: GetAllBlocks did not return correctly", getAllBlocks)
		}

		if getAllBlocksMinedByAddress, _ := ethdb.SGetAllBlocksMinedByAddress(blocks[0].Coinbase().String()); len(getAllBlocksMinedByAddress) == 0 {
			t.Fatalf("GetAllBlocksMinedByAddress [%v]: GetAllBlocksMinedByAddress did not return correctly", getAllBlocksMinedByAddress)
		}
	})
	t.Run("TestGetRecentBlock", func(t *testing.T) {
		response, _ := ethdb.SGetRecentBlock()
		byteRes := []byte(response)
		var recentBlock stypes.SBlock
		json.Unmarshal(byteRes, &recentBlock)

		blocks := CreateTestBlocks()

		if blocks[0].Hash().String() != recentBlock.Hash {
			t.Fatalf("Block Hash [%v]: Block hash not found, Expected: [%v]", blocks[0].Hash().String(), recentBlock.Hash)
		}
		if blocks[0].Coinbase().String() != recentBlock.Coinbase {
			t.Fatalf("Block coinbase [%v]: Block coinbase not found", blocks[0].Coinbase().String())
		}
		if blocks[0].Number().String() != recentBlock.Number {
			t.Fatalf("Block number [%v]: Block number not found", blocks[0].Number().String())
		}
		if blocks[0].GasUsed() != recentBlock.GasUsed {
			t.Fatalf("Gas Used [%v]: Gas used not found", blocks[0].GasUsed())
		}
		if blocks[0].GasLimit() != recentBlock.GasLimit {
			t.Fatalf("Gas Limit [%v]: Gas limit not found", blocks[0].GasLimit())
		}
		if blocks[0].Transactions().Len() != recentBlock.TxCount {
			t.Fatalf("Tx Count [%v]: Tx Count not found", blocks[0].Transactions().Len())
		}
		if len(blocks[0].Uncles()) != recentBlock.UncleCount {
			t.Fatalf("Uncle count [%v]: Uncle count not found", len(blocks[0].Uncles()))
		}
		if blocks[0].ParentHash().String() != recentBlock.ParentHash {
			t.Fatalf("Parent hash [%v]: Parent hash not found", blocks[0].ParentHash().String())
		}
		if blocks[0].UncleHash().String() != recentBlock.UncleHash {
			t.Fatalf("Uncle hash [%v]: Uncle hash not found", blocks[0].UncleHash().String())
		}
		if blocks[0].Size().String() != recentBlock.Size {
			t.Fatalf("Size [%v]: Size not found", blocks[0].Size().String())
		}
		if blocks[0].Nonce() != recentBlock.Nonce {
			t.Fatalf("Block nonce [%v]: Block nonce not found", blocks[0].Nonce())
		}

		if allTxsFromBlock, _ := ethdb.SGetAllTransactionsFromBlock(blocks[2].Number().String()); len(allTxsFromBlock) == 0 {
			t.Fatalf("GetAllTransactionsFromBlock [%v]: GetAllTransactionsFromBlock did not return correctly", allTxsFromBlock)
		}
	})

	t.Run("TestContractCreationTx", func(t *testing.T) {
		var contractAddressFromReciept string
		for _, receipt := range CreateTestReceipts() {
			contractAddressFromReciept = (*types.ReceiptForStorage)(receipt).ContractAddress.String()
		}

		blocks := CreateTestBlocks()

		for _, tx := range CreateTestContractTransactions() {
			txn, _ := ethdb.SGetTransaction(tx.Hash().String())
			byt := []byte(txn)
			var data stypes.ShyftTxEntryPretty
			json.Unmarshal(byt, &data)

			if tx.Hash().String() != data.TxHash {
				t.Fatalf("txHash [%v]: tx Hash not found, expected [%v] ::", tx.Hash().String(), data.TxHash)
			}
			if contractAddressFromReciept != data.To {
				t.Fatalf("Contract Addr [%v]: Contract addr not found", contractAddressFromReciept)
			}
			if strings.ToLower(tx.From().String()) != data.From {
				t.Fatalf("From Addr [%v]: From addr not found", tx.From().String())
			}
			if tx.Nonce() != data.Nonce {
				t.Fatalf("Nonce [%v]: Nonce not found", tx.Nonce())
			}
			if tx.Gas() != data.Gas {
				t.Fatalf("Gas [%v]: Gas not found", tx.Gas())
			}
			if tx.GasPrice().Uint64() != data.GasPrice {
				t.Fatalf("Gas Price [%v]: Gas price not found", tx.GasPrice().String())
			}
			if blocks[0].GasLimit() != data.GasLimit {
				t.Fatalf("Gas Limit [%v]: Gas limit not found", blocks[0].GasLimit())
			}
			if blocks[2].Hash().String() != data.BlockHash {
				t.Fatalf("Block Hash [%v]: Block hash not found", blocks[2].Hash().String())
			}
			if blocks[2].Number().String() != data.BlockNumber {
				t.Fatalf("Block Number [%v]: Block number not found", blocks[2].Number().String())
			}
			if tx.Value().String() != data.Amount {
				t.Fatalf("Amount [%v]: Amount not found", tx.Value().String())
			}
			if tx.Cost().String() != data.Cost {
				t.Fatalf("Cost [%v]: Cost not found", tx.Cost().String())
			}
			var status string
			if CreateTestReceipts()[0].Status == 1 {
				status = "SUCCESS"
			}
			if CreateTestReceipts()[0].Status == 0 {
				status = "FAIL"
			}
			if status != data.Status {
				t.Fatalf("Receipt status [%v]: Receipt status not found", status)
			}
			var isContract bool
			if tx.To() != nil {
				isContract = false
			} else {
				isContract = true
			}
			if isContract != data.IsContract {
				t.Fatalf("isContract [%v]: isContract bool is incorrect", isContract)
			}
		}
	})

	t.Run("TestTransactionsToReturnTransactions", func(t *testing.T) {
		for _, tx := range CreateTestTransactions() {
			txn, _ := ethdb.SGetTransaction(tx.Hash().String())
			byt := []byte(txn)
			var data stypes.ShyftTxEntryPretty
			json.Unmarshal(byt, &data)

			blocks := CreateTestBlocks()

			//TODO age, data
			if strings.ToLower(tx.Hash().String()) != data.TxHash {
				t.Fatalf("txHash [%v]: tx Hash not found, expected [%v] ::", tx.Hash().String(), data.TxHash)
			}
			if strings.ToLower(tx.From().String()) != data.From {
				t.Fatalf("From Addr [%v]: From addr not found", tx.From().String())
			}
			if strings.ToLower(tx.To().String()) != data.To {
				t.Fatalf("To Addr [%v]: To addr not found", tx.To().String())
			}
			if tx.Nonce() != data.Nonce {
				t.Fatalf("Nonce [%v]: Nonce not found", tx.Nonce())
			}
			if tx.Gas() != data.Gas {
				t.Fatalf("Gas [%v]: Gas not found", tx.Gas())
			}
			if tx.GasPrice().Uint64() != data.GasPrice {
				t.Fatalf("Gas Price [%v]: Gas price not found", tx.GasPrice().String())
			}
			if blocks[0].GasLimit() != data.GasLimit {
				t.Fatalf("Gas Limit [%v]: Gas limit not found", blocks[0].GasLimit())
			}
			if blocks[0].Hash().String() != data.BlockHash {
				t.Fatalf("Block Hash [%v]: Block hash not found", blocks[0].Hash().String())
			}
			if blocks[0].Number().String() != data.BlockNumber {
				t.Fatalf("Block Number [%v]: Block number not found", blocks[0].Number().String())
			}
			if tx.Value().String() != data.Amount {
				t.Fatalf("Amount [%v]: Amount not found", tx.Value().String())
			}
			if tx.Cost().String() != data.Cost {
				t.Fatalf("Cost [%v]: Cost not found", tx.Cost().String())
			}
			var status string
			if CreateTestReceipts()[0].Status == 1 {
				status = "SUCCESS"
			}
			if CreateTestReceipts()[0].Status == 0 {
				status = "FAIL"
			}
			if status != data.Status {
				t.Fatalf("Receipt status [%v]: Receipt status not found", status)
			}
			var isContract bool
			if tx.To() != nil {
				isContract = false
			} else {
				isContract = true
			}
			if isContract != data.IsContract {
				t.Fatalf("isContract [%v]: isContract bool is incorrect", isContract)
			}
		}
		if getAllTx, _ := ethdb.SGetAllTransactions(); len(getAllTx) == 0 {
			t.Fatalf("GetAllTransactions [%v]: GetAllTransactions did not return correctly", getAllTx)
		}
	})
}
//
//func TestDbCreationExistence(t *testing.T) {
//	//db, err := core.InitDB()
//
//	t.Run("Creates the Tables Required from the Migration Schema", func(t *testing.T) {
//		tableNameQuery := `select table_name from information_schema.tables where table_schema = 'public' AND table_type = 'BASE TABLE' order by table_name ASC;`
//		SQLDB = core.Connect(core.ShyftConnectStr())
//		rows, err := SQLDB.Query(tableNameQuery)
//		if err != nil {
//			panic(err)
//		}
//		defer rows.Close()
//		var tablenames string
//		var table string
//		notLast := rows.Next()
//		for notLast {
//			//... rows.Scan
//			err = rows.Scan(&table)
//			if err != nil {
//				panic(err)
//			}
//			notLast = rows.Next()
//			if notLast {
//				tablenames += table + ", "
//			} else {
//				tablenames += table
//			}
//		}
//		err = rows.Err()
//		if err != nil {
//			panic(err)
//		}
//		want := "accountblocks, accounts, blocks, internaltxs, txs"
//		if tablenames != want {
//			t.Errorf("Test Failed as wanted: %s  - got: %s", want, tablenames)
//		}
//	})
//	SQLDB.Close()
//}
//
//func deleteAllTables(db *sqlx.DB) {
//	db.MustExec("TRUNCATE accounts CASCADE;")
//	db.MustExec("TRUNCATE accountblocks CASCADE;")
//	db.MustExec("TRUNCATE blocks CASCADE;")
//	db.MustExec("TRUNCATE txs CASCADE;")
//	db.MustExec("TRUNCATE internalTxs CASCADE")
//}

//func TestCreateAccount(t *testing.T) {
//	t.Run("CreateAccount - creates an account in the PG db ", func(t *testing.T) {
//		db, err := core.InitDB()
//		addr := "0x7ef5a6135f1fd6a02593eedc869c6d41d934aef8"
//		balance, _ := new(big.Int).SetString("3500000000", 10)
//		accountNonce := strconv.Itoa(int(1))
//		err = core.CreateAccount(addr, balance.String(), accountNonce)
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//		newDbAccounts := []shyftschema.Account{}
//		err = db.Select(&newDbAccounts, "SELECT * FROM accounts WHERE addr = $1", addr)
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//		if len(newDbAccounts) > 1 {
//			t.Errorf("Got %v Accounts Created: Expected 1", len(newDbAccounts))
//		}
//		//stringBalance := strconv.FormatInt(newDbAccounts[0].Balance, 10)
//		stringBalance := newDbAccounts[0].Balance
//		if newDbAccounts[0].Addr != addr || stringBalance != "3500000000" || accountNonce != "1" {
//			t.Errorf("Account: Got %v Accounts Created: Expected addr: %s balance: %d nonce %s", newDbAccounts, addr, balance, accountNonce)
//		}
//	})
//	core.DeletePgDb(core.DbName)
//}
//
//func TestInsertTx(t *testing.T) {
//	// Set up a  test transaction
//	key, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
//	signer := types.NewEIP155Signer(big.NewInt(2147483647))
//
//	//Nonce, To Address,Value, GasLimit, Gasprice, data
//	tx1 := types.NewTransaction(1, common.BytesToAddress([]byte{0x11}), big.NewInt(20), 1111, big.NewInt(11111), []byte{0x11, 0x11, 0x11})
//	tx, _ := types.SignTx(tx1, signer, key)
//
//	blockHash := "0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056"
//	txData := stypes.ShyftTxEntryPretty{
//		TxHash:      tx.Hash().Hex(),
//		From:        tx.From().Hex(),
//		To:          tx.To().String(),
//		BlockHash:   blockHash,
//		BlockNumber: strconv.Itoa(21234),
//		Amount:      tx.Value().String(),
//		Cost:        tx.Cost().String(),
//		GasPrice:    tx.GasPrice().Uint64(),
//		GasLimit:    uint64(18000),
//		Gas:         tx.Gas(),
//		Nonce:       tx.Nonce(),
//		Age:         time.Now(),
//		Data:        tx.Data(),
//		Status:      "SUCCESS",
//		IsContract:  false,
//	}
//	t.Run("InsertTx - No Account exists inserts a transaction to the database and updates/creates accounts accordingly", func(t *testing.T) {
//		core.DeletePgDb(core.DbName)
//		db, _ := core.InitDB()
//		core.InsertTx(txData)
//		dbTransactions := []shyftschema.PgTransaction{}
//		err := db.Select(&dbTransactions, "SELECT * FROM txs")
//		if err != nil {
//			panic(err)
//		}
//		pgdb := dbTransactions[0]
//		txInput := txData
//		if len(dbTransactions) != 1 {
//			t.Errorf("Got %v db transactions created \nExpected 1", len(dbTransactions))
//		}
//		if pgdb.TxHash != txInput.TxHash && pgdb.Blockhash != txData.BlockHash && pgdb.To != txData.To && pgdb.From != txData.From && pgdb.Amount != txData.Amount {
//			t.Errorf("Got %+v \nExpected %+v", dbTransactions[0], txData)
//		}
//		newDbAccounts := []shyftschema.Account{}
//		err = db.Select(&newDbAccounts, "SELECT * FROM accounts")
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println("ACCOUNTS",newDbAccounts)
//		if len(newDbAccounts) != 2 {
//			t.Errorf("Got %v db transactions created -  Expected 2", len(newDbAccounts))
//		}
//		toAcct := newDbAccounts[0]
//		fromAcct := newDbAccounts[1]
//		if toAcct.Addr != txData.To && toAcct.Balance != tx.Value().String() && toAcct.Nonce != 1 {
//			t.Errorf("Got %+v \nExpected %s %s %d", toAcct, txData.To, txData.Amount, 1)
//		}
//		fromAcctBal, _ := strconv.Atoi(txData.Amount)
//		fromBalInt := -1 * fromAcctBal
//		product := new(big.Int)
//		product.Mul(new(big.Int).SetInt64(-1), tx.Value())
//		if fromAcct.Addr != txData.To && fromAcct.Balance != product.String() &&
//			fromAcct.Nonce != 1 {
//			t.Errorf("Got %+v \nExpected %s %d %d", fromAcct, txData.From, fromBalInt, 1)
//		}
//		newDbAccountBlocks := []shyftschema.AccountBlock{}
//		err = db.Select(&newDbAccountBlocks, "SELECT * FROM accountblocks")
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println("BLOCKS",newDbAccountBlocks)
//		if len(newDbAccountBlocks) != 2 {
//			t.Errorf("Got %d db accountblocks created -  Expected 2", len(newDbAccountBlocks))
//		}
//		toAcctBl := newDbAccountBlocks[0]
//		fromAcctBl := newDbAccountBlocks[1]
//		if toAcctBl.Acct != txData.To && toAcctBl.Blockhash != txData.BlockHash &&
//			strconv.Itoa(int(toAcctBl.Delta)) != txData.Amount && int(toAcctBl.TxCount) != 1 {
//			t.Errorf("Got %+v \nExpected %s %s %s", toAcctBl, txData.To, txData.BlockHash, txData.Amount)
//		}
//		if fromAcctBl.Acct != txData.To && fromAcctBl.Blockhash != txData.BlockHash &&
//			strconv.Itoa(int(fromAcctBl.Delta)*-1) != txData.Amount && int(toAcctBl.TxCount) != 1 {
//			t.Errorf("Got %+v \nExpected %s %s %s", toAcctBl, txData.To, txData.BlockHash, txData.Amount)
//		}
//	})
//	core.DeletePgDb(core.DbName)
//}
//
//func TestGenesisBlockCreationDeveloper(t *testing.T) {
//	db, _ := core.InitDB()
//	edb, _ := eth.NewShyftTestLDB()
//	shyftTracer := new(eth.ShyftTracer)
//	core.SetIShyftTracer(shyftTracer)
//
//	ethConf := &eth.Config{
//		Genesis:   core.DeveloperGenesisBlock(15, common.Address{}),
//		Etherbase: common.HexToAddress(testAddress),
//		Ethash: ethash.Config{
//			PowMode: ethash.ModeTest,
//		},
//	}
//
//	eth.SetGlobalConfig(ethConf)
//
//	t.Run("SetupGenesisBlock - populates the pg accounts, transactions, and accountblocks appropriately", func(t *testing.T) {
//		core.SetupGenesisBlock(edb, ethConf.Genesis)
//		newDbAccounts := []shyftschema.Account{}
//		err := db.Select(&newDbAccounts, "SELECT * FROM accounts")
//		if err != nil {
//			panic(err)
//		}
//		if len(newDbAccounts) != 9 {
//			t.Errorf("Got %v db transactions created -  Expected 9", len(newDbAccounts))
//		}
//		accountAddresses := []string{}
//		sqlStmnt := "SELECT addr FROM accounts WHERE addr = ANY($1)"
//		err = db.Select(&accountAddresses, sqlStmnt, pq.StringArray(GenesisAcctAddresses))
//		if err != nil {
//			panic(err)
//		}
//		if len(accountAddresses) != 9 {
//			t.Errorf("Got the following acct addresses %+v \n Expected %+v \n", accountAddresses, GenesisAcctAddresses)
//		}
//		var bal string
//		sqlStmnt = "SELECT balance FROM accounts WHERE addr = $1"
//		for _, addr := range GenesisAcctAddresses {
//			err = db.Get(&bal, sqlStmnt, addr)
//			if err != nil {
//				panic(err)
//			}
//			genesisBal := "115792089237316195423570985008687907853269984665640564039457584007913129639927"
//			if addr == "0x0000000000000000000000000000000000000000" {
//				if bal != genesisBal {
//					t.Errorf("Got for Genesis Account Balance %+v \n Expected %s", bal, genesisBal)
//				}
//			} else {
//				if bal != "1" {
//					t.Errorf("Got balance for acct %s: %+v \n Expected %s", addr, bal, genesisBal)
//				}
//			}
//		}
//		for _, acct := range newDbAccounts {
//			if acct.Nonce != 1 {
//				t.Errorf("For acct: %s - got Nonce of %d \n Expected %d", acct.Addr, acct.Nonce, 1)
//			}
//		}
//		dbTransactions := []shyftschema.PgTransaction{}
//		err = db.Select(&dbTransactions, "SELECT * FROM txs")
//		if err != nil {
//			panic(err)
//		}
//		if len(dbTransactions) != 9 {
//			t.Errorf("Got %v db transactions created \nExpected 9", len(dbTransactions))
//		}
//		// genesisFaucetBal := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(9))
//		for _, acct := range dbTransactions {
//			if acct.To != "0x0000000000000000000000000000000000000000" {
//				if acct.From != "genesis" || acct.Blocknumber != "0" || acct.Amount != "1" ||
//					!strings.Contains(acct.TxHash, "genesis") {
//					t.Errorf("Got %+v \n Expected DeveloperGenesisBlock", acct)
//				}
//			}
//		}
//	})
//	core.DeletePgDb(core.DbName)
//}
//
//var (
//	BlockAccounts map[string][]shyftschema.Account
//	BlockHashes   []string
//)
//
//func insertBlocksTransactions() (map[string][]shyftschema.Account, []string, *sqlx.DB) {
//	db, _ := core.InitDB()
//	key, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
//	signer := types.NewEIP155Signer(big.NewInt(2147483647))
//
//	//Nonce, To Address,Value, GasLimit, Gasprice, data
//	tx1 := types.NewTransaction(1, common.BytesToAddress([]byte{0x11}), big.NewInt(5), 1111, big.NewInt(11111), []byte{0x11, 0x11, 0x11})
//	mytx1, _ := types.SignTx(tx1, signer, key)
//	tx2 := types.NewTransaction(2, common.BytesToAddress([]byte{0x22}), big.NewInt(5), 2222, big.NewInt(22222), []byte{0x22, 0x22, 0x22})
//	mytx2, _ := types.SignTx(tx2, signer, key)
//	tx3 := types.NewTransaction(3, common.BytesToAddress([]byte{0x33}), big.NewInt(5), 3333, big.NewInt(33333), []byte{0x33, 0x33, 0x33})
//	mytx3, _ := types.SignTx(tx3, signer, key)
//	txs := []*types.Transaction{mytx1, mytx2}
//	txs1 := []*types.Transaction{mytx3}
//
//	receipt := &types.Receipt{
//		Status:            types.ReceiptStatusSuccessful,
//		CumulativeGasUsed: 1,
//		Logs: []*types.Log{
//			{Address: common.BytesToAddress([]byte{0x11})},
//			{Address: common.BytesToAddress([]byte{0x01, 0x11})},
//		},
//		TxHash:          common.BytesToHash([]byte{0x11, 0x11}),
//		ContractAddress: common.BytesToAddress([]byte{0x01, 0x11, 0x11}),
//		GasUsed:         111111,
//	}
//	receipts := []*types.Receipt{receipt}
//
//	block1 := types.NewBlock(&types.Header{Number: big.NewInt(323)}, txs, nil, receipts)
//	block2 := types.NewBlock(&types.Header{Number: big.NewInt(320)}, txs1, nil, receipts)
//	block3 := types.NewBlock(&types.Header{Number: big.NewInt(322)}, txs, nil, receipts)
//	blocks := []*types.Block{block1, block2, block3}
//	blockHashes := []string{}
//	blockAccounts := map[string][]shyftschema.Account{}
//	core.TruncateTables()
//	for _, bl := range blocks {
//		// Write and verify the block in the database
//		err := core.SWriteBlock(bl, receipts)
//		if err != nil {
//			panic(err)
//		}
//		newDbAccounts := []shyftschema.Account{}
//		err = db.Select(&newDbAccounts, "SELECT * FROM accounts")
//		if err != nil {
//			panic(err)
//		}
//		blockHashes = append(blockHashes, bl.Hash().Hex())
//
//		blockAccounts[bl.Transactions()[0].To().Hex()] = newDbAccounts
//		blockAccounts[bl.Transactions()[0].From().Hex()] = newDbAccounts
//	}
//	return blockAccounts, blockHashes, db
//}
//func TestRollbackReconcilesAccounts(t *testing.T) {
//	t.Run("PgRollback - of all blocks reverses all account balances", func(t *testing.T) {
//		_, blockHashes, db := insertBlocksTransactions()
//
//		// Rollback 1 blocks
//		core.RollbackPgDb(blockHashes[0:])
//		rollBackDbAccounts := []shyftschema.Account{}
//		err := db.Select(&rollBackDbAccounts, "SELECT * FROM accounts")
//		if err != nil {
//			panic(err)
//		}
//		if len(rollBackDbAccounts) > 0 {
//			t.Errorf("Rollback of the following blocks %+v expected %d accounts have %d\n", blockHashes[0:], 0, len(rollBackDbAccounts))
//		}
//		newDbAccountBlocks := []shyftschema.AccountBlock{}
//		err = db.Select(&newDbAccountBlocks, "SELECT * FROM accountblocks")
//		if err != nil {
//			panic(err)
//		}
//		if len(newDbAccountBlocks) != 0 {
//			t.Errorf("Got %d db accountblocks on rollback -  Expected 2", len(newDbAccountBlocks))
//		}
//		rollbackBlocks := []shyftschema.Block{}
//		err = db.Select(&rollbackBlocks, "SELECT * FROM blocks")
//		if err != nil {
//			panic(err)
//		}
//		if len(rollbackBlocks) != 0 {
//			t.Errorf("Got %d db blocks on rollback -  Expected 0", len(rollbackBlocks))
//		}
//		rollbackTxs := []shyftschema.PgTransaction{}
//		err = db.Select(&rollbackTxs, "SELECT * FROM txs")
//		if err != nil {
//			panic(err)
//		}
//		if len(rollbackTxs) != 0 {
//			t.Errorf("Got %d db transactions on rollback -  Expected 0", len(rollbackTxs))
//		}
//		core.DeletePgDb(core.DbName)
//	})
//	t.Run("PgRollback - 2 Blocks- reverses all account balances accordingly", func(t *testing.T) {
//		blockAccounts, blockHashes, db := insertBlocksTransactions()
//		fmt.Println("Rollback by 2 blocks should yield the following balances:")
//		fmt.Println("*********************************************************************")
//		fmt.Printf("\n@block insertion %s \n", blockHashes[0])
//		fmt.Println("*********************************************************************")
//		for _, acct := range blockAccounts {
//			fmt.Printf("%+v \n", acct)
//		}
//		fmt.Println("*********************************************************************")
//		// Rollback 2 blocks
//		core.RollbackPgDb(blockHashes[1:])
//		rollBackDbAccounts := []shyftschema.Account{}
//		err := db.Select(&rollBackDbAccounts, "SELECT * FROM accounts")
//		if err != nil {
//			panic(err)
//		}
//		if len(rollBackDbAccounts) != 5 {
//			t.Errorf("Rollback of the following blocks %+v expected %d accounts have %d\n", blockHashes[1:], 5, len(rollBackDbAccounts))
//		}
//		newDbAccountBlocks := []shyftschema.AccountBlock{}
//		err = db.Select(&newDbAccountBlocks, "SELECT * FROM accountblocks")
//		if err != nil {
//			panic(err)
//		}
//		if len(newDbAccountBlocks) != 5 {
//			t.Errorf("Got %d db accountblocks on rollback -  Expected 5", len(newDbAccountBlocks))
//		}
//		rollbackBlocks := []shyftschema.Block{}
//		err = db.Select(&rollbackBlocks, "SELECT * FROM blocks")
//		if err != nil {
//			panic(err)
//		}
//		if len(rollbackBlocks) != 1 {
//			t.Errorf("Got %d db blocks on rollback -  Expected 1", len(rollbackBlocks))
//		}
//		rollbackTxs := []shyftschema.PgTransaction{}
//		err = db.Select(&rollbackTxs, "SELECT * FROM txs")
//		if err != nil {
//			panic(err)
//		}
//		if len(rollbackTxs) != 2 {
//			t.Errorf("Got %d db transactions on rollback -  Expected 2", len(rollbackTxs))
//		}
//		for _, acct := range blockAccounts[blockHashes[0]] {
//			fetchDbBalanceStmnt := `SELECT * FROM accounts WHERE addr = $1`
//			acctCheck := shyftschema.Account{}
//			err = db.Get(&acctCheck, fetchDbBalanceStmnt, acct.Addr)
//			if err != nil {
//				panic(err)
//			}
//			if acctCheck.Balance != acct.Balance || acctCheck.Nonce != acct.Nonce {
//				t.Errorf("Got Balance: %s Nonce: %d Expected Balance: %s Nonce: %d - Addr: %s\n", acctCheck.Balance, acctCheck.Nonce, acct.Balance, acct.Nonce, acct.Addr)
//			}
//		}
//		core.DeletePgDb(core.DbName)
//	})
//	t.Run("PgRollback - 1 Blocks- reverses all account balances accordingly", func(t *testing.T) {
//		blockAccounts, blockHashes, db := insertBlocksTransactions()
//		fmt.Println("Rollback by 2 blocks should yield the following balances:")
//		fmt.Println("*********************************************************************")
//		fmt.Printf("\n@block insertion %s \n", blockHashes[1])
//		fmt.Println("*********************************************************************")
//		for _, acct := range blockAccounts[blockHashes[1]] {
//			fmt.Printf("%+v \n", acct)
//		}
//		fmt.Println("*********************************************************************")
//		// Rollback 2 blocks
//		core.RollbackPgDb(blockHashes[2:])
//		rollBackDbAccounts := []shyftschema.Account{}
//		err := db.Select(&rollBackDbAccounts, "SELECT * FROM accounts")
//		if err != nil {
//			panic(err)
//		}
//		if len(rollBackDbAccounts) != 6 {
//			t.Errorf("Rollback of the following blocks %+v expected %d accounts have %d\n", blockHashes[1:], 5, len(rollBackDbAccounts))
//		}
//		newDbAccountBlocks := []shyftschema.AccountBlock{}
//		err = db.Select(&newDbAccountBlocks, "SELECT * FROM accountblocks")
//		if err != nil {
//			panic(err)
//		}
//		if len(newDbAccountBlocks) != 9 {
//			t.Errorf("Got %d db accountblocks on rollback -  Expected 5", len(newDbAccountBlocks))
//		}
//		rollbackBlocks := []shyftschema.Block{}
//		err = db.Select(&rollbackBlocks, "SELECT * FROM blocks")
//		if err != nil {
//			panic(err)
//		}
//		if len(rollbackBlocks) != 2 {
//			t.Errorf("Got %d db blocks on rollback -  Expected 2", len(rollbackBlocks))
//		}
//		rollbackTxs := []shyftschema.PgTransaction{}
//		err = db.Select(&rollbackTxs, "SELECT * FROM txs")
//		if err != nil {
//			panic(err)
//		}
//		if len(rollbackTxs) != 3 {
//			t.Errorf("Got %d db transactions on rollback -  Expected 3", len(rollbackTxs))
//		}
//		for _, acct := range blockAccounts[blockHashes[1]] {
//			fetchDbBalanceStmnt := `SELECT * FROM accounts WHERE addr = $1`
//			acctCheck := shyftschema.Account{}
//			err = db.Get(&acctCheck, fetchDbBalanceStmnt, acct.Addr)
//			if err != nil {
//				panic(err)
//			}
//			if acctCheck.Balance != acct.Balance || acctCheck.Nonce != acct.Nonce {
//				t.Errorf("Got Balance: %s Nonce: %d Expected Balance: %s Nonce: %d - Addr: %s \n", acctCheck.Balance, acctCheck.Nonce, acct.Balance, acct.Nonce, acct.Addr)
//			}
//		}
//	})
//	core.DeletePgDb(core.DbName)
//}
