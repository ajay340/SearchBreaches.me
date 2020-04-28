package database

import "time"

type Session struct {
	session_key string
	username    string
	expire      int64
}

func DeleteSessionKey(key string) {
	db := InitDBConn()
	stmt, err := db.Prepare("DELETE FROM sessions WHERE session_key = $1")
	checkErr(err)

	_, err = stmt.Exec(key)
	checkErr(err)
	stmt.Close()
	db.Close()
}

func AddSessionKey(key string, username string) {
	db := InitDBConn()
	sql_addkey := `
	INSERT INTO sessions(
		session_key,
		username,
		expire
	) values($1, $2, $3)
	`
	stmt, err := db.Prepare(sql_addkey)
	checkErr(err)

	_, err = stmt.Exec(key, username, time.Now().Unix()+259200)
	checkErr(err)
	stmt.Close()
	db.Close()
}

func GetSessionUsername(session_key string) string {
	db := InitDBConn()
	rows, err := db.Query("SELECT * FROM Sessions WHERE session_key = $1", session_key)
	checkErr(err)

	session := Session{}
	for rows.Next() {
		err := rows.Scan(&session.session_key, &session.username, &session.expire)
		checkErr(err)
	}
	checkErr(err)
	rows.Close()
	db.Close()
	return session.username
}

func IsValidSession(key string) bool {
	db := InitDBConn()
	rows, err := db.Query("SELECT * FROM Sessions WHERE session_key = $1", key)
	checkErr(err)

	session := Session{}
	for rows.Next() {
		err := rows.Scan(&session.session_key, &session.username, &session.expire)
		checkErr(err)
	}
	checkErr(err)
	defer rows.Close()
	db.Close()
	if time.Now().Unix() < session.expire {
		return true
	} else {
		DeleteSessionKey(key)
		return false
	}
}

func DeleteExpireSessions() {
	db := InitDBConn()
	rows, err := db.Query("SELECT * FROM Sessions")
	checkErr(err)
	var expiredSessions []Session
	for rows.Next() {
		session := Session{}
		err := rows.Scan(&session.session_key, &session.username, &session.expire)
		checkErr(err)
		if time.Now().Unix() >= session.expire {
			expiredSessions = append(expiredSessions, session)
		}
	}
	defer rows.Close()
	db.Close()
	for _, sess := range expiredSessions {
		DeleteSessionKey(sess.session_key)
	}
}
