package src

import (
	"database/sql"
	"errors"
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
		return nil, fmt.Errorf("error while opening a connection: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("mysql isn't running: %s", err)
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
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error verifying the table: %s", err)
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
			return fmt.Errorf("couldn't create the table: %v", err)
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
		insertQuery := fmt.Sprintf("INSERT INTO %s (filepath, filename, thumbnail) VALUES(?, ?, ?)", shared.MySQLTableName)
		_, err := db.Exec(insertQuery, vd.Path, vd.Title, vd.Thumbnail)
		if err != nil {
			return fmt.Errorf("error inserting data: %v", err)
		}
	} else {
		log.Printf("%s already exists in the database.", vd.Path)
	}

	return nil
}

// This function removes the video if there's an error on the database.
func (vd *VideoData) safeErr(err error) {
	newErr := os.Remove(vd.Path)
	if newErr != nil {
		log.Printf("unable to remove '%s', please remove it manually.\n", newErr)
	}
	log.Printf("'%s' was removed successfully.", vd.Path)
	newErr = os.Remove(vd.Thumbnail)
	if newErr != nil {
		log.Printf("unable to remove '%s', please remove it manually.\n", newErr)
	}
	log.Printf("'%s' was removed successfully.", vd.Thumbnail)

	log.Panicln(err)
}

func (vd *VideoData) PushToBD() {
	dsn := fmt.Sprintf("%s/%s", shared.MySQLConn, shared.MySQLDBName)
	db, err := connectDB(dsn)
	if err != nil {
		vd.safeErr(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("[ERROR] closing database connection: %s", err)
		}
	}(db)

	db.SetMaxOpenConns(1) // Modify when concurrency is set
	db.SetMaxIdleConns(1)

	err = checkAndCreateTable(db)
	if err != nil {
		vd.safeErr(err)
	}

	err = upload(db, vd)
	if err != nil {
		vd.safeErr(err)
	}

	fmt.Println("Cycle finished.")
}
