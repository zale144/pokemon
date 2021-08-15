// +build integration

package client

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"
)

func TestPokemon_GetPokemon(t *testing.T) {
	type fields struct {
		baseURL string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				baseURL: defaultPokemonURL,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewPokemon()
			c.url = tt.fields.baseURL

			rand.Seed(time.Now().Unix())
			id := fmt.Sprint(rand.Intn(150) + 1)

			got, err := c.GetPokemon(id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPokemon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) == 0 {
				t.Errorf("GetPokemon() len(image) = 0")
				return
			}

			if err = ioutil.WriteFile(id+".png", got, 0777); err != nil {
				t.Errorf("GetPokemon() writeFile = %s", err)
				return
			}
		})
	}
}
