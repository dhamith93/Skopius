package store

import (
	"database/sql"

	"github.com/dhamith93/Skopius/internal/models"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

type Store struct {
	DB *sql.DB
}

func NewStore(dsn string) (*Store, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS results (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		agent_id TEXT,
		service TEXT,
		url TEXT,
		status TEXT,
		code INTEGER,
		latency INTEGER,
		error TEXT,
		timestamp TIMESTAMP,
		received TIMESTAMP
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	schema = `CREATE TABLE IF NOT EXISTS agents (
		id TEXT PRIMARY KEY,
		hostname TEXT,
		region TEXT,
		token TEXT,
		created_at DATETIME,
		last_seen DATETIME
	);`

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return &Store{DB: db}, nil
}

func (s *Store) SaveResult(res models.CheckResult) error {
	_, err := s.DB.Exec(`
        INSERT INTO results (agent_id, service, url, status, code, latency, error, timestamp, received)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		res.AgentID, res.Service, res.URL, res.Status,
		res.Code, res.Latency, res.Error, res.Timestamp, res.Received)
	return err
}

func (s *Store) GetLatestResults(limit int) ([]models.CheckResult, error) {
	rows, err := s.DB.Query(`
        SELECT id, agent_id, service, url, status, code, latency, error, timestamp, received
        FROM results
        ORDER BY received DESC
        LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.CheckResult
	for rows.Next() {
		var r models.CheckResult
		if err := rows.Scan(&r.ID, &r.AgentID, &r.Service, &r.URL, &r.Status,
			&r.Code, &r.Latency, &r.Error, &r.Timestamp, &r.Received); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

func (s *Store) RegisterAgent(id string, hostname string, region string, token string) error {
	_, err := s.DB.Exec(`
	INSERT INTO agents (id, hostname, region, token, created_at, last_seen)
	VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
	`, id, hostname, region, token)
	return err
}

func (s *Store) UpdateLastSeen(id string) error {
	_, err := s.DB.Exec(`
	UPDATE agents SET last_seen = datetime('now') WHERE id = ?
	`, id)
	return err
}

func (s *Store) GetAgentByID(id string) (string, error) {
	var token string
	err := s.DB.QueryRow(`SELECT token FROM agents WHERE id = ?`, id).Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}
