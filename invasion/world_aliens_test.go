package invasion

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlien_createWorldAliens(t *testing.T) {

	m := createWorldAliens(120)

	// the returned map should have the specified size
	assert.Len(t, m, 120)
	// alien objects inside the map should match the key with
	// its number
	for anum, a := range m {
		assert.NotNil(t, a)
		assert.Equal(t, anum, a.num)
		assert.Nil(t, a.curCity)
	}

}

func TestAlien_trapped(t *testing.T) {
	a := &alien{
		num:     1,
		curCity: &city{name: "C1"},
	}

	// when the current city does not have surrounding cities,
	// the alien is trapped.
	movedCity := a.moveRandom()
	assert.Nil(t, movedCity)
	assert.True(t, a.trapped)

	// after trapped the first time, next calls to moveRandom wil return
	// nil city
	movedCity = a.moveRandom()
	assert.Nil(t, movedCity)
	assert.True(t, a.trapped)
}

func TestAlien_moveRandom(t *testing.T) {

	wm := parseSmallMap(t)

	a := &alien{num: 1, curCity: wm.cities["C7"]}

	tAlienMoveNextCityFn = func(alien int, curCity *city, possibleCities []*city) *city {
		require.Equal(t, 1, alien)
		assert.Equal(t, wm.cities["C7"], curCity)
		// at C7 city, the alien can only move to 2 directions
		require.Len(t, possibleCities, 2)
		assert.Contains(t, possibleCities, wm.cities["C4"])
		return wm.cities["C4"]
	}
	movedCity := a.moveRandom()
	assert.Equal(t, wm.cities["C4"], movedCity)

	tAlienMoveNextCityFn = func(alien int, curCity *city, possibleCities []*city) *city {
		require.Equal(t, 1, alien)
		assert.Equal(t, wm.cities["C4"], curCity)
		// at C4 city, the alien can only move to 3 directions
		require.Len(t, possibleCities, 3)
		assert.Contains(t, possibleCities, wm.cities["C5"])
		return wm.cities["C5"]
	}
	movedCity = a.moveRandom()
	assert.Equal(t, wm.cities["C5"], movedCity)

	// at C5 city, alien can move in 4 directions, however if C4, C8 and C6 are destroyed
	// the alien can only move to C2
	wm.destroyCity("C4")
	wm.destroyCity("C8")
	wm.destroyCity("C6")
	tAlienMoveNextCityFn = nil
	movedCity = a.moveRandom()
	assert.Equal(t, wm.cities["C2"], movedCity)

	// if surrounding cities to C2 are destroyed (C1, C5 and C3) then alien can't move and
	// it will be trapped
	wm.destroyCity("C1")
	wm.destroyCity("C5")
	wm.destroyCity("C3")
	tAlienMoveNextCityFn = nil
	movedCity = a.moveRandom()
	assert.Nil(t, movedCity)
	assert.Equal(t, wm.cities["C2"], a.curCity)
	assert.True(t, a.trapped)

}
