package main

import (
	"bytes"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/zhenyang0405/gophercises/phone-number-normalizer/database"
)


const (
	host = "localhost"
	port = 5432
	user = "postgres"
	password = "04051125"
	dbname = "gophercises_phone"
)

func main() {
	postgresqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	err := phoneDB.Reset("postgres", postgresqlInfo, dbname)
	if err != nil {
		panic(err)
	}

	postgresqlInfo = fmt.Sprintf("%s dbname=%s", postgresqlInfo, dbname)
	err = phoneDB.Migrate("postgres", postgresqlInfo)
	if err != nil {
		panic(err)
	}

	db, err := phoneDB.Open("postgres", postgresqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Seed("contact-numbers.txt")
	if err != nil {
		panic(err)
	}

	phones, err := db.AllPhoneNumber()
	if err != nil {
		panic(err)
	}
	for _, p := range phones{
		fmt.Println("Printing...")
		num := normalizer(p.UnformulatedPhone)
		if num != p.UnformulatedPhone {
			existing, err := db.FindPhoneNumber(num)
			if err != nil {
				panic(err)
			}
			if existing != nil {
				fmt.Println("Deleting phone number: ", p.UnformulatedPhone)
				err := db.DeletePhoneNumber(p)
				if err != nil {
					panic(err)
				}
			} else {
				fmt.Println("Updating phone number: ", p.UnformulatedPhone)
				p.PhoneNumber = num
				err := db.UpdatePhoneNumber(p)
				if err != nil {
					panic(err)
				}
			}
		} else {
			fmt.Println("No changes required!")
			p.PhoneNumber = num
			err := db.UpdatePhoneNumber(p)
			if err != nil {
				panic(err)
			}
		}
	}
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