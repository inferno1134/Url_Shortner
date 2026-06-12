package service

import (
	"fmt"

	"github.com/tanishk-deore/url-shortner/store"
)

type ShortnerService struct {
	db    *store.PostgresStore
	cache *store.RedisCache
}

func NewShortnerService(db *store.PostgresStore, cache *store.RedisCache) *ShortnerService {
	return &ShortnerService{db: db, cache: cache}
}

func (s *ShortnerService) CreateShortURL(longUrl string) (string, error) {
	if s == nil || s.db == nil {
		return "", fmt.Errorf("storage not initialized")
	}

	code, err := s.db.Save(longUrl)
	if err != nil {
		return "", err
	}

	if s.cache != nil {
		_ = s.cache.Set(code, longUrl)
	}

	return code, nil
}

func (s *ShortnerService) GetOriginalURL(shortCode string) (string, error) {
	if s == nil || s.db == nil {
		return "", fmt.Errorf("storage not initialized")
	}

	if s.cache != nil {
		if v, err := s.cache.Get(shortCode); err == nil && v != "" {
			return v, nil
		}
	}

	longUrl, err := s.db.Get(shortCode)
	if err != nil {
		return "", err
	}

	if s.cache != nil {
		_ = s.cache.Set(shortCode, longUrl)
	}

	return longUrl, nil
}
