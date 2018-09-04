package shyfttest

import "os/exec"

//@SHYFT NOTE: Side effects from PG database therefore need to reset before running

// PgTestDbSetup - reinitializes the pg database
func PgTestDbSetup() {
	cmdStr := "$GOPATH/src/github.com/ShyftNetwork/go-empyrean/shyftdb/postgres_setup_test/init_test_db.sh"
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	_, err := cmd.Output()
	PgRecreateTables()
	if err != nil {
		println(err.Error())
		return
	}
}

// PgTestTearDown - resets the pg test database
func PgTestTearDown() {
	PgTestDbSetup()
}

// PgRecreateTables - recreates pg database tables
func PgRecreateTables() {
	cmdStr := "$GOPATH/src/github.com/ShyftNetwork/go-empyrean/shyftdb/postgres_setup_test/recreate_tables_test.sh"
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	_, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}
}
