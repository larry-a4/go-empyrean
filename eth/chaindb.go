package eth

import (
	"github.com/ShyftNetwork/go-empyrean/ethdb"
	"github.com/ShyftNetwork/go-empyrean/node"
)

var Chaindb_global ethdb.Database
var Shyftdb_global ethdb.SDatabase

func SetChainDB(db ethdb.Database) {
	Chaindb_global = db
}

func SetShyftChainDB(db ethdb.SDatabase) {
	Shyftdb_global = db
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

func shyftdb(ctx *node.ServiceContext, cfg *Config) (ethdb.SDatabase, error) {
	if cfg.Postgres == false {
		return nil, nil
	} else {
		if Shyftdb_global != nil {
			return Shyftdb_global, nil
		}
		shyftDb, err := CreateShyftDB(ctx)
		if err == nil {
			SetShyftChainDB(shyftDb)
			return Shyftdb_global, nil
		}
		return nil, err
	}
}