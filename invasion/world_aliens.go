package invasion

import (
	"math/rand"
)

// tAlienNextCityFn exists only for tests purpose. This variable should never
// be changed in files different from *_test.go
var tAlienMoveNextCityFn func(alien int, curCity *city, possibleCities []*city) *city

// createWorldAliens initializes a map with "numOfAliens" aliens.
//
func createWorldAliens(numOfAliens int) map[int]*alien {
	aliens := make(map[int]*alien, numOfAliens)
	for i := 0; i < numOfAliens; i++ {
		aliens[i] = &alien{num: i}
	}
	return aliens
}

// alien represents alien model.
type alien struct {
	num     int
	curCity *city
	trapped bool
}

// setCurCity sets the current city where the alien is.
//
func (a *alien) setCurCity(c *city) {
	a.curCity = c
}

// moveRandom moves the alien to any of its current city boundaries and
// returns the city it moved. If the alien cannot move, then this function
// returns nil.
func (a *alien) moveRandom() *city {
	if a.trapped {
		return nil
	}
	possibleCities := a.curCity.surroundingCities()
	if len(possibleCities) == 0 {
		// alien was trapped
		a.trapped = true
		return nil
	}
	if len(possibleCities) == 1 {
		a.curCity = possibleCities[0]
	} else {
		if tAlienMoveNextCityFn == nil {
			a.curCity = possibleCities[rand.Intn(len(possibleCities))]
		} else {
			a.curCity = tAlienMoveNextCityFn(a.num, a.curCity, possibleCities)
		}
	}
	return a.curCity
}
