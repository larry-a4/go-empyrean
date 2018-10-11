package shyfttest

import (
	"github.com/ShyftNetwork/go-empyrean/core"
)

//@SHYFT NOTE: Side effects from PG database therefore need to reset before running

// PgTestDbSetup - reinitializes the pg database
func PgTestDbSetup() {
	core.DeletePgDb(core.DbName())
	_, err := core.DBConnection()
	if err != nil {
		println(err.Error())
		return
	}
}

// PgTestTearDown - resets the pg test database
func PgTestTearDown() {
	core.DeletePgDb(core.DbName())
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
