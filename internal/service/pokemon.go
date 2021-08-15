package service

import (
	"fmt"
	"log"
)

type Pokemon struct {
	client pokemonClient
	cache  pokemonCache
}

func NewPokemon(client pokemonClient, cache pokemonCache) Pokemon {
	return Pokemon{
		client: client,
		cache:  cache,
	}
}

func (p Pokemon) GetImage(id string) ([]byte, error) {
	img, err := p.cache.Get(id)
	if err == nil {
		return img, nil
	}

	log.Printf("pokemon %s not found in cache - fallback to external API\n", id)

	img, err = p.client.GetPokemon(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find pokemon from external API: %w", err)
	}

	if err = p.cache.Set(id, img); err != nil {
		log.Printf("failed to cache pokemon %s: %s", id, err)
	}

	return img, nil
}
