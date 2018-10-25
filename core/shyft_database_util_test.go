package core

import (
	"testing"

	"github.com/ShyftNetwork/go-empyrean/core/types"
	"github.com/ShyftNetwork/go-empyrean/core/sTypes"
	"encoding/json"
	"strings"
	"github.com/ShyftNetwork/go-empyrean/ethdb"
	"fmt"
)

func TestBlock(t *testing.T) {
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
