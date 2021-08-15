package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Pokemon struct {
	url    string
	client *http.Client
}

const (
	defaultPokemonURL     = "https://pokeapi.co/api/v2/pokemon/"
	defaultPokemonTimeout = 10 * time.Second
)

func NewPokemon() Pokemon {
	return Pokemon{
		url: defaultPokemonURL,
		client: &http.Client{
			Timeout: defaultPokemonTimeout,
		},
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
	start := time.Now()

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	since := time.Since(start)
	log.Printf("image request to pokemon API took %f seconds", since.Seconds())

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
	resp, err := c.client.Get(c.url + "/" + id)
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
