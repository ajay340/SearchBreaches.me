package database

import (
	"fmt"
	"strconv"
	"strings"
)

type Breach struct {
	ID, Name_of_Covered_Entity, State, Business_Associate_Involved, Individuals_Affected, Date_of_Breach, Type_of_Breach, Location_of_Breached_Information, Date_Posted_or_Updated, Summary, Breach_start, Breach_end, Year, Industry string
}

func FindRowBreach(rowValue string, column string) *Breach {
	db := InitDBConn()
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM breaches WHERE %s = $1", column), strings.ToUpper(rowValue))
	checkErr(err)
	breach := Breach{}
	for rows.Next() {
		err = rows.Scan(&breach.ID, &breach.Name_of_Covered_Entity, &breach.State, &breach.Business_Associate_Involved, &breach.Individuals_Affected, &breach.Date_of_Breach, &breach.Type_of_Breach, &breach.Location_of_Breached_Information, &breach.Date_Posted_or_Updated, &breach.Summary, &breach.Breach_start, &breach.Breach_end, &breach.Year, &breach.Industry)
		checkErr(err)
	}
	rows.Close()
	db.Close()
	return &breach
}

func FindRowsBreach(rowValue string, column string) []Breach {
	db := InitDBConn()
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM breaches WHERE %s like $1", column), "%"+strings.ToUpper(rowValue)+"%")
	checkErr(err)
	var foundBreaches []Breach
	for rows.Next() {
		breach := Breach{}
		err = rows.Scan(&breach.ID, &breach.Name_of_Covered_Entity, &breach.State, &breach.Business_Associate_Involved, &breach.Individuals_Affected, &breach.Date_of_Breach, &breach.Type_of_Breach, &breach.Location_of_Breached_Information, &breach.Date_Posted_or_Updated, &breach.Summary, &breach.Breach_start, &breach.Breach_end, &breach.Year, &breach.Industry)
		checkErr(err)
		if len(breach.Summary) <= 0 {
			breach.Summary = "Breach of " + breach.Name_of_Covered_Entity + " by " + breach.Type_of_Breach
		}
		foundBreaches = append(foundBreaches, breach)
	}
	rows.Close()
	db.Close()
	return foundBreaches
}

func SearchDB(term string) []Breach {
	db := InitDBConn()
	term = strings.ToUpper(term)
	statement := fmt.Sprintf("select * from breaches where name_of_covered_entity like '%%%s%%' or State like '%%%s%%' or Business_Associate_Involved like '%%%s%%' or Individuals_Affected like '%%%s%%' or Date_of_Breach like '%%%s%%' or Type_of_Breach like '%%%s%%' or Location_of_Breached_Information like '%%%s%%' or Date_Posted_or_Updated like '%%%s%%' or Summary like '%%%s%%' or Industry like '%%%s%%' limit 10", term, term, term, term, term, term, term, term, term, term)
	rows, err := db.Query(statement)
	checkErr(err)
	var searchResults []Breach
	for rows.Next() {
		breach := Breach{}
		err = rows.Scan(&breach.ID, &breach.Name_of_Covered_Entity, &breach.State, &breach.Business_Associate_Involved, &breach.Individuals_Affected, &breach.Date_of_Breach, &breach.Type_of_Breach, &breach.Location_of_Breached_Information, &breach.Date_Posted_or_Updated, &breach.Summary, &breach.Breach_start, &breach.Breach_end, &breach.Year, &breach.Industry)
		checkErr(err)
		if len(breach.Summary) <= 0 {
			breach.Summary = "Breach of " + breach.Name_of_Covered_Entity + " by " + breach.Type_of_Breach
		}
		searchResults = append(searchResults, breach)
	}
	rows.Close()
	db.Close()
	return searchResults
}

func AddBreach(Name_of_Covered_Entity string, State string, Business_Associate_Involved string, Individuals_Affected string, Date_of_Breach string, Type_of_Breach string, Location_of_Breached_Information string, Date_Posted_or_Updated string, Summary string, Breach_start string, Breach_end string, Year string, Industry string) {
	db := InitDBConn()
	sql_adduser := `
	INSERT INTO breaches(
		ID, 
		Name_of_Covered_Entity, 
		State, 
		Business_Associate_Involved, 
		Individuals_Affected, 
		Date_of_Breach, 
		Type_of_Breach, 
		Location_of_Breached_Information, 
		Date_Posted_or_Updated, 
		Summary, 
		Breach_start, 
		Breach_end, 
		Year, 
		Industry
	) values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	stmt, err := db.Prepare(sql_adduser)
	checkErr(err)

	id := GenerateID()
	for FindRowBreach(strconv.Itoa(id), "id").ID != "" {
		id = GenerateID()
	}
	idstr := strconv.Itoa(id)
	_, err = stmt.Exec(
		strings.ToUpper(idstr),
		strings.ToUpper(Name_of_Covered_Entity),
		strings.ToUpper(State),
		strings.ToUpper(Business_Associate_Involved),
		strings.ToUpper(Individuals_Affected),
		strings.ToUpper(Date_of_Breach),
		strings.ToUpper(Type_of_Breach),
		strings.ToUpper(Location_of_Breached_Information),
		strings.ToUpper(Date_Posted_or_Updated),
		strings.ToUpper(Summary),
		strings.ToUpper(Breach_start),
		strings.ToUpper(Breach_end),
		strings.ToUpper(Year),
		strings.ToUpper(Industry))
	checkErr(err)
	stmt.Close()
	db.Close()
}
