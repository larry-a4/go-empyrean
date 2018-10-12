package shyfttest

import (
	"strconv"

	"github.com/ShyftNetwork/go-empyrean/core"
)

const (
	defaultTestDb = "shyftdbtest"
)

var testDbInstances []string

// AssignTestDbInstanceName - returns a Db Name for testing
func assignTestDbInstanceName() string {
	// returns name of current test db
	var dbNameAssigned string
	if len(testDbInstances) == 0 {
		dbNameAssigned = defaultTestDb + "_1"
		testDbInstances = append(testDbInstances, dbNameAssigned)
		return dbNameAssigned
	}
	newDbSuffix := len(testDbInstances) + 1
	dbNameAssigned = defaultTestDb + "_" + strconv.Itoa(newDbSuffix)
	return dbNameAssigned
}

func sliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

//@SHYFT NOTE: Side effects from PG database therefore need to reset before running

// PgTestDbSetup - reinitializes the pg database and returns the name of the testdb
func PgTestDbSetup() string {
	// Check Db Instances - and get a db name to use
	core.ActiveTestDb = assignTestDbInstanceName()
	// core.DeletePgDb(.DbName()core)
	_, err := core.DBConnection()
	if err != nil {
		println(err.Error())
		return ""
	}
	return core.ActiveTestDb
}

// PgTestTearDown - resets the pg test database
func PgTestTearDown(dbname string) {
	// remove db from list of active dbs
	index := sliceIndex(len(testDbInstances), func(i int) bool { return testDbInstances[i] == dbname })
	testDbInstances = append(testDbInstances[:index], testDbInstances[index+1:]...)
	core.DeletePgDb(dbname)
}

// PgRecreateTables - recreates pg database tables
func PgRecreateTables() {
	core.DeletePgDb(core.DbName())
	_, err := core.DBConnection()
	if err != nil {
		println(err.Error())
		return
	}
}
