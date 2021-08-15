package service

type pokemonService interface {
	GetImage(id string) ([]byte, error)
}

type catService interface {
	GetImage(sizeLimitPx int) ([]byte, string, error)
}

type pokemonCatCache interface {
	Get(key string) ([]byte, error)
	Set(key string, data []byte) error
}

type pokemonClient interface {
	GetPokemon(key string) ([]byte, error)
}

type pokemonCache interface {
	Get(key string) ([]byte, error)
	Set(key string, data []byte) error
}

type catClient interface {
	GetRandomCat(sizeLimitPx int) (string, string, error)
	GetCatImage(url string) ([]byte, error)
}

type catCache interface {
	Random() ([]byte, string, error)
	Get(key string) ([]byte, error)
	Set(key string, data []byte) error
}
