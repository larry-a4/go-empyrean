package core

import "strconv"

var TestDbInstances []string

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

// AssignTestDbInstanceName - returns a Db Name for testing
func AssignTestDbInstanceName() string {
	// returns name of current test db
	var dbNameAssigned string
	if len(TestDbInstances) == 0 {
		dbNameAssigned = defaultTestDb + "_1"
		TestDbInstances = append(TestDbInstances, dbNameAssigned)
		return dbNameAssigned
	}
	newDbSuffix := len(TestDbInstances) + 1
	dbNameAssigned = defaultTestDb + "_" + strconv.Itoa(newDbSuffix)
	return dbNameAssigned
}