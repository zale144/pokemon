// +build integration

package client

import (
	"io/ioutil"
	"testing"
)

func TestCat_GetRandomCat(t *testing.T) {
	type fields struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				url: defaultCatURL,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCat()
			c.url = tt.fields.url

			id, url, err := c.GetRandomCat(900)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRandomCat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(id) == 0 {
				t.Errorf("GetRandomCat() len(id) = 0")
			}

			image, err := c.GetCatImage(url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCatImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(image) == 0 {
				t.Errorf("GetRandomCat() len(image) = 0")
				return
			}

			if err = ioutil.WriteFile(id+".png", image, 0777); err != nil {
				t.Errorf("GetRandomCat() writeFile = %s", err)
				return
			}
		})
	}
}
