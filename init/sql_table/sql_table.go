package sql_table

import (
	"log/slog"
	_ "github.com/go-sql-driver/mysql"
	mariadb "suglider-auth/internal/database"
	"io/ioutil"
	"os"
	"strings"
)

func SugliderTableInit() {

	env := os.Getenv("ENV")

	if len(env) == 0 {
		panic("ENV variable is empty!")
	}

	filePath := "./configs/sql_syntax/" + env + ".sql"

	sqlBytes := readSQLFile(filePath)

	// Split SQL syntax
	sqlStatements := strings.Split(string(sqlBytes), ";")

	// Excute every SQL syntax
	for _, sqlStatement := range sqlStatements {
		sqlStatement = strings.TrimSpace(sqlStatement)
		if sqlStatement == "" {
			continue
		}

		mariadb.DataBase.Exec(sqlStatement)
	}

	slog.Info("SQL file excution complete!")

}

func readSQLFile(filePath string) ([]byte) {

	// Read SQL file
	sqlFile, err := os.Open(filePath)
	if err != nil {
		slog.Error(err.Error())
		
		panic("Can not open " + filePath)
	}
	defer sqlFile.Close()

	// Read SQL file content
	sqlBytes, err := ioutil.ReadAll(sqlFile)
	if err != nil {
		slog.Error(err.Error())
		panic("Can not read content " + filePath)
	}

	return sqlBytes
}