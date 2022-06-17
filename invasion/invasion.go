package invasion

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
	"strings"
)

// tStarterCityForAlienFn is only used for test purposes and it should not be
// changed in files different from *_test.go
var tStarterCityNameForAlienFn func(alien int) string

// Start starts the invasion.
//
func Start(out io.Writer, wmap *WorldMap, numOfAliens int) {

	if wmap == nil || numOfAliens == 0 {
		return
	}

	if out == nil {
		out = os.Stdout
	}

	aliensMap := createWorldAliens(numOfAliens)

	cityNames := make([]string, 0, len(wmap.cities))
	for cn := range wmap.cities {
		cityNames = append(cityNames, cn)
	}

	// key (string) = city name, value (*alientSet) = aliens in the city
	cityAliensMap := make(map[string]*alienSet, len(cityNames))

	// unleash aliens at random cities
	for i := 0; i < numOfAliens; i++ {
		var cname string
		if tStarterCityNameForAlienFn == nil {
			rnum := rand.Intn(len(cityNames))
			cname = cityNames[rnum]
		} else {
			cname = tStarterCityNameForAlienFn(i)
		}
		aliensMap[i].setCurCity(wmap.cities[cname])
		addAlienToCity(cityAliensMap, cname, i)
	}

	// start alien invasion
	nonTrappedAlienMoves := 0
	for {

		if len(aliensMap) == 0 {
			fmt.Fprintln(out, "All the aliens have been destroyed!")
			break
		}

		if nonTrappedAlienMoves >= 10000 {
			if !remAliensCanReachEachOther(aliensMap) {
				fmt.Fprintln(out, "Remaining aliens can't reach each other!")
			} else {
				fmt.Fprintln(out, "Each non-trapped alien has moved 10000 times!")
			}
			break
		}

		trappedAliens := 0
		for _, a := range aliensMap {
			curCityName := a.curCity.name
			movedCity := a.moveRandom()
			if movedCity == nil {
				trappedAliens++
				continue
			}
			removeAlienFromCity(cityAliensMap, curCityName, a.num)
			addAlienToCity(cityAliensMap, movedCity.name, a.num)
		}

		if (len(aliensMap) - trappedAliens) <= 1 {
			if trappedAliens == len(aliensMap) {
				fmt.Fprintf(out, "All the remaining aliens (%d) were trapped!\n", trappedAliens)
			} else {
				fmt.Fprintln(out, "There is just 1 free alien, then no city can be destroyed!")
			}
			break
		}

		nonTrappedAlienMoves++

		for c, aSet := range cityAliensMap {
			if aSet.len() >= 2 {
				wmap.destroyCity(c)
				fmt.Fprintf(out, "%s has been destroyed by %s!\n", c, aSet)
				for a := range aSet.data {
					delete(aliensMap, a)
				}
				delete(cityAliensMap, c)
			}
		}

	}

	fmt.Fprintln(out, "\nResult map:")
	wmap.print(out)
}

// ===============================================================
// Util functions
// ===============================================================

// remAliensCanReachEachOther checks if at least one pair of the remaining
// aliens can reach each other.
func remAliensCanReachEachOther(aliensMap map[int]*alien) bool {
	cities := make([]*city, 0, len(aliensMap))
	for _, a := range aliensMap {
		if a.trapped {
			continue
		}
		cities = append(cities, a.curCity)
	}
	for startPos := 0; startPos < len(cities); startPos++ {
		rc := cities[startPos].reachedCities()
		for i := startPos + 1; i < len(cities); i++ {
			if _, reached := rc[cities[i].name]; reached {
				return true
			}
		}
	}
	return false
}

// addAlienToCity adds an alien to the specified city using the city-alienSet map.
// If the alienSet map for "city" is nil, this function will create a newAlienSet
// and then adds the alien.
func addAlienToCity(m map[string]*alienSet, city string, alien int) {
	if m[city] == nil {
		m[city] = newAlienSet()
	}
	m[city].add(alien)
}

// removeAlienFromCity removes the specified alien from the given city in the city-alienSet
// map. If the city does not exist in the map, then this function does nothing.
func removeAlienFromCity(m map[string]*alienSet, city string, alien int) {
	if m[city] == nil {
		return
	}
	m[city].remove(alien)
}

// ===============================================================
// Alien set type
// ===============================================================

type alienSet struct {
	data map[int]struct{}
}

// newAlienSet creates a new alien instantiating the internal data map.
func newAlienSet() *alienSet {
	return &alienSet{data: make(map[int]struct{})}
}

// len returns the alienSet length.
func (s *alienSet) len() int {
	return len(s.data)
}

// add adds one alien to the set.
func (s *alienSet) add(n int) {
	s.data[n] = struct{}{}
}

// remove removes one alien from the set.
func (s *alienSet) remove(n int) {
	delete(s.data, n)
}

// String implements fmt.Stringer. This will return
// a string with all the aliens in the form:
// 		`alien x, alien y and alien z`
// And sorted ascending.
func (s *alienSet) String() string {
	if s.len() == 0 {
		return ""
	}
	aliens := make([]int, 0, len(s.data))
	for anum := range s.data {
		aliens = append(aliens, anum)
	}
	sort.Slice(aliens, func(i, j int) bool {
		return aliens[i] < aliens[j]
	})
	sb := &strings.Builder{}
	sb.Grow(9*s.len() + 3)
	for i, a := range aliens {
		if i > 0 && i == len(aliens)-1 {
			sb.WriteString(" and ")
		} else if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("alien %d", a))
	}
	return sb.String()
}
