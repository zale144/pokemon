package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"pokemon/pkg/unit"
)

type Cat struct {
	url    string
	client *http.Client
}

const (
	defaultCatURL     = "https://api.thecatapi.com/v1/images/search?mime_types=png"
	defaultCatTimeout = 5 * time.Second
	maxIDAttempts     = 10
)

func NewCat() Cat {
	return Cat{
		url: defaultCatURL,
		client: &http.Client{
			Timeout: defaultCatTimeout,
		},
	}
}

func (c Cat) GetCatImage(url string) ([]byte, error) {
	start := time.Now()
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("failed to get image: code %d; status: %s", resp.StatusCode, resp.Status))
	}

	out := new(bytes.Buffer)
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	body := out.Bytes()
	log.Printf("downloaded cat image in %f seconds", time.Since(start).Seconds())
	log.Printf("cat image size is %d KB", len(body)/unit.KB)

	return body, nil
}

type searchResponse struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func (c Cat) GetRandomCat(sizeLimitPx int) (string, string, error) {
	type result struct {
		id, url string
	}
	resCh := make(chan *result)

	wg := sync.WaitGroup{}
	wg.Add(maxIDAttempts)
	// up to `maxIDAttempts` attempts, in parallel
	for i := 0; i < maxIDAttempts; i++ {
		go func() {
			defer wg.Done()

			id, url, err := c.getRandomCat(sizeLimitPx)
			if err != nil {
				log.Printf("image id not found: %s - retrying\n", err)
				return
			}
			resCh <- &result{
				id:  id,
				url: url,
			}
		}()
	}

	start := time.Now()

	go func() {
		wg.Wait()
		close(resCh)
	}()

	response := <-resCh
	if response == nil {
		return "", "", errors.New("cat id not found for requested size")
	}
	log.Printf("image ID with requested size found on cat API in %f seconds\n", time.Since(start).Seconds())
	return response.id, response.url, nil
}

func (c Cat) getRandomCat(sizeLimitPx int) (string, string, error) {
	resp, err := c.client.Get(c.url)
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
		return "", "", fmt.Errorf("image size exceeds maximum: W = %d; H = %d", response.Width, response.Height)
	}

	return response.ID, response.URL, nil
}
