package core

import (
	"github.com/ShyftNetwork/go-empyrean/common"
	"github.com/ShyftNetwork/go-empyrean/crypto"
	"math/big"
	"github.com/ShyftNetwork/go-empyrean/core/types"
)

func CreateTestTransactions() []*types.Transaction {
	key, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	signer := types.NewEIP155Signer(big.NewInt(2147483647))

	//Nonce, To Address,Value, GasLimit, Gasprice, data
	tx1 := types.NewTransaction(1, common.BytesToAddress([]byte{0x11}), big.NewInt(5), 1111, big.NewInt(11111), []byte{0x11, 0x11, 0x11})
	mytx1, _ := types.SignTx(tx1, signer, key)
	tx2 := types.NewTransaction(2, common.BytesToAddress([]byte{0x22}), big.NewInt(5), 2222, big.NewInt(22222), []byte{0x22, 0x22, 0x22})
	mytx2, _ := types.SignTx(tx2, signer, key)
	tx3 := types.NewTransaction(3, common.BytesToAddress([]byte{0x33}), big.NewInt(5), 3333, big.NewInt(33333), []byte{0x33, 0x33, 0x33})
	mytx3, _ := types.SignTx(tx3, signer, key)
	txs := []*types.Transaction{mytx1, mytx2, mytx3}

	return txs
}

func CreateTestContractTransactions() []*types.Transaction {
	key, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	signer := types.NewEIP155Signer(big.NewInt(2147483647))

	//Nonce,Value, GasLimit, Gasprice, data
	contractCreation := types.NewContractCreation(1, big.NewInt(111), 1111, big.NewInt(11111), []byte{0x11, 0x11, 0x11})
	mytx4, _ := types.SignTx(contractCreation, signer, key)
	txs2 := []*types.Transaction{mytx4}

	return txs2
}

func CreateTestReceipts() []*types.Receipt {
	receipt := &types.Receipt{
		Status:            types.ReceiptStatusSuccessful,
		CumulativeGasUsed: 1,
		Logs: []*types.Log{
			{Address: common.BytesToAddress([]byte{0x11})},
			{Address: common.BytesToAddress([]byte{0x01, 0x11})},
		},
		TxHash:          common.BytesToHash([]byte{0x11, 0x11}),
		ContractAddress: common.BytesToAddress([]byte{0x01, 0x11, 0x11}),
		GasUsed:         111111,
	}
	receipts := []*types.Receipt{receipt}

	return receipts
}

func CreateTestBlocks() []*types.Block {
	block1 := types.NewBlock(&types.Header{Number: big.NewInt(323)}, CreateTestTransactions(), nil, CreateTestReceipts())
	block2 := types.NewBlock(&types.Header{Number: big.NewInt(322)}, CreateTestTransactions(), nil, CreateTestReceipts())
	block3 := types.NewBlock(&types.Header{Number: big.NewInt(321)}, CreateTestContractTransactions(), nil, CreateTestReceipts())
	blocks := []*types.Block{block1, block2, block3}

	return blocks
}
