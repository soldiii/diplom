package store

type Store struct {
	Config *StoreConfig
}

func NewStore(cfg *StoreConfig) *Store {
	return &Store{
		Config: cfg,
	}
}
