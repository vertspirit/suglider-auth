package database

import (
	"log"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"suglider-auth/configs"
)

var DataBase *sqlx.DB

func init() {

	var dbErr error

	var DatabaseConfig = configs.ApplicationConfig.Database

	// DB connection setting, no specify any Database
	dbConfig := mysql.Config{
		User:      DatabaseConfig.User,
		Passwd:    DatabaseConfig.Password,
		Net:       "tcp",
		Addr:      DatabaseConfig.Host + ":" + DatabaseConfig.Port,
		AllowNativePasswords: true,
	}

	// Generate MariaDB connection string
	DataBaseURL := dbConfig.FormatDSN()

	// Connect to MariaDB
	DataBase, dbErr = sqlx.Connect("mysql", DataBaseURL)
	if dbErr != nil {
		log.Println("Can not connect to database:", dbErr)
		panic(dbErr)
	}

	// Set DB max connection
    DataBase.SetMaxOpenConns(100)

	log.Println("Connected to DataBase successfully！")

}

// Close DB connection
func Close() {
	DataBase.Close()
	return
}
