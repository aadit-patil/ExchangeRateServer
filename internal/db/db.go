package db

//go:generate mockgen -destination=../mocks/mock_db.go -package=mocks exchange/internal/db Database

type Database interface {
	GetRate(from, to, date string) (float64, error)
	InsertRate(from, to, date string, rate float64) error
}

var DBImpl Database

func SetDatabase(db Database) {
	DBImpl = db
}
