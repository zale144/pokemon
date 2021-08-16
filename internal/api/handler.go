package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type pokemonCatService interface {
	GetPokeCat(pokID string, sizeLimitPx int) ([]byte, error)
}

func PokemonHandler(svc pokemonCatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		sizeLimitPx, err := strconv.ParseInt(c.DefaultQuery("size_limit_px", "0"), 10, 64)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		img, err := svc.GetPokeCat(id, int(sizeLimitPx))
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("failed to create pokemon-cat image: %s", err))
			return
		}

		c.Data(http.StatusOK, "image/png", img)
	}
}
