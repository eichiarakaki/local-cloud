package src

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	shared "shared_mods"

	_ "github.com/go-sql-driver/mysql"
)

// Checks connectivity and returns the database pointer
func connectDB(dsn string) (*sql.DB, error) {
	log.Println("Checking Database availability.")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		errorfmt := fmt.Sprintf("Error while opening a connection: %s", err)
		return nil, errors.New(errorfmt)
	}

	err = db.Ping()
	if err != nil {
		errorfmt := fmt.Sprintf("Mysql isn't running: %s", err)
		return nil, errors.New(errorfmt)
	} else {
		log.Println("MySQL working propertly.")
	}
	return db, nil
}

func (vd *VideoData) PushToBD() {
	dsn := fmt.Sprintf("%s@tcp(%s)/%s", shared.MySQLUserPass, shared.MySQLSocket, shared.MySQLDBName)
	db, err := connectDB(dsn)
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(1) // Modify when concurrency is set
	db.SetMaxIdleConns(1)

	fmt.Println("Executed correctly") // Debug
}
