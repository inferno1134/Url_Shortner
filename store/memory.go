package store 

import (
	"sync"
)

type MemoryStore struct {

	urls map[string ] string 

	mu sync.RWMutex


}


func NewMemoryStore() *MemoryStore {

	return &MemoryStore{
		urls: make(map[string]string),
	}
}


func (s *MemoryStore) Save(shortCode string , longUrl string ){
	s.mu.Lock()
	defer s.mu.Unlock()


	s.urls[shortCode]=longUrl
}


func (s *MemoryStore) Get (shortCode string)(string,bool){

	s.mu.RLock()
	defer s.mu.RUnlock()
	url, exists:= s.urls[shortCode]

	return url, exists

}
