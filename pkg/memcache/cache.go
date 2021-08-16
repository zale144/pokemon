package memcache

import (
	"errors"
	"log"

	"pokemon/pkg/unit"
)

type Memcache struct {
	getCh          chan getReq
	randCh         chan randReq
	setCh          chan setReq
	delCh          chan delReq
	maxCacheSizeKb int
}

func NewMemcache(maxCacheSizeKb int) *Memcache {
	provider := Memcache{
		getCh:          make(chan getReq),
		randCh:         make(chan randReq),
		setCh:          make(chan setReq),
		delCh:          make(chan delReq),
		maxCacheSizeKb: maxCacheSizeKb,
	}
	go provider.cache()
	return &provider
}

type item struct {
	key string
	raw []byte
}

func (mp *Memcache) Get(key string) ([]byte, error) {
	req := getReq{
		key:  key,
		resp: make(chan item),
	}
	mp.getCh <- req
	d := <-req.resp
	if len(d.raw) == 0 {
		return nil, errors.New("cache miss")
	}
	return d.raw, nil
}

func (mp *Memcache) Random() ([]byte, string, error) {
	req := randReq{
		resp: make(chan item),
	}
	mp.randCh <- req
	d := <-req.resp
	if len(d.raw) == 0 {
		return nil, d.key, errors.New("cache miss")
	}
	return d.raw, d.key, nil
}

func (mp *Memcache) Set(key string, data []byte) error {
	req := setReq{
		data: item{
			key: key,
			raw: data,
		},
		key:  key,
		resp: make(chan error),
	}
	mp.setCh <- req
	return <-req.resp
}

func (mp *Memcache) Delete(key string) error {
	req := delReq{
		key: key,
	}
	mp.delCh <- req
	return nil
}

type getReq struct {
	key  string
	resp chan item
}

type randReq struct {
	resp chan item
}

type setReq struct {
	data item
	key  string
	resp chan error
}

type delReq struct {
	key  string
	resp chan error
}

func (mp *Memcache) cache() {
	m := make(map[string]item)

	for {
		select {
		case gd := <-mp.getCh:
			gd.resp <- m[gd.key]
		case rd := <-mp.randCh:
			var i item
			for k := range m {
				i = m[k]
				break
			}
			rd.resp <- i
		case sd := <-mp.setCh:
			if mp.maxCacheSizeKb > 0 {
				size := getCacheSizeKB(m)
				if int(size)/unit.KB > mp.maxCacheSizeKb {
					for k := range m {
						log.Printf("cache full, removing item: %s\n", k)
						delete(m, k)
						break
					}
				}
			}
			m[sd.key] = sd.data
			sd.resp <- nil
		case d := <-mp.delCh:
			delete(m, d.key)
			d.resp <- nil
		}
	}
}

func getCacheSizeKB(mc map[string]item) int {
	l := 0
	for _, v := range mc {
		l += len(v.raw)
	}
	return l * unit.KB
}
