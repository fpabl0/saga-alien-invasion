package mapgen

import (
	"bytes"
	"fmt"
	"strings"
)

// Generator represents map generator object.
type Generator struct {
	width, height int
}

// NewGenerator creates a new map generator instance.
//
func NewGenerator(width, height int) *Generator {
	return &Generator{width: width, height: height}
}

// Generate generates a world map using the width and height as the number of citys at its borders.
//
func (g *Generator) Generate() []byte {
	buf := &bytes.Buffer{}
	buf.Grow(g.width * g.height * 40)

	for j := 0; j < g.height; j++ {
		for i := 0; i < g.width; i++ {
			buf.WriteString(g.createCityLine(i, j))
			buf.WriteByte('\n')
		}
	}

	return buf.Bytes()
}

// createCityLine creates a formatted city line in the form:
// <city> north=<north_city> south=<south_city> east=<east_city> west=<west_city>
func (g *Generator) createCityLine(i, j int) string {
	north, south := j-1, j+1
	east, west := i+1, i-1

	sb := strings.Builder{}
	sb.Grow(40)

	sb.WriteString(g.getCityName(i, j))
	if north >= 0 && north < g.height {
		sb.WriteByte(' ')
		sb.WriteString(fmt.Sprintf("north=%s", g.getCityName(i, north)))
	}
	if south >= 0 && south < g.height {
		sb.WriteByte(' ')
		sb.WriteString(fmt.Sprintf("south=%s", g.getCityName(i, south)))
	}
	if east >= 0 && east < g.width {
		sb.WriteByte(' ')
		sb.WriteString(fmt.Sprintf("east=%s", g.getCityName(east, j)))
	}
	if west >= 0 && west < g.width {
		sb.WriteByte(' ')
		sb.WriteString(fmt.Sprintf("west=%s", g.getCityName(west, j)))
	}
	return sb.String()
}

// getCityName returns the city name based on the map coordinates.
//
func (g *Generator) getCityName(i, j int) string {
	return fmt.Sprintf("C%d", j*g.width+i+1)
}
