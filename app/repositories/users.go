package repositories

import (
	"ticket-sync-mvc/framework/database"
	"ticket-sync-mvc/framework/support"
)

func UserByEmail(email string) support.M {
	rows := database.Rows("SELECT * FROM users WHERE email=?", email)
	if len(rows) == 0 {
		return nil
	}
	return rows[0]
}
func CreateUser(name, email, hash string) (int64, error) {
	result, err := database.DB.Exec("INSERT INTO users(name,email,password_hash,created_at) VALUES(?,?,?,?)", name, email, hash, support.Now())
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
