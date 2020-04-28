package database

import (
	"time"
)

type User struct {
	Username, Password, InsertedDatetime string
}

func AddUser(username string, password string) {
	db := InitDBConn()
	sql_adduser := `
	INSERT INTO users(
		Username,
		Password,
		InsertedDatetime
	) values($1, $2, $3)
	`
	stmt, err := db.Prepare(sql_adduser)
	checkErr(err)

	_, err = stmt.Exec(username, password, time.Now().Unix())
	checkErr(err)
	stmt.Close()
	db.Close()
}

func DeleteUser(username string) {
	db := InitDBConn()
	stmt, err := db.Prepare("DELETE FROM users WHERE Username = $1")
	checkErr(err)

	_, err = stmt.Exec(username)
	checkErr(err)
	stmt.Close()
	db.Close()
}

func UpdateUser(username string, newpassword string) {
	db := InitDBConn()
	sql_adduser := `Update users set Password=$1 where Username=$2`
	stmt, err := db.Prepare(sql_adduser)
	checkErr(err)

	_, err = stmt.Exec(newpassword, username)
	checkErr(err)
	stmt.Close()
	db.Close()
}

func ReadUserInfo(username string) User {
	db := InitDBConn()
	rows, err := db.Query("SELECT * FROM Users WHERE Username = $1", username)
	checkErr(err)

	user := User{}
	for rows.Next() {
		err := rows.Scan(&user.Username, &user.Password, &user.InsertedDatetime)
		checkErr(err)
	}
	checkErr(err)
	rows.Close()
	db.Close()
	return user
}

func UserExists(username string) bool {
	user := ReadUserInfo(username)
	if user.Username == "" {
		return false
	} else {
		return true
	}
}
