package service

import (
	"crypto/rand"
	"url-shortner/store"
)

const charset="abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type ShortnerService struct {
	store *store.MemoryStore
}

func NewShortnerService (store *store.MemoryStore) *ShortnerService{

	return &ShortnerService{
		store:store,
	}



}


func generateShortCode (length int)(string , error){


	bytes := make([]byte, length )

	_, err:= rand.Read(bytes)

	if err!= nil {
		return "", err
	}


	for i := range bytes {
		bytes[i]= charset[int(bytes[i])%len(charset)]
	}

	return string(bytes),nil

}

func (s *ShortnerService) CreateShortURL (longUrl string ) (string, error){

	shortCode, err:= generateShortCode(6)

	if err!= nil {
		return "", err
	}

	s.store.Save(shortCode,longUrl)

	return shortCode, nil



}

func (s *ShortnerService) GetOriginalURL(shortCode string) (string , bool){

	return s.store.Get(shortCode)


}