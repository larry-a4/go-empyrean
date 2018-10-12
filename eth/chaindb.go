package eth

import (
	"io/ioutil"
	"os"

	"github.com/ShyftNetwork/go-empyrean/ethdb"
	"github.com/ShyftNetwork/go-empyrean/node"
)

var Chaindb_global ethdb.Database

func SetChainDB(db ethdb.Database) {
	Chaindb_global = db
}

func chaindb(ctx *node.ServiceContext, config *Config) (ethdb.Database, error) {
	if Chaindb_global != nil {
		return Chaindb_global, nil
	}

	chainDb, err := CreateDB(ctx, config, "chaindata")
	if err == nil {
		SetChainDB(chainDb)
		return Chaindb_global, nil
	}
	return nil, err
}

// NewShyftTestLDB - returns a (*ethdb.LDBDatabase, func())
func NewShyftTestLDB() (*ethdb.LDBDatabase, error) {
	if Chaindb_global != nil {
		Chaindb_global.Close()
		os.RemoveAll("shyftdb_test_")
	}
	dirname, err := ioutil.TempDir(os.TempDir(), "shyftdb_test_")
	if err != nil {
		panic("failed to create test file: " + err.Error())
	}
	db, err := ethdb.NewLDBDatabase(dirname, 0, 0)
	if err != nil {
		panic("failed to create test database: " + err.Error())
	}
	Chaindb_global = db
	return db, err
}
