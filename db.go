package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type ServerRow struct {
	ID        int
	IP        string
	QueryPort int
}

func OpenDB(cfg *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func GetServers(db *sql.DB, table string) ([]ServerRow, error) {
	q := fmt.Sprintf("SELECT id, ip, queryPort FROM %s WHERE queryPort IS NOT NULL", table)
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ServerRow
	for rows.Next() {
		var s ServerRow
		if err := rows.Scan(&s.ID, &s.IP, &s.QueryPort); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func GetNextServer(db *sql.DB, table string) (ServerRow, error) {
	q := fmt.Sprintf("SELECT id, ip, queryPort FROM %s WHERE queryPort IS NOT NULL ORDER BY lastUpdate LIMIT 1", table)
	var s ServerRow
	row := db.QueryRow(q)
	if err := row.Scan(&s.ID, &s.IP, &s.QueryPort); err != nil {
		return s, err
	}
	return s, nil
}

func UpdateServerStatus(db *sql.DB, table string, id int, players int, maxPlayers int) error {
	// very lazy implementation of average players, but it works well enough for our use case
	q := fmt.Sprintf("UPDATE %s SET players = ?, maxPlayers = ?, avgPlayers = (avgPlayers * 0.9 + ? * 0.1), lastUpdate = NOW() WHERE id = ?", table)
	_, err := db.Exec(q, players, maxPlayers, players, id)
	return err
}
