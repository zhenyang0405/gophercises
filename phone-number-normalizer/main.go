package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

type phoneNumber struct {
	Id int
	Number string
}

const (
	host = "localhost"
	port = 5432
	user = "postgres"
	password = "04051125"
	dbname = "gophercises_phone"
)

func main() {
	postgresqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	//db, err := sql.Open("postgres", postgresqlInfo)
	//if err != nil {
	//	panic(err)
	//}
	//err = resetDB(db, dbname)
	//if err != nil {
	//	panic(err)
	//}
	//err = createDB(db, dbname)
	//if err != nil {
	//	panic(err)
	//}
	//db.Close()

	postgresqlInfo = fmt.Sprintf("%s dbname=%s", postgresqlInfo, dbname)
	db, err := sql.Open("postgres", postgresqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = createPhoneNumberTable(db)
	if err != nil {
		panic(err)
	}
	//_, err = insertPhoneNumber(db, "0128787989")
	//if err != nil {
	//	panic(err)
	//}
	//_, err = insertPhoneNumber(db, "012-8787989")
	//if err != nil {
	//	panic(err)
	//}
	//_, err = insertPhoneNumber(db, "012-8787-989")
	//if err != nil {
	//	panic(err)
	//}
	//_, err = insertPhoneNumber(db, "+60-12-8787989")
	//if err != nil {
	//	panic(err)
	//}
	//_, err = insertPhoneNumber(db, "(012)8787989")
	//if err != nil {
	//	panic(err)
	//}
	//_, err = insertPhoneNumber(db, "016-90908989")
	//if err != nil {
	//	panic(err)
	//}
	//_, err = insertPhoneNumber(db, "+6016-9090-8989")
	//if err != nil {
	//	panic(err)
	//}
	//_, err = insertPhoneNumber(db, "01790908989")
	//if err != nil {
	//	panic(err)
	//}
	//_, err = insertPhoneNumber(db, "6019802080208")
	//if err != nil {
	//	panic(err)
	//}

	numbers, err := getAllPhoneNumber(db)
	if err != nil {
		panic(err)
	}
	for _, p := range numbers {
		fmt.Printf("%+v\n", p)
	}
}

func getPhoneNumber(db *sql.DB, id int) (string, error) {
	var number string
	err := db.QueryRow("SELECT value FROM phone_numbers WHERE id=$1", id).Scan(&number)
	if err != nil {
		return "", err
	}
	return number, nil
}

func getAllPhoneNumber(db *sql.DB) ([]phoneNumber, error) {
	var allPhone []phoneNumber
	rows, err := db.Query("SELECT id, value FROM phone_numbers")
	if err != nil {
		return allPhone, err
	}
	defer rows.Close()
	var phone phoneNumber
	for rows.Next() {
		if err := rows.Scan(&phone.Id, &phone.Number); err != nil {
			log.Fatal(err)
		}
		allPhone = append(allPhone, phone)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return allPhone, nil
}

func insertPhoneNumber(db *sql.DB, phone string) (int, error) {
	statement := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`
	var id int
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func createPhoneNumberTable(db *sql.DB) error {
	statement := `
	CREATE TABLE IF NOT EXISTS phone_numbers (
		id SERIAL,
		value VARCHAR(255)
	)`
	_, err := db.Exec(statement)
	return err
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	return nil
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return nil
}

func normalizer(phone string) string {
	var buf bytes.Buffer
	if strings.HasPrefix(phone, "+6") {
		phone = strings.TrimPrefix(phone, "+6")
	} else if strings.HasPrefix(phone, "6") {
		phone = strings.TrimPrefix(phone, "6")
	}
	for _, char := range phone {
		if char >= '0' && char <= '9' {
			buf.WriteRune(char)
		}
	}
	if len(buf.String()) > 11 {
		return "false"
	}
	return buf.String()
}