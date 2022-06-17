package invasion

import (
	"fmt"
	"strings"
)

// city defines a city in the map with its surroundings.
type city struct {
	name string
	dirs [4]*city
}

// surroundingCities returns the surrounding cities using the dirs field.
//
func (c *city) surroundingCities() []*city {
	ret := make([]*city, 0, 4)
	for _, d := range c.dirs {
		if d == nil {
			continue
		}
		ret = append(ret, d)
	}
	return ret
}

// reachedCities returns all possible cities that can be reached
// from this city.
func (c *city) reachedCities() map[string]struct{} {
	rc := make(map[string]struct{})
	rc = reachedCitiesFrom(rc, c)
	delete(rc, c.name)
	return rc
}

// String implements fmt.Stringer.
// This string follow the map file format for a city.
func (c *city) String() string {
	sb := strings.Builder{}
	sb.WriteString(c.name)
	for i, surCity := range c.dirs {
		if surCity == nil {
			continue
		}
		sb.WriteByte(' ')
		sb.WriteString(fmt.Sprintf("%s=%s", direction(i), surCity.name))
	}
	return sb.String()
}

// reachedCitiesFrom returns all the possible cities that can be reached from the `fromCity`
// including `fromCity`.
func reachedCitiesFrom(reachedCities map[string]struct{}, fromCity *city) map[string]struct{} {
	if reachedCities == nil || fromCity == nil {
		return nil
	}
	for _, sc := range fromCity.dirs {
		if sc == nil {
			continue
		}
		if _, reached := reachedCities[sc.name]; reached {
			continue
		}
		reachedCities[sc.name] = struct{}{}
		reachedCities = reachedCitiesFrom(reachedCities, sc)
	}
	return reachedCities
}
