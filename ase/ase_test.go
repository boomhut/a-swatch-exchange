package ase

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Make real HTTP requests for the color name lookup
func TestGetColorName(t *testing.T) {
	//	lookup the color name
	name := GetColorName("#FF0000")
	assert.Equal(t, "Red", name)

}

// Test the creation of a color struct
func TestColor(t *testing.T) {
	c := Color{
		Name: "Red",
		R:    255,
		G:    0,
		B:    0,
		Hex:  "#FF0000",
	}

	assert.Equal(t, "Red", c.Name)
	assert.Equal(t, uint8(255), c.R)
	assert.Equal(t, uint8(0), c.G)
	assert.Equal(t, uint8(0), c.B)
	assert.Equal(t, "#FF0000", c.Hex)
}

// Test the encoding and decoding of an ASE file
func TestEncodeDecode(t *testing.T) {
	colors := []Color{
		{
			Name: "Red",
			R:    255,
			G:    0,
			B:    0,
			Hex:  "#FF0000",
		},
		{
			Name: "Green",
			R:    0,
			G:    255,
			B:    0,
			Hex:  "#00FF00",
		},
		{
			Name: "Blue",
			R:    0,
			G:    0,
			B:    255,
			Hex:  "#0000FF",
		},
	}

	//	encode the colors
	err := ExportASE(colors, "test.ase")
	assert.Nil(t, err)

	//	decode the colors
	decoded, err := DecodeASE("test.ase")
	assert.Nil(t, err)

	//	check the colors
	assert.Equal(t, colors, decoded)
}

// Test the export to a JSON file
func TestExportJSON(t *testing.T) {
	colors := []Color{
		{
			Name: "Red",
			R:    255,
			G:    0,
			B:    0,
			Hex:  "#FF0000",
		},
		{
			Name: "Green",
			R:    0,
			G:    255,
			B:    0,
			Hex:  "#00FF00",
		},
		{
			Name: "Blue",
			R:    0,
			G:    0,
			B:    255,
			Hex:  "#0000FF",
		},
	}

	//	export the colors
	err := ExportJSON(colors, "test.json")
	assert.Nil(t, err)

	// read the file back in
	decoded, err := ImportJSON("test.json")
	assert.Nil(t, err)

	//	check the colors
	assert.Equal(t, colors, decoded)

}

// Clean up the test files
func TestCleanUp(t *testing.T) {

	// sleep for a second to allow the files to be closed
	// before we try to delete them
	// time.Sleep(1 * time.Second)

	err := os.Remove("./test.ase")
	assert.Nil(t, err)

	err = os.Remove("./test.json")
	assert.Nil(t, err)
}
