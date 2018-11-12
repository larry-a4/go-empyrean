package ethdb

import "testing"

func TestDbCreationExistence(t *testing.T) {
	db, _ := NewShyftDatabase()
	db.TruncateTables()
	t.Run("Creates the Tables Required from the Migration Schema", func(t *testing.T) {
		tableNameQuery := `select table_name from information_schema.tables where table_schema = 'public' AND table_type = 'BASE TABLE' order by table_name ASC;`
		rows, err := db.db.Query(tableNameQuery)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		var tablenames string
		var table string
		notLast := rows.Next()
		for notLast {
			//... rows.Scan
			err = rows.Scan(&table)
			if err != nil {
				panic(err)
			}
			notLast = rows.Next()
			if notLast {
				tablenames += table + ", "
			} else {
				tablenames += table
			}
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}
		want := "accountblocks, accounts, blocks, internaltxs, txs"
		if tablenames != want {
			t.Errorf("Test Failed as wanted: %s  - got: %s", want, tablenames)
		}
	})
	db.db.Close()
}
