package src

import (
	"database/sql"
	"fmt"
	"log"
	shared "shared_mods"

	_ "github.com/go-sql-driver/mysql"
)

// Checks connectivity and returns the database pointer
func ConnectDB() (*sql.DB, error) {
	log.Println("Checking Database availability.")
	dsn := fmt.Sprintf("%s/%s", shared.MySQLConn, shared.MySQLDBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("Error while opening a connection: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Mysql isn't running: %s", err)
	} else {
		log.Println("MySQL working propertly.")
	}
	return db, nil
}
