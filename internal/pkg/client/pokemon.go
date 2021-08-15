package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Pokemon struct {
	url string
}

const defaultPokemonURL = "https://pokeapi.co/api/v2/pokemon/"

func NewPokemon() Pokemon {
	return Pokemon{
		url: defaultPokemonURL,
	}
}

func (c Pokemon) GetPokemon(id string) ([]byte, error) {
	pokemonImageURL, err := c.getImageURL(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get pokemon: %w", err)
	}

	image, err := c.getImage(pokemonImageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	return image, nil
}

func (c Pokemon) getImage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("failed to get image: code %d; status: %s", resp.StatusCode, resp.Status))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

func (c Pokemon) getImageURL(id string) (string, error) {
	resp, err := http.Get(c.url + "/" + id)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("failed to search: code %d; status: %s", resp.StatusCode, resp.Status))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	result := new(restResponse)
	if err = json.Unmarshal(body, result); err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if len(result.Sprites.FrontDefault) == 0 {
		return "", errors.New("front_default is empty")
	}

	return result.Sprites.FrontDefault, nil
}

type restResponse struct {
	Sprites struct {
		FrontDefault string `json:"front_default"`
	} `json:"sprites"`
}
