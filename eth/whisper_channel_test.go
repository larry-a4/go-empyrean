package eth

import (
	"fmt"
	"testing"

	"github.com/ShyftNetwork/go-empyrean/common"
	"github.com/ShyftNetwork/go-empyrean/core/types"
	"github.com/ShyftNetwork/go-empyrean/crypto"

	"github.com/ShyftNetwork/go-empyrean/consensus"
	"github.com/ShyftNetwork/go-empyrean/core/vm"
	"github.com/ShyftNetwork/go-empyrean/ethdb"
	"github.com/ShyftNetwork/go-empyrean/params"

	"github.com/ShyftNetwork/go-empyrean/consensus/ethash"
	"github.com/ShyftNetwork/go-empyrean/core"
)

// So we can deterministically seed different blockchains
var (
	canonicalSeed = 1
)

var (
	key1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	addr1   = crypto.PubkeyToAddress(key1.PublicKey)
)

// makeBlockChain creates a deterministic chain of blocks rooted at parent.
func makeBlockChain(parent *types.Block, n int, engine consensus.Engine, db ethdb.Database, shyftdb ethdb.SDatabase, seed int) []*types.Block {
	blocks, _ := core.GenerateChain(params.TestChainConfig, parent, engine, db, shyftdb, n, func(i int, b *core.BlockGen) {
		b.SetCoinbase(common.Address{0: byte(seed), 19: byte(i)})
	})
	return blocks
}

func makeHeaderChain(parent *types.Header, n int, engine consensus.Engine, db ethdb.Database, shyftdb ethdb.SDatabase, seed int) []*types.Header {
	blocks := makeBlockChain(types.NewBlockWithHeader(parent), n, engine, db, shyftdb, seed)
	headers := make([]*types.Header, len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	return headers
}

// newCanonical creates a chain database, and injects a deterministic canonical
// chain. Depending on the full flag, if creates either a full block chain or a
// header only chain.
func newCanonical(engine consensus.Engine, n int, full bool) (ethdb.Database, ethdb.SDatabase, *core.BlockChain, error) {
	var (
		db         = ethdb.NewMemDatabase()
		shyftdb, _ = ethdb.NewShyftDatabase()
		genesis    = new(core.Genesis).MustCommit(db)
	)

	// Initialize a fresh chain with only a genesis block
	blockchain, _ := core.NewBlockChain(db, shyftdb, nil, params.AllEthashProtocolChanges, engine, vm.Config{}, nil)
	// Create and inject the requested chain
	if n == 0 {
		return db, shyftdb, blockchain, nil
	}
	if full {
		// Full block-chain requested

		blocks := makeBlockChain(genesis, n, engine, db, shyftdb, canonicalSeed)
		shyftdb.TruncateTables()
		_, err := blockchain.InsertChain(blocks)
		return db, shyftdb, blockchain, err
	}
	// Header-only chain requested
	headers := makeHeaderChain(genesis.Header(), n, engine, db, shyftdb, canonicalSeed)
	fmt.Println(headers)
	_, err := blockchain.InsertHeaderChain(headers, 1)
	return db, shyftdb, blockchain, err
}

// This tests the backend whisper message listening logic to determine if
// A valid blockhash is placed on the listening channel the rollback occurs to the designated blockhash
// TODO: Consider an integration test to test the messaging end to end
func TestWhisperListener(t *testing.T) {
	_, shyftdb, blockchain, _ := newCanonical(ethash.NewFaker(), 20, true)
	t.Run("Whisper rolls back the blockchain to the designated block", func(t *testing.T) {
		rollbackTo := blockchain.GetHeaderByNumber(10).Hash().Hex()
		rollbackFn(rollbackTo, blockchain, nil, shyftdb, addr1)
		wantedHeader := rollbackTo
		actualHeader := blockchain.CurrentHeader().Hash().Hex()
		if wantedHeader != actualHeader {
			t.Errorf("Rollback test failed wanted: %+s  - got: %+s", wantedHeader, actualHeader)
		}
	})
}
