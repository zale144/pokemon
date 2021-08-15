package service

import (
	"bytes"
	"fmt"
	"image"
	"log"

	"github.com/disintegration/imaging"
)

type PokemonCat struct {
	pokemon pokemonService
	cat     catService
	cache   pokemonCatCache
}

func NewPokemonCat(pokemon pokemonService, cat catService, cache pokemonCatCache) PokemonCat {
	return PokemonCat{
		pokemon: pokemon,
		cat:     cat,
		cache:   cache,
	}
}

func (p PokemonCat) GetPokeCat(pokID string, sizeLimitPx int) ([]byte, error) {
	catImageBytes, catID, err := p.cat.GetImage(sizeLimitPx)
	if err != nil {
		return nil, fmt.Errorf("failed to get cat image: %w", err)
	}

	key := getKey(pokID, catID)

	pokeCat, err := p.cache.Get(key)
	if err == nil {
		return pokeCat, nil
	}

	log.Printf("pokemon-cat %s not found in cache - generating\n", key)

	pokemonImageBytes, err := p.pokemon.GetImage(pokID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pokemon image: %w", err)
	}

	pokeCat, err = p.generate(pokemonImageBytes, catImageBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate pokemon-cat image: %w", err)
	}

	if err = p.cache.Set(key, pokeCat); err != nil {
		log.Printf("failed to cache generated pokemon-cat %s: %s", key, err)
	}

	return pokeCat, nil
}

func getKey(pokID, catID string) string {
	return fmt.Sprintf("%s-%s", pokID, catID)
}

func (PokemonCat) generate(pokemonImageBytes, catImageBytes []byte) ([]byte, error) {
	catImage, _, err := image.Decode(bytes.NewBuffer(catImageBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode cat image: %w", err)
	}

	pokemonImage, _, err := image.Decode(bytes.NewBuffer(pokemonImageBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode pokemon image: %w", err)
	}

	const (
		catAspectRatio    = 0.75
		catPokHeightRatio = 5.0
	)

	catInHeight := catImage.Bounds().Dy()
	catInWidth := catImage.Bounds().Dx()
	pokInHeight := pokemonImage.Bounds().Dy()

	catOutHeight := pokInHeight * catPokHeightRatio
	catScaleDown := float64(catOutHeight) / float64(catInHeight)
	catOutWidth := int(float64(catInWidth) * catScaleDown)

	resizedCatImage := imaging.Resize(catImage, catOutWidth, catOutHeight, imaging.Lanczos)

	catCroppedWidth := float64(catOutHeight) * catAspectRatio
	catMinX := (catOutWidth - int(catCroppedWidth)) / 2
	catMaxX := catOutWidth - catMinX

	catBounds := image.Rectangle{
		Min: image.Point{
			X: catMinX,
			Y: 0,
		},
		Max: image.Point{
			X: catMaxX,
			Y: catOutHeight,
		},
	}

	croppedCatImage := imaging.Crop(resizedCatImage, catBounds)

	catInHeight = croppedCatImage.Bounds().Dy()
	pokY := catInHeight / catPokHeightRatio
	y := catInHeight - pokY
	x := 0

	pokeCatImage := imaging.Overlay(croppedCatImage, pokemonImage, image.Pt(x, y), 1)

	w := new(bytes.Buffer)
	if err = imaging.Encode(w, pokeCatImage, imaging.PNG); err != nil {
		return nil, fmt.Errorf("failed to encode pokemon-cat image: %w", err)
	}

	return w.Bytes(), nil
}
