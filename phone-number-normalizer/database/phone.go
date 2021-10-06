package phoneDB

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
)


type Phone struct {
	ID                int
	UnformulatedPhone string
	PhoneNumber       string
}

type DB struct {
	db *sql.DB
}

func Open(driverName, dataSource string) (*DB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}


func (db *DB) AllPhoneNumber() ([]Phone, error) {
	var allPhone []Phone
	rows, err := db.db.Query("SELECT id, unformulatedPhone FROM phone_numbers")
	if err != nil {
		return allPhone, err
	}
	defer rows.Close()
	var phone Phone
	for rows.Next() {
		if err := rows.Scan(&phone.ID, &phone.UnformulatedPhone); err != nil {
			log.Fatal(err)
		}
		allPhone = append(allPhone, phone)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return allPhone, nil
}

func (db *DB) UpdatePhoneNumber(p Phone) error {
	statement := `UPDATE phone_numbers SET phoneNumber=$2 WHERE id=$1`
	_, err := db.db.Exec(statement, p.ID, p.PhoneNumber)
	return err
}

func (db *DB) DeletePhoneNumber(p Phone) error {
	statement := `DELETE FROM phone_numbers WHERE unformulatedPhone=$1`
	_, err := db.db.Exec(statement, p.UnformulatedPhone)
	return err
}

func (db *DB) FindPhoneNumber(number string) (*Phone, error) {
	var p Phone
	row := db.db.QueryRow("SELECT * FROM phone_numbers WHERE unformulatedPhone=$1", number)
	err := row.Scan(&p.ID, &p.UnformulatedPhone, &p.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &p, nil
}


func (db *DB) Seed(path string) error {
	phoneSlices, err := readData(path)
	if err != nil {
		panic(err)
	}
	for _, number := range phoneSlices {
		_, err := insertPhoneNumber(db.db, number)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func readData(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Failed to open file: ", err)
		return nil, err
	}
	var nums []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		nums = append(nums, scanner.Text())
	}
	return nums, nil
}

func insertPhoneNumber(db *sql.DB, phone string) (int, error) {
	statement := `INSERT INTO phone_numbers(unformulatedPhone, phoneNumber) VALUES($1, '') RETURNING id`
	var id int
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func Migrate(driverName, dataSource string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		panic(err)
	}

	err = createPhoneNumberTable(db)
	if err != nil {
		panic(err)
	}
	return db.Close()
}

func createPhoneNumberTable(db *sql.DB) error {
	statement := `
	CREATE TABLE IF NOT EXISTS phone_numbers (
		id SERIAL,
		unformulatedPhone VARCHAR(255),
	    phoneNumber VARCHAR(11)
	)`
	_, err := db.Exec(statement)
	return err
}

func Reset(driverName, dataSource, dbname string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		panic(err)
	}
	err = resetDB(db, dbname)
	if err != nil {
		panic(err)
	}
	return db.Close()
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
	return createDB(db, name)
}

//func getPhoneNumber(db *sql.DB, id int) (string, error) {
//	var number string
//	err := db.QueryRow("SELECT unformatted_phone FROM phone_numbers WHERE id=$1", id).Scan(&number)
//	if err != nil {
//		return "", err
//	}
//	return number, nil
//}