package core

import (
	"regexp"
	"strconv"
)

var TestDbInstances []string

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

func stripNumber(str string) int {
	re := regexp.MustCompile(`\w*_([0-9]+)$`)
	match := re.FindStringSubmatch(str)
	d, err := strconv.Atoi(match[1])
	if err != nil {
		return -1
	}
	return d
}

func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
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

	var dbNumbersUsed []int
	for _, x := range TestDbInstances {
		dbNumbersUsed = append(dbNumbersUsed, stripNumber(x))
	}

	dbNum := false
	dbInt := 1
	for dbNum == false {
		if intInSlice(dbInt, dbNumbersUsed) {
			dbInt++
		} else {
			dbNum = true
		}
	}
	dbNameAssigned = defaultTestDb + "_" + strconv.Itoa(dbInt)
	TestDbInstances = append(TestDbInstances, dbNameAssigned)
	return dbNameAssigned
}
