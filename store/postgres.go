package store

import (
	"context"
	"fmt"
	"strings"
	"time"
	"log"

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

		log.Printf("[PostgresStore] Attempting databse Store connection : %d/15", i+1)

		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, err = pgxpool.New(ctx, dsn)
		cancel()
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("Pool Connection not established: %w", err)
	}

	log.Printf("[PostgresStore] Database Connection Established ")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("[PostgresStore] Creating Urls table...")

	_, err = pool.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS urls (
            id BIGSERIAL PRIMARY KEY,
            url TEXT NOT NULL
        );
    `)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("Pool Exec: %w",err)
	}

	log.Printf("[PostgresStore] Url table is ready")

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

func (s *PostgresStore) Save(longUrl string) (string,error) {

	log.Printf("[PostgresStore.Save] Saving Url...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int64

	
	err := s.pool.QueryRow(ctx, `INSERT INTO urls (url) VALUES ($1) RETURNING id`, longUrl).Scan(&id)
	if err != nil {
		return "",fmt.Errorf("Database Insertion Failed : %w", err)
	}

	log.Printf("[PostgresStore.Save] Url inserted in table")

	// //doing the base 62 encoding here in the db itself
	// //kept it simple for nw
	// code := toBase62(id)
	// _, err = s.pool.Exec(ctx, `UPDATE urls SET code=$1 WHERE id=$2`, code, id)
	// if err != nil {
	// 	return code
	// }

	return toBase62(id),nil
}

func (s *PostgresStore) Get(shortCode string) (string, error) {

	log.Printf("[PostgresStore.Get] Fetching the Url")


	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	id, err:= fromBase62(shortCode)

	if err!=nil{
		return "", fmt.Errorf("Error Decoding short code %q: %w", shortCode, err)
	}

	var url string 	

	err= s.pool.QueryRow(ctx,`SELECT url FROM urls WHERE id=$1`,id).Scan(&url)

	if err!= nil {
		return "", fmt.Errorf("Error querying for id %d: %w", id, err)
	}

	log.Printf("[PostgresStore.Get] Url found from base 62")

	// var url string
	// err := s.pool.QueryRow(ctx, `SELECT url FROM urls WHERE code=$1`, shortCode).Scan(&url)
	// if err != nil {
	// 	return "", false
	// }
	return url, nil
}


//this is just the fucntion to convert back the code into id and then fetch the longurl 

//get fucntion will use it to get the long url and redirect function
func fromBase62(code string )(int64 , error){

	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    const base = 62

	var value int64

	for _,ch := range code{

		idx:= strings.IndexRune(charset, ch)

		if idx < 0{
			return 0, fmt.Errorf("Invalid Base62 character : %c", ch)
		} 

		value=value*base + int64(idx)

	}

	return value,nil

	
}