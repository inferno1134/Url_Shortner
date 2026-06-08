package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

//created a pool for implementing the postgres
//added basci structure for now

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(dsn string) (*PostgresStore, error) {

	//creating pool 
	//not a instance

	var pool *pgxpool.Pool
	var err error

	for i := 0; i < 15; i++ {

		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, err = pgxpool.New(ctx, dsn)
		cancel()
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = pool.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS urls (
            id BIGSERIAL PRIMARY KEY,
            code TEXT UNIQUE,
            url TEXT NOT NULL
        );
    `)
	if err != nil {
		pool.Close()
		return nil, err
	}

	return &PostgresStore{pool: pool}, nil
}

///we need to make seperate funcionsfor  save and get same logic as earlier

//base 62 enocding function copied as earlier

func toBase62(num int64) string {

	//same charset varibale
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const base = 62

	if num == 0 {
		return "0"
	}

	result := ""
	for num > 0 {
		remainder := num % base
		result = string(charset[remainder]) + result
		num = num / base
	}
	return result
}

func (s *PostgresStore) Save(longUrl string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int64
	err := s.pool.QueryRow(ctx, `INSERT INTO urls (url) VALUES ($1) RETURNING id`, longUrl).Scan(&id)
	if err != nil {
		return fmt.Sprintf("err-%d", time.Now().UnixNano())
	}

	//doing the base 62 encoding here in the db itself
	//kept it simple for nw
	code := toBase62(id)
	_, err = s.pool.Exec(ctx, `UPDATE urls SET code=$1 WHERE id=$2`, code, id)
	if err != nil {
		return code
	}

	return code
}

func (s *PostgresStore) Get(shortCode string) (string, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var url string
	err := s.pool.QueryRow(ctx, `SELECT url FROM urls WHERE code=$1`, shortCode).Scan(&url)
	if err != nil {
		return "", false
	}
	return url, true
}
