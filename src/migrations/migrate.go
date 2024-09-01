package migrations

import (
	"fmt"
	"os"
	"path/filepath"
	"spacesona-go-backend/db"
	"strings"
)

func Migrate() {
	fmt.Println("migrating database")
	path := filepath.Join("migrations", "0-init.sql")
	runSqlFile(path)
}

func DummyData() {
	path := filepath.Join("migrations", "dummyData.sql")
	runSqlFile(path)
}

func runSqlFile(path string) {
	dat, ioErr := os.ReadFile(path)
	if ioErr != nil {
		fmt.Println("Error:", ioErr)
		panic(ioErr)
	}
	sqlStatments := strings.Split(string(dat), ";")
	for _, sqlStatement := range sqlStatments {
		sql := strings.TrimSpace(sqlStatement)
		res, dbErr := db.UseSQL().Exec(sql)
		if dbErr != nil {
			fmt.Println("Error:", dbErr)
		} else {
			fmt.Println(res.RowsAffected())
		}
	}
}
