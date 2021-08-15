package main

import (
	"log"
	"time"

	"pokemon/internal/api"
	"pokemon/internal/pkg/client"
	"pokemon/internal/service"
	"pokemon/pkg/memcache"
	"pokemon/pkg/request"
	"pokemon/pkg/unit"
)

func main() {
	pokemonClient := client.NewPokemon()
	cacheSizeMB := 100 * unit.MB
	pokemonCache := memcache.NewMemcache(cacheSizeMB)
	pokemonSvc := service.NewPokemon(pokemonClient, pokemonCache)

	maxCatRequests := 10
	limitPeriod := time.Duration(60)
	reqLimiter := request.FrequencyLimiter(maxCatRequests, limitPeriod)

	catCache := memcache.NewMemcache(cacheSizeMB)
	catClient := client.NewCat()

	catSvc := service.NewCat(catClient, catCache, reqLimiter)
	cacheSvc := memcache.NewMemcache(cacheSizeMB)

	svc := service.NewPokemonCat(pokemonSvc, catSvc, cacheSvc)
	log.Fatal(api.Server(svc))
}
