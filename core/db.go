package core

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ShyftNetwork/go-empyrean/shyft_schema"
	"github.com/jmoiron/sqlx"
)

var blockExplorerDb *sqlx.DB

const (
	defaultTestDb  = "shyftdbtest"
	defaultDb      = "shyftdb"
	connStrTest    = "user=postgres password=docker sslmode=disable"
	connStrDocker  = "user=postgres host=pg password=docker sslmode=disable"
	connStrDefault = "user=postgres host=localhost sslmode=disable"
)

// InitDB - initalizes a Postgresql Database for use by the Blockexplorer
func InitDB() (*sqlx.DB, error) {
	// To set the environment you can run the program with an ENV variable DBENV.
	// DBENV defaults to local for purposes of running the correct local
	// database connection parameters but will use docker connection parameters if DBENV=docker

	// Check for existence of Database
	exist, _ := DbExists(DbName())
	if !exist {
		// create the db
		CreatePgDb(DbName())
	}
	// connect to the designated db & create tables if necessary
	blockExplorerDb = Connect(ShyftConnectStr())
	fmt.Println(shyftschema.MakeTableQuery())
	blockExplorerDb.MustExec(shyftschema.MakeTableQuery())
	return blockExplorerDb, nil
}

// ShyftConnectStr - Returns the Connection String With The appropriate database
func ShyftConnectStr() string {
	return fmt.Sprintf("%s%s%s", ConnectionStr(), " dbname=", DbName())
}

// Connect - return a connection to a postgres database wi
func Connect(connectURL string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", connectURL)
	if err != nil {
		fmt.Println("ERROR OPENING DB, NOT INITIALIZING")
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

// DbName - gets the correct db name based on the environment
func DbName() string {
	if flag.Lookup("test.v") == nil {
		return defaultDb
	} else {
		return defaultTestDb
	}
}

// CreatePgDb - Creates a DB
func CreatePgDb(dbName string) {
	conn := Connect(ConnectionStr())
	sqlCmd := fmt.Sprintf(`CREATE DATABASE %s;`, dbName)
	_, err := conn.Exec(sqlCmd)
	if err != nil {
		panic(err)
	}
	conn.Close()
}

// DeletePgDb - Deletes the designated DB
func DeletePgDb(dbName string) {
	conn := Connect(ConnectionStr())
	q1 := fmt.Sprintf(`SELECT pg_terminate_backend(pid)FROM pg_stat_activity WHERE datname = '%s';`, dbName)
	_, err1 := conn.Exec(q1)
	if err1 != nil || err1 == sql.ErrNoRows {
		panic(err1)
	}
	q2 := fmt.Sprintf(`DROP DATABASE IF EXISTS %s;`, dbName)
	_, err2 := conn.Exec(q2)
	if err2 != nil || err2 == sql.ErrNoRows {
		panic(err2)
	}
	conn.Close()
}

// ConnectionStr - return a Connection to the PG admin database
func ConnectionStr() string {
	dbEnv := os.Getenv("DBENV")
	if flag.Lookup("test.v") == nil {
		switch dbEnv {
		default:
			return connStrDefault
		case "docker":
			return connStrDocker
		}
	} else {
		return connStrTest
	}
}

// DbExists - Checks whether the named database exists returns true or false
func DbExists(dbname string) (bool, error) {
	sqldb := Connect(ConnectionStr())
	var exists bool
	sqlStatement := fmt.Sprintf(`select exists(SELECT datname FROM pg_catalog.pg_database WHERE lower(datname) = '%s');`, strings.ToLower(dbname))
	error := sqldb.QueryRow(sqlStatement).Scan(&exists)
	switch {
	case error == sql.ErrNoRows:
		sqldb.Close()
		return false, error
	case error != nil:
		return false, error
		panic(error)
	default:
		sqldb.Close()
		return exists, error
	}
}

// DBConnection returns a connection to the PG BlockExporer DB
func DBConnection() (*sqlx.DB, error) {
	if blockExplorerDb == nil {
		_, err := InitDB()
		if err != nil {
			return nil, err
		}
	}
	conn := blockExplorerDb
	conn.Ping()
	return conn, nil
}

func ClearTables() {
	sqldb, err := DBConnection()
	if err != nil {
		panic(err)
	}

	sqlStatementTx := `DELETE FROM txs`
	_, err = sqldb.Exec(sqlStatementTx)
	if err != nil {
		panic(err)
	}

	sqlStatementAcc := `DELETE FROM accounts`
	_, err = sqldb.Exec(sqlStatementAcc)
	if err != nil {
		panic(err)
	}

	sqlStatement := `DELETE FROM blocks`
	_, err = sqldb.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
}

// TruncateTables - Is primarily user to clear the pg database between unit tests
func TruncateTables() {
	sqldb, err := DBConnection()
	if err != nil {
		panic(err)
	}
	tx, _ := sqldb.Begin()
	sqlStatement := `TRUNCATE TABLE txs, accounts, blocks, internaltxs RESTART IDENTITY CASCADE;`
	_, err = tx.Exec(sqlStatement)
	tx.Commit()
	if err != nil {
		panic(err)
	}
}
