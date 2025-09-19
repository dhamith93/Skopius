package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func getDb() *sql.DB {
	db, err := sql.Open("sqlite3", "./agents.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS agents (
		id TEXT PRIMARY KEY,
		hostname TEXT,
		region TEXT,
		token TEXT,
		created_at DATETIME,
		last_seen DATETIME
	);`)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func RegisterAgent(id string, hostname string, region string, token string) error {
	db := getDb()
	_, err := db.Exec(`
	INSERT INTO agents (id, hostname, region, token, created_at, last_seen)
	VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
	`, id, hostname, region, token)
	return err
}

func UpdateLastSeen(id string) error {
	db := getDb()
	_, err := db.Exec(`
	UPDATE agents SET last_seen = datetime('now') WHERE id = ?
	`, id)
	return err
}

func GetAgentByID(id string) (string, error) {
	var token string
	db := getDb()
	err := db.QueryRow(`SELECT token FROM agents WHERE id = ?`, id).Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}
