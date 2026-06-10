package service

import (
	"fmt"
	"log"
)

// "crypto/rand" it was earlier used for random generation of code

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

//changed this to use interface here 
//as in now it will chek what kind of sevice we are using for the databse

type Store interface {
	Save(longUrl string) (string,error)
	Get(shortCode string) (string, error)
}

type ShortnerService struct {
	store Store
}

func NewShortnerService(  store Store) *ShortnerService {

	return &ShortnerService{
		store: store,
	}
}

func (s *ShortnerService) CreateShortURL  (longUrl string) (string, error) {
	
	log.Printf("[Service]: Creating Short Url")

	shortCode, err := s.store.Save(longUrl)

	if err!=nil {
		return "", fmt.Errorf("Error Creating short Url: %w", err)
	}

	log.Printf("[Service]: Shortcode is Generated")

	return shortCode, nil
}

func (s *ShortnerService) GetOriginalURL( shortCode string) (string, error) {

	log.Printf("[Service]: Fetching the Original Url")

	longUrl ,err := s.store.Get(shortCode)

	if err!=nil {
		return "", fmt.Errorf("Error retrieving original url :%w", err)
	}

	return longUrl,nil
}
