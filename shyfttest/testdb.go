package shyfttest

import (
	"github.com/ShyftNetwork/go-empyrean/core"
)

const (
	defaultTestDb = "shyftdbtest"
)

// PgTestDbSetup - reinitializes the pg database and returns the name of the testdb
func PgTestDbSetup() string {
	// Check Db Instances - and get a db name to use
	core.ActiveTestDb = core.AssignTestDbInstanceName()
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
	index := core.SliceIndex(len(core.TestDbInstances), func(i int) bool { return core.TestDbInstances[i] == dbname })
	core.TestDbInstances = append(core.TestDbInstances[:index], core.TestDbInstances[index+1:]...)
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
