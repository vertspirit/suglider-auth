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

func UserSignUp(username, password, comfirm_pwd, mail, address string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	sqlStr := "INSERT INTO suglider.user_info(user_id, username, password, comfirm_pwd, mail, address) VALUES (UNHEX(REPLACE(UUID(), '-', '')),?,?,?,?,?)"
	_, err = DataBase.ExecContext(ctx, sqlStr, username, password, comfirm_pwd, mail, address)
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

	err = DataBase.GetContext(ctx, &userInfo, "SELECT username, password, LOWER(HEX(user_id)) AS user_id FROM suglider.user_info WHERE username=?", username)

	return userInfo, err
}

func TotpStoreSecret(user_id, totp_secret, totp_url string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	sqlStr := "INSERT INTO suglider.totp(user_id, totp_secret, totp_url) VALUES (UNHEX(?),?,?)"
	_, err = DataBase.ExecContext(ctx, sqlStr, user_id, totp_secret, totp_url)
	return err
}

func TotpUpdateSecret(username, totp_secret, totp_url string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	sqlStr := "UPDATE suglider.totp " + 
			"JOIN suglider.user_info ON user_info.user_id = totp.user_id " +
			"SET totp_secret = ?, totp_url = ? " +
			"WHERE user_info.username = ?"
	_, err = DataBase.ExecContext(ctx, sqlStr, totp_secret, totp_url, username)
	return err
}

func TotpUserCheck(user_id string) (rowCount int, err error) {
	var count int
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	sqlStr := "SELECT COUNT(*) FROM suglider.totp WHERE user_id=UNHEX(?)"
	err = DataBase.GetContext(ctx, &count, sqlStr, user_id)
	return count, err
}

func TotpUserData(username string) (totpUserInfo TotpUserInfo, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	sqlStr := "SELECT totp.user_id, totp.totp_enabled, totp_verified, totp.totp_secret, totp_url " + 
			"FROM suglider.user_info " + 
			"INNER JOIN suglider.totp ON user_info.user_id = totp.user_id " +
			"WHERE user_info.username=?"
	err = DataBase.GetContext(ctx, &totpUserInfo, sqlStr, username)
	return totpUserInfo, err
}

func TotpUpdateVerify(username string, totp_enabled, totp_verified bool) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	sqlStr := "UPDATE suglider.totp " + 
			"JOIN suglider.user_info ON user_info.user_id = totp.user_id " +
			"SET totp_enabled = ?, totp_verified = ? " +
			"WHERE username = ?"
	_, err = DataBase.ExecContext(ctx, sqlStr, totp_enabled, totp_verified, username)
	return err
}

func TotpUpdateEnabled(username string, totp_enabled bool) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	sqlStr := "UPDATE suglider.totp " + 
			"JOIN suglider.user_info ON user_info.user_id = totp.user_id " + 
			"SET totp_enabled = ? " +
			"WHERE username = ?"
	_, err = DataBase.ExecContext(ctx, sqlStr, totp_enabled, username)
	return err
}

func LookupUserID(username string) (userIDInfo UserIDInfo ,err error){
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	err = DataBase.GetContext(ctx, &userIDInfo, "SELECT LOWER(HEX(user_id)) AS user_id FROM suglider.user_info WHERE username=?", username)
	return userIDInfo, err
}