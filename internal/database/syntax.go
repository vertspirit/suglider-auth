package database

import (
	"database/sql"
	"context"
	"time"
	"fmt"
	"log/slog"
	"suglider-auth/configs"
)

var dbTimeOut time.Duration

func init() {
	var DatabaseConfig = configs.ApplicationConfig.Database
	var err error

	dbTimeOut, err = time.ParseDuration(DatabaseConfig.Timeout)

	if err != nil {
		errorMessage := fmt.Sprintf("DB timeout string convert to duration failed: %v", err)
		slog.Error(errorMessage)

		panic(err)
	}
}

func UserSignUp(username, password, mail, address string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	sqlStr := "INSERT INTO suglider.user_info(user_id, username, password, mail, address) VALUES (UNHEX(REPLACE(UUID(), '-', '')),?,?,?,?)"
	_, err = DataBase.ExecContext(ctx, sqlStr, username, password, mail, address)
	return err
}

func UserDelete(username, mail string) (result sql.Result, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	sqlStr := "DELETE FROM suglider.user_info WHERE username=? AND mail=?"
	result, err = DataBase.ExecContext(ctx, sqlStr, username, mail)
	return result, err
}

func UserDeleteByUUID(user_id, username, mail string) (result sql.Result, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	// UNHEX(?) can convert user_id to binary(16)
	sqlStr := "DELETE FROM suglider.user_info WHERE user_id=UNHEX(?) AND username=? AND mail=?"
	result, err = DataBase.ExecContext(ctx, sqlStr,user_id ,username, mail)
	return result, err
}

func UserLogin(username string) (userInfo UserDBInfo ,err error){
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	err = DataBase.GetContext(ctx, &userInfo, "SELECT username, password FROM suglider.user_info WHERE username=?", username)

	return userInfo, err
}