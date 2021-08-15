package service

import (
	"embed"
	"io"
	"reflect"
	"testing"
)

//go:embed testdata
var testdata embed.FS

func TestPokemonCat_generate(t *testing.T) {
	type args struct {
		pokemonInputFileName string
		catInputFileName     string
	}
	tests := []struct {
		name         string
		args         args
		wantFileName string
		wantErr      bool
	}{
		{
			name: "success1",
			args: args{
				pokemonInputFileName: "pok1.png",
				catInputFileName:     "cat1.png",
			},
			wantFileName: "test_out1.png",
			wantErr:      false,
		}, {
			name: "success2",
			args: args{
				pokemonInputFileName: "pok1.png",
				catInputFileName:     "cat2.png",
			},
			wantFileName: "test_out2.png",
			wantErr:      false,
		}, {
			name: "success3",
			args: args{
				pokemonInputFileName: "pok1.png",
				catInputFileName:     "cat3.png",
			},
			wantFileName: "test_out3.png",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PokemonCat{}

			pokFile, _ := testdata.Open("testdata/" + tt.args.pokemonInputFileName)
			pokBytes, _ := io.ReadAll(pokFile)

			catFile, _ := testdata.Open("testdata/" + tt.args.catInputFileName)
			catBytes, _ := io.ReadAll(catFile)

			got, err := p.generate(pokBytes, catBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			outFile, _ := testdata.Open("testdata/" + tt.wantFileName)
			outBytes, _ := io.ReadAll(outFile)

			if !reflect.DeepEqual(got, outBytes) {
				t.Error("generate() got = != want")
			}
		})
	}
}
