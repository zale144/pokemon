package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Cat struct {
	url string
}

const defaultCatURL = "https://api.thecatapi.com/v1/images/search?mime_types=png"

func NewCat() Cat {
	return Cat{
		url: defaultCatURL,
	}
}

func (c Cat) GetCatImage(url string) ([]byte, error) {
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

	err = ioutil.WriteFile("cat1.png", body, 0777)
	if err != nil {
		log.Println(err)
	}

	return body, nil
}

type searchResponse struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func (c Cat) GetRandomCat(sizeLimitPx int) (string, string, error) {
	resp, err := http.Get(c.url)
	if err != nil {
		return "", "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", "", errors.New(fmt.Sprintf("failed to search: code %d; status: %s", resp.StatusCode, resp.Status))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %w", err)
	}

	var searchResult []*searchResponse
	if err = json.Unmarshal(body, &searchResult); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	var response *searchResponse
	if len(searchResult) > 0 {
		response = searchResult[0]
	}

	if sizeLimitPx > 0 && (response.Width > sizeLimitPx || response.Height > sizeLimitPx) {
		log.Printf("image size exceeds maximum: W = %d; H = %d - retrying \n", response.Width, response.Height)
		return c.GetRandomCat(sizeLimitPx)
	}

	return response.ID, response.URL, nil
}
