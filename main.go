package main

import (
	"log"

	"github.com/boomhut/a-swatch-exchange/ase"
)

func main() {

	// 	Decode the ASE file
	colors, err := ase.DecodeASE("palette.ase")
	if err != nil {
		log.Println(err)
	}

	// 	Print the colors
	for _, c := range colors {
		log.Printf("%+v\n", c)
	}

	// export the colors to a json file
	err = ase.ExportJSON(colors, "palette.json")
	if err != nil {
		log.Println(err)
	}

}
