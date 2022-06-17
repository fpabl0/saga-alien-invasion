package mapgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCityName_Uniqueness(t *testing.T) {

	g := Generator{width: 20, height: 20}

	cityNames := make(map[string]struct{}, g.width*g.height)

	for i := 0; i < g.width; i++ {
		for j := 0; j < g.height; j++ {
			cn := g.getCityName(i, j)
			_, exist := cityNames[cn]
			require.False(t, exist)
			cityNames[cn] = struct{}{}
		}
	}

	require.Len(t, cityNames, g.width*g.height)

}

func TestCityLine_Create(t *testing.T) {
	g := Generator{width: 8, height: 5}

	cl := g.createCityLine(2, 1)
	assert.Equal(t, "C11 north=C3 south=C19 east=C12 west=C10", cl)

	cl = g.createCityLine(7, 2)
	assert.Equal(t, "C24 north=C16 south=C32 west=C23", cl)

	// -- cities on map boundaries only have 2 surrounding cities

	cl = g.createCityLine(0, 0)
	assert.Equal(t, "C1 south=C9 east=C2", cl)

	cl = g.createCityLine(0, 4)
	assert.Equal(t, "C33 north=C25 east=C34", cl)

	cl = g.createCityLine(7, 0)
	assert.Equal(t, "C8 south=C16 west=C7", cl)

	cl = g.createCityLine(7, 4)
	assert.Equal(t, "C40 north=C32 west=C39", cl)

}

func TestMapGeneration(t *testing.T) {
	g := Generator{width: 3, height: 3}

	data := g.Generate()

	assert.Equal(t, `C1 south=C4 east=C2
C2 south=C5 east=C3 west=C1
C3 south=C6 west=C2
C4 north=C1 south=C7 east=C5
C5 north=C2 south=C8 east=C6 west=C4
C6 north=C3 south=C9 west=C5
C7 north=C4 east=C8
C8 north=C5 east=C9 west=C7
C9 north=C6 west=C8
`, string(data))

}
