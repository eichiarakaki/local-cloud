package src

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	shared "shared_mods"

	_ "github.com/go-sql-driver/mysql"
)

// Checks connectivity and returns the database pointer
func connectDB(dsn string) (*sql.DB, error) {
	log.Println("Checking Database availability.")

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

// Check if table exists, if not, creates it.
func checkAndCreateTable(db *sql.DB) error {
	query := fmt.Sprintf("SHOW TABLES LIKE '%s'", shared.MySQLTableName)
	var table string

	err := db.QueryRow(query).Scan(&table)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("Error verifying the table: %s", err)
	}

	if table == "" {
		// Read the SQL file
		currentDir, _ := os.Getwd()
		sqlBytes, err := os.ReadFile(filepath.Join(currentDir, "/sql/createTable.sql"))
		if err != nil {
			return fmt.Errorf("Couldn't read from the file - Verify the path. %s\n", err)
		}
		newTable := string(sqlBytes)

		_, err = db.Exec(newTable)
		if err != nil {
			return fmt.Errorf("Couldn't create the table: %v", err)
		}
		log.Println("Table successfully created")
	}

	return nil
}

func upload(db *sql.DB, vd *VideoData) error {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE filepath = \"%s\"", shared.MySQLTableName, vd.Path)
	var count int

	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return fmt.Errorf("Error executing query: %v", err)
	}

	if count == 0 {
		insertQuery := fmt.Sprintf("INSERT INTO %s (filepath, filename) VALUES(?, ?)", shared.MySQLTableName)
		_, err := db.Exec(insertQuery, vd.Path, vd.Title)
		if err != nil {
			return fmt.Errorf("Error inserting data: %v", err)
		}
	} else {
		log.Printf("%s already exists in the database.", vd.Path)
	}

	return nil
}

func (vd *VideoData) PushToBD() {
	dsn := fmt.Sprintf("%s/%s", shared.MySQLConn, shared.MySQLDBName)
	db, err := connectDB(dsn)
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(1) // Modify when concurrency is set
	db.SetMaxIdleConns(1)

	err = checkAndCreateTable(db)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	err = upload(db, vd)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	fmt.Println("Cycle finished.")
}
