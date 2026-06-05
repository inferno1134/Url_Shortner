package service

import (
	// "crypto/rand" it was earlier used for random generation of code 
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

//no need of this geneerateshortCode function
//now whole shortening of code is being done in store itself


// func generateShortCode (length int)(string , error){


// 	bytes := make([]byte, length )

// 	_, err:= rand.Read(bytes)

// 	if err!= nil {
// 		return "", err
// 	}


// 	for i := range bytes {
// 		bytes[i]= charset[int(bytes[i])%len(charset)]
// 	}

// 	return string(bytes),nil

// }


//this one is before base62 encoding
// func (s *ShortnerService) CreateShortURL (longUrl string ) (string, error){

// 	shortCode, err:= generateShortCode(6)

// 	if err!= nil {
// 		return "", err
// 	}

// 	s.store.Save(shortCode,longUrl)

// 	return shortCode, nil



// }

func (s *ShortnerService) CreateShortURL(longUrl string) (string,error){

	//now storing the base62 encoding 
	shortCode:= s.store.Save(longUrl)
	return shortCode, nil

}

func (s *ShortnerService) GetOriginalURL(shortCode string) (string , bool){

	return s.store.Get(shortCode)


}