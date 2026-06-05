package store 

import (
	"sync"
)

type MemoryStore struct {

	urls map[string ] string 

	//adding auto incrementing id here 
	counter int64

	mu sync.RWMutex


}


func NewMemoryStore() *MemoryStore {

	return &MemoryStore{
		urls: make(map[string]string),
		//adding and initialize the countere to 0
		counter: 0,
	}
}


//adding the new function to convert the integer to base 62 string 

func (s *MemoryStore) toBase62(num int64) string {

	//base charset variable to store the string for the base 62 encoding 
	const charset="0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const base = 62

	if num ==0 {
		return  "0"
	}

	result:=""

	for num >0 {

		//getting the last digit in base62
		remainder:=num%base

		//prepend the character in result
		result =string(charset[remainder])+result

		//now removing the last ldigit added above in the result
		num=num/base
	}

	return result


}

//old save method which was saving the url code in memory it was taking shortcode parameter
// func (s *MemoryStore) Save(shortCode string , longUrl string ){
// 	s.mu.Lock()
// 	defer s.mu.Unlock()


// 	s.urls[shortCode]=longUrl
// }


func (s *MemoryStore) Save(longUrl string ) string{
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	shortCode := s.toBase62(s.counter)



	s.urls[shortCode]=longUrl

	return shortCode
}

func (s *MemoryStore) Get (shortCode string)(string,bool){

	s.mu.RLock()
	defer s.mu.RUnlock()
	url, exists:= s.urls[shortCode]

	return url, exists

}
