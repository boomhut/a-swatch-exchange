package ase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ARolek/ase"
)

// type Color to store the color information in a struct format
type Color struct {
	Name string `json:"name,omitempty" xml:"name,omitempty" yaml:"name,omitempty"` // look for the name of the color on colorname.org and add it here
	R    uint8  `json:"r" xml:"r" yaml:"r" validate:"required"`                    // Red value of the color (0-255)
	G    uint8  `json:"g" xml:"g" yaml:"g" validate:"required"`                    // Green value of the color (0-255)
	B    uint8  `json:"b" xml:"b" yaml:"b" validate:"required"`                    // Blue value of the color (0-255)
	Hex  string `json:"hex" xml:"hex" yaml:"hex" validate:"required"`              // Hex code of the color
}

// Function to decode an ASE file and return a slice of colors.
func DecodeASE(file string) (colors []Color, err error) {
	//	open the file
	f, err := os.Open(file)
	if err != nil {
		return
	}

	defer f.Close()

	//	decode can take in any io.Reader
	ase, err := ase.Decode(f)
	if err != nil {
		return
	}

	//  get the colors from the ase
	for _, c := range ase.Colors {
		// {nameLen:8 Name:#993955 Model:RGB Values:[0.6 0.22352941 0.33333334] Type:Global}

		// check if color has already a real name (not just the hex code) and try to get the name from colorname.org
		// first calc the hex code to compare it with the name
		if c.Name[0] == '#' {
			// starts with #, so it's a hex code.
			// compare the name with the hex code
			if c.Name == Hex(uint8(c.Values[0]*255), uint8(c.Values[1]*255), uint8(c.Values[2]*255)) {
				// it's a hex code, so get the name from colorname.org
				c.Name = GetColorName(c.Name)
			}

		}

		//  convert the color to our struct
		colors = append(colors, Color{
			Name: c.Name,
			R:    uint8(c.Values[0] * 255),
			G:    uint8(c.Values[1] * 255),
			B:    uint8(c.Values[2] * 255),
			Hex:  "#" + toHex(uint8(c.Values[0]*255)) + toHex(uint8(c.Values[1]*255)) + toHex(uint8(c.Values[2]*255)),
		})
	}

	return
}

// Lookup the color name on colorname.org and add it to the struct
func GetColorName(hex string) string {
	// get the color name from colorname.org
	// https://colornames.org/search/json/?hex=FF0000
	// Returns:
	// {"hexCode":"ff0000","name":"Red"}

	// create http client
	client := getHTTPClient()

	// make the request
	resp, err := client.Get("https://colornames.org/search/json/?hex=" + hex[1:])
	if err != nil {
		return ""
	}

	// close the response body
	defer resp.Body.Close()

	// decode the response
	var colorName struct {
		Hex  string `json:"hexCode"`
		Name string `json:"name"`
	}

	// decode the response
	if err := json.NewDecoder(resp.Body).Decode(&colorName); err != nil {
		return ""
	}

	return colorName.Name

}

// Function to export the colors to a JSON file
func ExportJSON(colors []Color, file string) (err error) {
	//	open the file
	f, err := os.Create(file)
	if err != nil {
		return
	}

	defer f.Close()

	//	encode the colors
	if err := json.NewEncoder(f).Encode(colors); err != nil {
		return err
	}

	return
}

// Import colors from a JSON file
func ImportJSON(file string) (colors []Color, err error) {
	//	open the file
	f, err := os.Open(file)
	if err != nil {
		return
	}

	defer f.Close()

	//	decode the colors
	if err := json.NewDecoder(f).Decode(&colors); err != nil {
		return nil, err
	}

	return
}

// Function to export the colors, including the color name, to a new ASE file
func ExportASE(colors []Color, file string) (err error) {
	//	create a new ASE file
	aseFile := ase.ASE{}

	//	add the colors to the ASE file
	for _, c := range colors {
		//	convert the color to the ase.Color type
		aseColor := ase.Color{
			Name:   c.Hex,
			Model:  "RGB",
			Values: []float32{float32(c.R) / 255, float32(c.G) / 255, float32(c.B) / 255},
			Type:   "Global",
		}

		//	add the color to the file
		aseFile.Colors = append(aseFile.Colors, aseColor)
	}

	//	open the file
	f, err := os.Create(file)
	if err != nil {
		return
	}

	defer f.Close()

	//	encode the file
	if err := ase.Encode(aseFile, f); err != nil {
		return err
	}

	return
}

// HTTP client
func getHTTPClient() http.Client {
	return http.Client{
		Timeout: time.Second * 5,
	}

}

// function to calc the hex code from the RGB values
func Hex(r, g, b uint8) string {
	// convert the color to hex code and return it, withouth using ase package
	// get the hex value for the red, green and blue values
	// return the hex value
	return "#" + toHex(r) + toHex(g) + toHex(b)

}

// function to convert a uint8 to a hex string
func toHex(v uint8) string {
	// convert the uint8 to a hex string
	// return the hex value
	return fmt.Sprintf("%02X", v)
}
