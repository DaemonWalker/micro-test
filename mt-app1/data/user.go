package data

import "mt-app1/entity"

func GetUserByEMail(email string) (entity.User, error) {
	sql := "select id,name,email from user t where t.email=?"
	var u entity.User
	err := db.QueryRow(sql, email).Scan(&u.Id, &u.Name, &u.EMail)
	return u, err
}
