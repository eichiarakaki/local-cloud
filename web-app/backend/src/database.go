package src

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	shared "shared_mods"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Checks connectivity and returns the database pointer
func ConnectDB() (*sql.DB, error) {
	log.Println("Checking Database availability.")
	dsn := fmt.Sprintf("%s/%s", shared.MySQLConn, shared.MySQLDBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, errors.New("mysql isn't available after 5 attempts")
	}

	for i := 0; i < 5; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		log.Printf("MySQL isn't running yet, retrying... (%s)\n", err)
		time.Sleep(5 * time.Second)
	}
	log.Println("MySQL working propertly.")

	return db, nil
}
