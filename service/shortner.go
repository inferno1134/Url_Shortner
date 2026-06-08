package service

// "crypto/rand" it was earlier used for random generation of code

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

//changed this to use interface here 
//as in now it will chek what kind of sevice we are using for the databse

type Store interface {
	Save(longUrl string) string
	Get(shortCode string) (string, bool)
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
	shortCode := s.store.Save(longUrl)

	return shortCode, nil
}

func (s *ShortnerService) GetOriginalURL( shortCode string) (string, bool) {
	return s.store.Get(shortCode)
}
