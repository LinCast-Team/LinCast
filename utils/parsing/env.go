package parsing

import (
	"log"
	"os"
	"strconv"
)

func ParseEnv() (dbPort int, dbHost, dbUser, dbPassword, dbName string) {
	dbHost = os.Getenv("DB_HOST")
	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_NAME")

	dbPort, err := strconv.Atoi((os.Getenv("DB_PORT")))
	if err != nil {
		log.Fatal(err.Error())
	}

	return
}
