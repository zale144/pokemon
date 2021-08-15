package service

import (
	"fmt"
	"log"
)

type Cat struct {
	client         catClient
	cache          catCache
	isLimitReached requestLimiter
}

type requestLimiter func() bool

func NewCat(client catClient, cache catCache, limitFn requestLimiter) Cat {
	return Cat{
		client:         client,
		cache:          cache,
		isLimitReached: limitFn,
	}
}

func (c Cat) GetImage(sizeLimitPx int) ([]byte, string, error) {
	var (
		img []byte
		id  string
		err error
	)
	if c.isLimitReached != nil && c.isLimitReached() {
		log.Println("limit for calling external API reached - fallback to cache")
		return c.getRandomFromCache()
	}

	id, url, err := c.client.GetRandomCat(sizeLimitPx)
	if err != nil {
		log.Printf("failed to get random cat image ID from external API: %s - fallback to cache\n", err)
		return c.getRandomFromCache()
	}

	img, err = c.cache.Get(id)
	if err == nil {
		return img, id, nil
	}

	log.Printf("cat with id %s not found in cache - getting from external API\n", id)

	img, err = c.client.GetCatImage(url)
	if err != nil {
		log.Printf("failed to get random cat image from external API: %s - fallback to cache\n", err)
		return c.getRandomFromCache()
	}

	if err = c.cache.Set(id, img); err != nil {
		log.Printf("failed to cache cat %s: %s", id, err)
	}

	return img, id, nil
}

func (c Cat) getRandomFromCache() ([]byte, string, error) {
	img, id, err := c.cache.Random()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get random image from cache: %w", err)
	}
	return img, id, nil
}
