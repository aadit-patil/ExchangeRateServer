package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLDB struct {
	conn *sql.DB
}

func InitMySQL(dsn string) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("DB ping failed:", err)
	}
	SetDatabase(&MySQLDB{conn: db})
}

func (m *MySQLDB) GetRate(from, to, date string) (float64, error) {
	var rate float64
	err := m.conn.QueryRow(
		`SELECT rate FROM exchange_rates WHERE rate_date=? AND base_currency=? AND target_currency=?`,
		date, from, to).Scan(&rate)
	return rate, err
}

func (m *MySQLDB) InsertRate(from, to, date string, rate float64) error {
	_, err := m.conn.Exec(
		`INSERT IGNORE INTO exchange_rates (rate_date, base_currency, target_currency, rate) VALUES (?, ?, ?, ?)`,
		date, from, to, rate)
	return err
}

func (m *MySQLDB) InsertMultipleRates(base, date string, rates map[string]float64, ttl time.Time) error {
	tx, err := m.conn.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT IGNORE INTO exchange_rates (rate_date, base_currency, target_currency, rate, expires_at) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for target, rate := range rates {
		if _, ok := supportedCurrencies[target]; !ok {
			continue
		}
		if _, err := stmt.Exec(date, base, target, rate, ttl); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

var supportedCurrencies = map[string]bool{
	"USD": true,
	"INR": true,
	"EUR": true,
	"JPY": true,
	"GBP": true,
}
