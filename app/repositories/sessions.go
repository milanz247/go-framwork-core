package repositories

import (
	"ticket-sync-mvc/framework/database"
	"ticket-sync-mvc/framework/support"
)

func FindSession(token string) support.M {
	rows := database.Rows("SELECT u.id,u.name,u.email,s.csrf FROM users u JOIN go_sessions s ON s.user_id=u.id WHERE s.token=? AND s.expires>?", token, support.Now())
	if len(rows) == 0 {
		return nil
	}
	return rows[0]
}
func SaveSession(token string, uid int64, csrf, expires string) {
	database.Exec("DELETE FROM go_sessions WHERE expires<?", support.Now())
	database.Exec("INSERT INTO go_sessions(token,user_id,csrf,expires) VALUES(?,?,?,?)", token, uid, csrf, expires)
}
func DeleteSession(token string) { database.Exec("DELETE FROM go_sessions WHERE token=?", token) }
