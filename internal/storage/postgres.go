package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"my_crypto_project/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const createTableSQL = `
CREATE TABLE IF NOT EXISTS prices (
    id SERIAL PRIMARY KEY,
    exchange TEXT NOT NULL,
    price NUMERIC NOT NULL,
    timestamp TIMESTAMP NOT NULL
);`

type Storage struct {
	db *sql.DB
}

func New(dsn string) (*Storage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Printf("Не удалось создать таблицу: %v", err)
		return nil, err
	}

	return &Storage{db: db}, nil

}

func (s *Storage) SavePrice(ctx context.Context, p models.PriceResult) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO prices (exchange, price, timestamp) VALUES ($1,$2,$3)",
		p.Exchange, p.Price, p.Timestamp)
	return err
}

func (s *Storage) GetPrices(ctx context.Context, interval string) ([]models.PriceResult, error) {
	query := fmt.Sprintf("SELECT exchange, price, timestamp FROM prices WHERE timestamp > NOW() - INTERVAL '%s'", interval)
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.PriceResult
	for rows.Next() {
		var p models.PriceResult
		if err := rows.Scan(&p.Exchange, &p.Price, &p.Timestamp); err != nil {
			continue
		}
		results = append(results, p)
	}
	return results, nil
}

func (s *Storage) GetStats(ctx context.Context) (models.DataForStats, error) {
	var stats models.DataForStats
	query := "SELECT COUNT(id), COALESCE(MAX(price), 0), COALESCE(MIN(price), 0) FROM prices"
	row := s.db.QueryRowContext(ctx, query)
	err := row.Scan(&stats.Count, &stats.MaxPrice, &stats.MinPrice)
	if err != nil {
		return models.DataForStats{}, fmt.Errorf("ошибка при получении стратистики: %w", err)
	}
	return stats, nil
}
