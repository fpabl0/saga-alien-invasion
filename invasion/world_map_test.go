package invasion

import (
	"bufio"
	"bytes"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorldMap_ParseWorldMap(t *testing.T) {

	t.Run("wrong map line format", func(t *testing.T) {
		f := openTestdataFile(t, "bad_line_format_map.txt")
		defer f.Close()
		wm, err := ParseWorldMap(bufio.NewScanner(f))
		assert.Nil(t, wm)
		assert.EqualError(t, err, "Invalid map line: city with many names?")
	})

	t.Run("empty new lines should be skipped", func(t *testing.T) {
		f := openTestdataFile(t, "empty_newlines_map.txt")
		defer f.Close()
		wm, err := ParseWorldMap(bufio.NewScanner(f))
		assert.NoError(t, err)
		assert.Len(t, wm.cities, 9)
	})

	t.Run("inconsistent map error", func(t *testing.T) {
		f := openTestdataFile(t, "inconsistent_map.txt")
		defer f.Close()
		wm, err := ParseWorldMap(bufio.NewScanner(f))
		assert.Nil(t, wm)
		assert.EqualError(t, err, "Cannot parse the map: inconsistent map")
	})

	t.Run("success small map", func(t *testing.T) {
		f := openTestdataFile(t, "small_map.txt")
		defer f.Close()
		wm, err := ParseWorldMap(bufio.NewScanner(f))
		assert.NoError(t, err)
		assert.Len(t, wm.cities, 9)

		cs := []*city{
			nil,
			wm.cities["C1"], wm.cities["C2"], wm.cities["C3"],
			wm.cities["C4"], wm.cities["C5"], wm.cities["C6"],
			wm.cities["C7"], wm.cities["C8"], wm.cities["C9"],
		}

		assert.Equal(t, &city{
			name: "C1",
			dirs: [4]*city{dirNorth: nil, dirSouth: cs[4], dirEast: cs[2], dirWest: nil},
		}, cs[1])

		assert.Equal(t, &city{
			name: "C2",
			dirs: [4]*city{dirNorth: nil, dirSouth: cs[5], dirEast: cs[3], dirWest: cs[1]},
		}, cs[2])

		assert.Equal(t, &city{
			name: "C3",
			dirs: [4]*city{dirNorth: nil, dirSouth: cs[6], dirEast: nil, dirWest: cs[2]},
		}, cs[3])

		assert.Equal(t, &city{
			name: "C4",
			dirs: [4]*city{dirNorth: cs[1], dirSouth: cs[7], dirEast: cs[5], dirWest: nil},
		}, cs[4])

		assert.Equal(t, &city{
			name: "C5",
			dirs: [4]*city{dirNorth: cs[2], dirSouth: cs[8], dirEast: cs[6], dirWest: cs[4]},
		}, cs[5])

		assert.Equal(t, &city{
			name: "C6",
			dirs: [4]*city{dirNorth: cs[3], dirSouth: cs[9], dirEast: nil, dirWest: cs[5]},
		}, cs[6])

		assert.Equal(t, &city{
			name: "C7",
			dirs: [4]*city{dirNorth: cs[4], dirSouth: nil, dirEast: cs[8], dirWest: nil},
		}, cs[7])

		assert.Equal(t, &city{
			name: "C8",
			dirs: [4]*city{dirNorth: cs[5], dirSouth: nil, dirEast: cs[9], dirWest: cs[7]},
		}, cs[8])

		assert.Equal(t, &city{
			name: "C9",
			dirs: [4]*city{dirNorth: cs[6], dirSouth: nil, dirEast: nil, dirWest: cs[8]},
		}, cs[9])
	})

	t.Run("success big map", func(t *testing.T) {
		f := openTestdataFile(t, "big_map.txt")
		defer f.Close()
		wm, err := ParseWorldMap(bufio.NewScanner(f))
		assert.NoError(t, err)
		assert.Len(t, wm.cities, 120000)

		// C119433 north=C119033 south=C119833 east=C119434 west=C119432
		ctest := wm.cities["C119433"]
		assert.Equal(t, "C119433", ctest.name)
		assert.Equal(t, "C119033", ctest.dirs[dirNorth].name)
		assert.Equal(t, "C119833", ctest.dirs[dirSouth].name)
		assert.Equal(t, "C119434", ctest.dirs[dirEast].name)
		assert.Equal(t, "C119432", ctest.dirs[dirWest].name)

		// C119601 north=C119201 east=C119602
		ctest = wm.cities["C119601"]
		assert.Equal(t, "C119601", ctest.name)
		assert.Equal(t, "C119201", ctest.dirs[dirNorth].name)
		assert.Nil(t, ctest.dirs[dirSouth])
		assert.Equal(t, "C119602", ctest.dirs[dirEast].name)
		assert.Nil(t, ctest.dirs[dirWest])
	})
}

func TestWorldMap_getOrCreateCity(t *testing.T) {

	city1 := &city{name: "City1"}
	wm := &WorldMap{
		cities: map[string]*city{
			city1.name: city1,
		},
	}

	t.Run("empty city name", func(t *testing.T) {
		c := wm.getOrCreateCity("")
		assert.Nil(t, c)
	})

	t.Run("city does exist", func(t *testing.T) {
		c1 := wm.getOrCreateCity("City1")
		assert.Equal(t, city1, c1)
	})

	t.Run("city does not exist", func(t *testing.T) {
		c2 := wm.getOrCreateCity("City2")
		assert.NotNil(t, c2)
		ret, exist := wm.cities["City2"]
		assert.True(t, exist)
		assert.Equal(t, c2, ret)
	})

}

func TestWorldMap_destroyCity(t *testing.T) {

	t.Run("trying to destroy a non-existing city", func(t *testing.T) {
		wm := parseSmallMap(t)
		prev := make(map[string]*city, len(wm.cities))
		for k, v := range wm.cities {
			prev[k] = &*v
		}
		wm.destroyCity("not-exist")
		assert.Equal(t, prev, wm.cities)
	})

	t.Run("destroy an existing city that has all 4 surrounding cities", func(t *testing.T) {
		wm := parseSmallMap(t)

		assert.Len(t, wm.cities, 9)

		wm.destroyCity("C5")

		assert.Len(t, wm.cities, 8)

		_, exist := wm.cities["C5"]
		assert.False(t, exist)

		// surrounding cities should update their directions that were pointed
		// to C5.
		assert.Nil(t, wm.cities["C4"].dirs[dirEast])
		assert.Nil(t, wm.cities["C6"].dirs[dirWest])
		assert.Nil(t, wm.cities["C8"].dirs[dirNorth])
		assert.Nil(t, wm.cities["C2"].dirs[dirSouth])
	})

	t.Run("destroy an existing city that has 3 surrounding cities", func(t *testing.T) {
		wm := parseSmallMap(t)

		assert.Len(t, wm.cities, 9)

		wm.destroyCity("C4")

		assert.Len(t, wm.cities, 8)

		_, exist := wm.cities["C4"]
		assert.False(t, exist)

		// surrounding cities should update their directions that were pointed
		// to C4.
		assert.Nil(t, wm.cities["C1"].dirs[dirSouth])
		assert.Nil(t, wm.cities["C5"].dirs[dirWest])
		assert.Nil(t, wm.cities["C7"].dirs[dirNorth])
	})

	t.Run("destroy an existing city that has 2 surrounding cities", func(t *testing.T) {
		wm := parseSmallMap(t)

		assert.Len(t, wm.cities, 9)

		wm.destroyCity("C1")

		assert.Len(t, wm.cities, 8)

		_, exist := wm.cities["C1"]
		assert.False(t, exist)

		// surrounding cities should update their directions that were pointed
		// to C1.
		assert.Nil(t, wm.cities["C2"].dirs[dirWest])
		assert.Nil(t, wm.cities["C4"].dirs[dirNorth])
	})

}

func TestWorldMap_print(t *testing.T) {
	wm := parseSmallMap(t)
	wm.destroyCity("C5")
	wm.destroyCity("C2")
	wm.destroyCity("C7")
	wm.destroyCity("C3")
	wm.destroyCity("C1")

	t.Run("print non-empty map", func(t *testing.T) {
		buf := &bytes.Buffer{}
		wm.print(buf)
		assert.Equal(t, `C4
C6 south=C9
C8 east=C9
C9 north=C6 west=C8
`, buf.String())
	})

	t.Run("print empty map", func(t *testing.T) {
		wm.destroyCity("C4")
		wm.destroyCity("C6")
		wm.destroyCity("C8")
		wm.destroyCity("C9")

		buf := &bytes.Buffer{}
		wm.print(buf)
		assert.Equal(t, "No cities in the map.\n", buf.String())
	})
}

func TestWorldMap_decodeMapLine(t *testing.T) {
	t.Run("invalid format wrong number of spaces", func(t *testing.T) {
		d, err := decodeMapLine("My city north=N south=S east=E west=W")
		assert.Nil(t, d)
		assert.EqualError(t, err, "Invalid map line: invalid format - wrong number of spaces")
	})

	t.Run("invalid number of '=' characters", func(t *testing.T) {
		d, err := decodeMapLine("C1 north==N")
		assert.Nil(t, d)
		assert.EqualError(t, err, "Invalid map line: invalid number of '=' characters")
	})

	t.Run("invalid city with many names", func(t *testing.T) {
		d, err := decodeMapLine("C1 C1 north=N")
		assert.Nil(t, d)
		assert.EqualError(t, err, "Invalid map line: city with many names?")
	})

	t.Run("invalid direction", func(t *testing.T) {
		d, err := decodeMapLine("C1 nnorth=N")
		assert.Nil(t, d)
		assert.EqualError(t, err, "Invalid map line: nnorth is not a valid direction")
	})

	t.Run("multiple direction declaration", func(t *testing.T) {
		d, err := decodeMapLine("C1 north=N south=S south=S2")
		assert.Nil(t, d)
		assert.EqualError(t, err, "Invalid map line: multiple south directions")
	})

	t.Run("success with all directions", func(t *testing.T) {
		d, err := decodeMapLine("C1 west=C5 south=C3 east=C4 north=C2")
		assert.NoError(t, err)
		assert.NotNil(t, d)
		assert.Equal(t, "C1", d.name)
		assert.Equal(t, "C2", d.dirs[dirNorth])
		assert.Equal(t, "C3", d.dirs[dirSouth])
		assert.Equal(t, "C4", d.dirs[dirEast])
		assert.Equal(t, "C5", d.dirs[dirWest])
	})

	t.Run("success with incomplete directions", func(t *testing.T) {
		d, err := decodeMapLine("C1 east=C2 north=C3")
		assert.NoError(t, err)
		assert.Equal(t, "C1", d.name)
		assert.Equal(t, "C3", d.dirs[dirNorth])
		assert.Equal(t, "", d.dirs[dirSouth])
		assert.Equal(t, "C2", d.dirs[dirEast])
		assert.Equal(t, "", d.dirs[dirWest])
	})
}

// ===============================================================
// test utils
// ===============================================================

func openTestdataFile(t *testing.T, fname string) *os.File {
	f, err := os.Open(path.Join("testdata", fname))
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func parseSmallMap(t *testing.T) *WorldMap {
	f := openTestdataFile(t, "small_map.txt")
	defer f.Close()
	s := bufio.NewScanner(f)
	worldMap, err := ParseWorldMap(s)
	if err != nil {
		t.Fatal(err)
	}
	return worldMap
}

func parseNormalMap(t *testing.T) *WorldMap {
	f := openTestdataFile(t, "normal_map.txt")
	defer f.Close()
	s := bufio.NewScanner(f)
	worldMap, err := ParseWorldMap(s)
	if err != nil {
		t.Fatal(err)
	}
	return worldMap
}
