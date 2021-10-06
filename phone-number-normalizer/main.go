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
	UnformattedNumber string
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
	db, err := sql.Open("postgres", postgresqlInfo)
	if err != nil {
		panic(err)
	}
	err = resetDB(db, dbname)
	if err != nil {
		panic(err)
	}
	err = createDB(db, dbname)
	if err != nil {
		panic(err)
	}
	db.Close()

	postgresqlInfo = fmt.Sprintf("%s dbname=%s", postgresqlInfo, dbname)
	db, err = sql.Open("postgres", postgresqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = createPhoneNumberTable(db)
	if err != nil {
		panic(err)
	}
	_, err = insertPhoneNumber(db, "0128787989")
	if err != nil {
		panic(err)
	}
	_, err = insertPhoneNumber(db, "012-8787989")
	if err != nil {
		panic(err)
	}
	_, err = insertPhoneNumber(db, "012-8787-989")
	if err != nil {
		panic(err)
	}
	_, err = insertPhoneNumber(db, "+60-12-8787989")
	if err != nil {
		panic(err)
	}
	_, err = insertPhoneNumber(db, "(012)8787989")
	if err != nil {
		panic(err)
	}
	_, err = insertPhoneNumber(db, "016-90908989")
	if err != nil {
		panic(err)
	}
	_, err = insertPhoneNumber(db, "+6016-9090-8989")
	if err != nil {
		panic(err)
	}
	_, err = insertPhoneNumber(db, "01790908989")
	if err != nil {
		panic(err)
	}
	_, err = insertPhoneNumber(db, "6019802080208")
	if err != nil {
		panic(err)
	}

	numbers, err := getAllPhoneNumber(db)
	if err != nil {
		panic(err)
	}
	for _, p := range numbers {
		fmt.Println("Printing...")
		num := normalizer(p.UnformattedNumber)
		if num != p.UnformattedNumber {
			existing, err := findPhoneNumber(db, num)
			if err != nil {
				panic(err)
			}
			if existing != nil {
				fmt.Println("Deleting phone number: ", p.UnformattedNumber)
				err := deletePhoneNumber(db, p)
				if err != nil {
					panic(err)
				}
			} else {
				fmt.Println("Updating phone number: ", p.UnformattedNumber)
				p.Number = num
				err := updatePhoneNumber(db, p)
				if err != nil {
					panic(err)
				}
			}
		} else {
			fmt.Println("No changes required!")
			p.Number = num
			err := updatePhoneNumber(db, p)
			if err != nil {
				panic(err)
			}
		}
	}
}

func updatePhoneNumber(db *sql.DB, p phoneNumber) error {
	statement := `UPDATE phone_numbers SET formatted_phone=$2 WHERE id=$1`
	_, err := db.Exec(statement, p.Id, p.Number)
	return err
}

func deletePhoneNumber(db *sql.DB, p phoneNumber) error {
	statement := `DELETE FROM phone_numbers WHERE unformatted_phone=$1`
	_, err := db.Exec(statement, p.UnformattedNumber)
	return err
}

func findPhoneNumber(db *sql.DB, number string) (*phoneNumber, error) {
	var p phoneNumber
	row := db.QueryRow("SELECT * FROM phone_numbers WHERE unformatted_phone=$1", number)
	err := row.Scan(&p.Id, &p.UnformattedNumber, &p.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &p, nil
}

func getPhoneNumber(db *sql.DB, id int) (string, error) {
	var number string
	err := db.QueryRow("SELECT unformatted_phone FROM phone_numbers WHERE id=$1", id).Scan(&number)
	if err != nil {
		return "", err
	}
	return number, nil
}

func getAllPhoneNumber(db *sql.DB) ([]phoneNumber, error) {
	var allPhone []phoneNumber
	rows, err := db.Query("SELECT id, unformatted_phone FROM phone_numbers")
	if err != nil {
		return allPhone, err
	}
	defer rows.Close()
	var phone phoneNumber
	for rows.Next() {
		if err := rows.Scan(&phone.Id, &phone.UnformattedNumber); err != nil {
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
	statement := `INSERT INTO phone_numbers(unformatted_phone, formatted_phone) VALUES($1, '') RETURNING id`
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
		unformatted_phone VARCHAR(255),
	    formatted_phone VARCHAR(11)
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