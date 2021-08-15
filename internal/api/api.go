package api

import "github.com/gin-gonic/gin"

func Server(svc pokemonCatService) error {
	r := gin.Default()
	r.GET("/api/v1/pokemon/:id", PokemonHandler(svc))
	return r.Run()
}
