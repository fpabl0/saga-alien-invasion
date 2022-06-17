package invasion

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
)

// WorldMap represents the simulated world map.
type WorldMap struct {
	cities map[string]*city
}

// ParseWorldMap parses the simulated world map.
//
func ParseWorldMap(s *bufio.Scanner) (*WorldMap, error) {

	wmap := &WorldMap{cities: make(map[string]*city)}

	for s.Scan() {
		// read one line
		line := s.Text()
		if line == "" {
			continue
		}

		data, err := decodeMapLine(line)
		if err != nil {
			return nil, err
		}
		curCity := wmap.getOrCreateCity(data.name)
		for i, d := range data.dirs {
			curCity.dirs[i] = wmap.getOrCreateCity(d)
		}
		for i, surCity := range curCity.dirs {
			if surCity == nil {
				continue
			}
			oppositeDir := direction(i).opposite()
			if surCity.dirs[oppositeDir] == nil {
				surCity.dirs[oppositeDir] = curCity
			} else if surCity.dirs[oppositeDir].name != curCity.name {
				return nil, errors.New("Cannot parse the map: inconsistent map")
			}
		}
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return wmap, nil
}

// getOrCreateCity gets or creates a city with the specified name in the map.
//
func (m *WorldMap) getOrCreateCity(name string) *city {
	if name == "" {
		return nil
	}
	c, ok := m.cities[name]
	if !ok {
		c = &city{name: name}
		m.cities[name] = c
	}
	return c
}

// destroyCity destroys a city in the world map.
//
func (m *WorldMap) destroyCity(name string) {
	c, ok := m.cities[name]
	if !ok {
		return
	}
	for i, sc := range c.dirs {
		if sc == nil {
			continue
		}
		sc.dirs[direction(i).opposite()] = nil
	}
	delete(m.cities, name)
}

// print prints the world map.
//
func (m *WorldMap) print(out io.Writer) {
	if len(m.cities) == 0 {
		fmt.Fprintln(out, "No cities in the map.")
		return
	}
	cs := make([]*city, 0, len(m.cities))
	for _, c := range m.cities {
		cs = append(cs, c)
	}
	sort.SliceStable(cs, func(i, j int) bool {
		return cs[i].name < cs[j].name
	})
	for _, c := range cs {
		fmt.Fprintln(out, c)
	}
}

// ===============================================================
// Utils
// ===============================================================

// rawCityData represents the raw city data in the map
type rawCityData struct {
	name string
	dirs [4]string
}

// decodeMapLine decode a single line from a map file.
//
func decodeMapLine(line string) (*rawCityData, error) {
	l := strings.TrimSpace(line)
	parts := strings.Fields(l)
	if len(parts) > 5 {
		return nil, errors.New("Invalid map line: invalid format - wrong number of spaces")
	}
	d := &rawCityData{}
	for _, p := range parts {
		m := strings.Split(p, "=")
		if len(m) > 2 {
			return nil, errors.New("Invalid map line: invalid number of '=' characters")
		}
		// get the name
		if len(m) == 1 {
			if d.name != "" {
				return nil, errors.New("Invalid map line: city with many names?")
			}
			d.name = m[0]
			continue
		}
		// find the directions
		dir, err := directionFromString(m[0])
		if err != nil {
			return nil, fmt.Errorf("Invalid map line: %v", err)
		}
		if d.dirs[dir] != "" {
			return nil, fmt.Errorf("Invalid map line: multiple %s directions", m[0])
		}
		d.dirs[dir] = m[1]
	}
	return d, nil
}
