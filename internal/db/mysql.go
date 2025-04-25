package db

import (
	"database/sql"
	"log"

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
	err := m.conn.QueryRow(`SELECT rate FROM exchange_rates WHERE rate_date=? AND base_currency=? AND target_currency=?`, date, from, to).Scan(&rate)
	return rate, err
}

func (m *MySQLDB) InsertRate(from, to, date string, rate float64) error {
	_, err := m.conn.Exec(`INSERT IGNORE INTO exchange_rates (rate_date, base_currency, target_currency, rate) VALUES (?, ?, ?, ?)`, date, from, to, rate)
	return err
}
