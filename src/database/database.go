package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/lib/pq"
)

var databasePath string = ""

func SetDBPath(filepath string) {
	databasePath = filepath
	if !checkTable("breaches") {
		statement := `CREATE TABLE breaches(
			ID text, 
			Name_of_Covered_Entity text, 
			State text, 
			Business_Associate_Involved text, 
			Individuals_Affected text, 
			Date_of_Breach text, 
			Type_of_Breach text, 
			Location_of_Breached_Information text, 
			Date_Posted_or_Updated text, 
			Summary text, 
			Breach_start text, 
			Breach_end text, 
			Year text, 
			Industry text
			)`
		createTable(statement)
		fmt.Println("Created table breaches")
	}
	if !checkTable("sessions") {
		statement := `CREATE TABLE sessions(
			session_key text,
			username text,
			expire integer
		)`
		createTable(statement)
		fmt.Println("Created table sessions")
	}
	if !checkTable("users") {
		statement := `CREATE TABLE users(
			Username text,
			Password text,
			InsertedDatetime integer
		)`
		createTable(statement)
		fmt.Println("Created table users")
	}
}

func ReadDBPath() string {
	return databasePath
}

type PostgresConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Endpoint string `json:"endpoint"`
	Port     string `json:"port"`
	Database string `json:"database"`
}

func createTable(tableStatement string) {
	db := InitDBConn()
	db.Exec(tableStatement)
	db.Close()
}

func checkTable(tableName string) bool {
	db := InitDBConn()
	rows, err := db.Query("select tablename from pg_tables where schemaname='public' and tablename=$1", tableName)
	checkErr(err)
	var table string
	for rows.Next() {
		err = rows.Scan(&table)
		checkErr(err)
	}
	rows.Close()
	db.Close()
	if table != "" {
		return true
	} else {
		return false
	}
}

func InitDBConn() *sql.DB {
	jsonFile, err := os.Open(ReadDBPath())
	if err != nil {
		panic(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var PDBConfig PostgresConfig
	json.Unmarshal(byteValue, &PDBConfig)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		PDBConfig.Endpoint, PDBConfig.Port, PDBConfig.Username, PDBConfig.Password, PDBConfig.Database)

	db, err := sql.Open("postgres", psqlInfo)
	checkErr(err)
	return db
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
