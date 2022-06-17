package invasion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCity_surroundingCities(t *testing.T) {

	c := &city{
		name: "C1",
		dirs: [4]*city{
			dirNorth: {name: "C2"},
			dirSouth: {name: "C3"},
			dirEast:  {name: "C4"},
			dirWest:  {name: "C5"},
		},
	}

	surCities := c.surroundingCities()
	assert.Equal(t, []*city{
		{name: "C2"}, {name: "C3"}, {name: "C4"}, {name: "C5"},
	}, surCities)

	// -- case after removing some directions

	c.dirs[dirEast] = nil
	c.dirs[dirNorth] = nil

	surCities = c.surroundingCities()

	assert.Equal(t, []*city{
		{name: "C3"}, {name: "C5"},
	}, surCities)

}

func TestCity_String(t *testing.T) {
	c := &city{
		name: "C1",
		dirs: [4]*city{
			dirNorth: {name: "C2"},
			dirSouth: {name: "C3"},
			dirEast:  {name: "C4"},
			dirWest:  {name: "C5"},
		},
	}

	assert.Equal(t, "C1 north=C2 south=C3 east=C4 west=C5", c.String())

	// -- case after removing some directions

	c.dirs[dirEast] = nil
	c.dirs[dirNorth] = nil

	assert.Equal(t, "C1 south=C3 west=C5", c.String())
}

func TestCity_reachedCities(t *testing.T) {

	wm := parseSmallMap(t)

	reachedCities := wm.cities["C4"].reachedCities()
	assert.Equal(t, map[string]struct{}{
		"C1": {}, "C2": {}, "C3": {}, "C5": {},
		"C6": {}, "C7": {}, "C8": {}, "C9": {},
	}, reachedCities)

	reachedCities = wm.cities["C1"].reachedCities()
	assert.Equal(t, map[string]struct{}{
		"C2": {}, "C3": {}, "C4": {}, "C5": {},
		"C6": {}, "C7": {}, "C8": {}, "C9": {},
	}, reachedCities)

	wm.destroyCity("C3")
	wm.destroyCity("C5")
	wm.destroyCity("C8")

	reachedCities = wm.cities["C1"].reachedCities()
	assert.Equal(t, map[string]struct{}{
		"C2": {}, "C4": {}, "C7": {},
	}, reachedCities)

	wm.destroyCity("C4")

	reachedCities = wm.cities["C1"].reachedCities()
	assert.Equal(t, map[string]struct{}{"C2": {}}, reachedCities)

	wm.destroyCity("C2")

	reachedCities = wm.cities["C1"].reachedCities()
	assert.Equal(t, map[string]struct{}{}, reachedCities)
}
