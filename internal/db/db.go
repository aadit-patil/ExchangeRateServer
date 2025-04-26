package db

import "time"

type Database interface {
	GetRate(from, to, date string) (float64, error)
	InsertRate(from, to, date string, rate float64) error
	InsertMultipleRates(base, date string, rates map[string]float64, ttl time.Time) error
}

var DBImpl Database

func SetDatabase(db Database) {
	DBImpl = db
}
