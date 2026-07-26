package core

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ShortenedUrl struct {
	Url       string `gorm:"primaryKey"`
	Original  string
	Status    string // active, expired
	Clicks    uint
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostgresDB struct {
	db *gorm.DB
}

func NewPostgresDB(host, port string,
	user, password, database string,
	sslmode, timezone string) (*PostgresDB, error) {

	db, err := gorm.Open(postgres.Open(
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
			host, port, user, password, database, sslmode, timezone)))
	if err != nil {
		return nil, fmt.Errorf("Failed to open postgres database. Error: %w", err)
	}

	initializeDB(db)

	return &PostgresDB{db}, nil
}

func initializeDB(db *gorm.DB) {
	db.AutoMigrate(&ShortenedUrl{})
}

func (pg *PostgresDB) Save(ctx context.Context,
	original, shortened, status string,
	expiresAt time.Time) error {

	row := ShortenedUrl{
		Url:       shortened,
		Original:  original,
		ExpiresAt: expiresAt,
		Status:    status,
	}

	if err := gorm.G[ShortenedUrl](pg.db).Create(ctx, &row); err != nil {
		return err
	}

	return nil
}
