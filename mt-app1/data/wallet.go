package data

import "database/sql"

func UpdateMoney(userId int, change int, exexecutor executor) (sql.Result, error) {
	updateSql := `UPDATE
	wallet
SET
	money = money + ?
WHERE
	user_id = ?`
	return exexecutor.Exec(updateSql, change, userId)
}
