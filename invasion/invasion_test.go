package invasion

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvasion_Start(t *testing.T) {

	t.Run("aliens not destroyed at starter city", func(t *testing.T) {
		wm := parseSmallMap(t)

		moves := map[int]*moveList{
			0: newMoveList(wm, "C4", "C5", "C6", "C5"),
			1: newMoveList(wm, "C6", "C3", "C2"),
			2: newMoveList(wm, "C8", "C9", "C8", "C5"),
			3: newMoveList(wm, "C4", "C1", "C2"),
			4: newMoveList(wm, "C4", "C7", "C4", "C5"),
		}

		tStarterCityNameForAlienFn = func(alien int) string {
			mvs, ok := moves[alien]
			assert.True(t, ok, "unexpected alien number")
			return mvs.poll().name
		}
		tAlienMoveNextCityFn = func(alien int, curCity *city, possibleCities []*city) *city {
			mvs, ok := moves[alien]
			require.True(t, ok)
			ret := mvs.poll()
			assert.Contains(t, possibleCities, ret)
			assert.Contains(t, possibleCities, ret, fmt.Sprintf("alien %d can't go to %s from %s", alien, ret.name, curCity.name))
			return ret
		}

		buf := &bytes.Buffer{}
		Start(buf, wm, len(moves))
		ret := buf.String()

		assert.Contains(t, ret, "C2 has been destroyed by alien 1 and alien 3!")
		assert.NotContains(t, ret, "C4 has been destroyed by alien 0, alien 3 and alien 4!")
		_, exist := wm.cities["C4"]
		assert.True(t, exist)
	})

	// case 1
	t.Run("all aliens have been destroyed", func(t *testing.T) {
		wm := parseSmallMap(t)

		moves := map[int]*moveList{
			0: newMoveList(wm, "C4", "C5", "C6", "C5"),
			1: newMoveList(wm, "C6", "C3", "C2"),
			2: newMoveList(wm, "C8", "C9", "C8", "C5"),
			3: newMoveList(wm, "C4", "C1", "C2"),
			4: newMoveList(wm, "C4", "C7", "C4", "C5"),
		}

		tStarterCityNameForAlienFn = func(alien int) string {
			mvs, ok := moves[alien]
			assert.True(t, ok, "unexpected alien number")
			return mvs.poll().name
		}
		tAlienMoveNextCityFn = func(alien int, curCity *city, possibleCities []*city) *city {
			mvs, ok := moves[alien]
			require.True(t, ok)
			ret := mvs.poll()
			assert.Contains(t, possibleCities, ret, fmt.Sprintf("alien %d can't go to %s from %s", alien, ret.name, curCity.name))
			return ret
		}

		buf := &bytes.Buffer{}
		Start(buf, wm, len(moves))

		assertFileResult(t, "result_case_1.txt", buf.String())
	})

	// case 2
	t.Run("each non-trapped alien has moved 10000 times", func(t *testing.T) {
		wm := parseNormalMap(t)

		moves := map[int]*moveList{
			3: newMoveList(wm, "C4", "C3", "C2", "C1"),
			1: newMoveList(wm, "C13", "C7"),
			5: newMoveList(wm, "C13", "C7"),
			2: newMoveList(wm, "C20", "C14", "C8", "C2"),
			6: newMoveList(wm, "C15", "C9", "C3", "C2"),
			4: newRepeatedMoveList(wm, 100001, "C11", "C17"),
			0: newRepeatedMoveList(wm, 100001, "C28", "C29"),
		}

		tStarterCityNameForAlienFn = func(alien int) string {
			mvs, ok := moves[alien]
			assert.True(t, ok, "unexpected alien number")
			return mvs.poll().name
		}
		tAlienMoveNextCityFn = func(alien int, curCity *city, possibleCities []*city) *city {
			mvs, ok := moves[alien]
			require.True(t, ok)
			ret := mvs.poll()
			assert.Contains(t, possibleCities, ret, fmt.Sprintf("alien %d can't go to %s from %s", alien, ret.name, curCity.name))
			return ret
		}

		buf := &bytes.Buffer{}
		Start(buf, wm, len(moves))

		assertFileResult(t, "result_case_2.txt", buf.String())
	})

	// case 3
	t.Run("remaining aliens cannot reach each other", func(t *testing.T) {
		wm := parseNormalMap(t)

		moves := map[int]*moveList{
			0:  newMoveList(wm, "C3", "C4"),
			1:  newMoveList(wm, "C5", "C4"),
			2:  newMoveList(wm, "C8", "C9", "C10"),
			3:  newMoveList(wm, "C12", "C11", "C10"),
			4:  newMoveList(wm, "C13", "C14", "C15", "C16"),
			5:  newMoveList(wm, "C20", "C21", "C22", "C16"),
			6:  newMoveList(wm, "C19", "C20", "C21", "C22", "C23"),
			7:  newMoveList(wm, "C26", "C27", "C28", "C29", "C23"),
			8:  newMoveList(wm, "C30", "C24", "C30", "C24", "C30", "C24"),
			9:  newMoveList(wm, "C18", "C12", "C18", "C12", "C18", "C24"),
			10: newRepeatedMoveList(wm, 10001, "C6", "C5"),
			11: newRepeatedMoveList(wm, 10001, "C8", "C7"),
		}

		tStarterCityNameForAlienFn = func(alien int) string {
			mvs, ok := moves[alien]
			assert.True(t, ok, "unexpected alien number")
			return mvs.poll().name
		}
		tAlienMoveNextCityFn = func(alien int, curCity *city, possibleCities []*city) *city {
			mvs, ok := moves[alien]
			require.True(t, ok)
			ret := mvs.poll()
			require.NotNil(t, ret, fmt.Sprintf("alien %d cannot have nil moved city", alien))
			assert.Contains(t, possibleCities, ret, fmt.Sprintf("alien %d can't go to %s from %s", alien, ret.name, curCity.name))
			return ret
		}

		buf := &bytes.Buffer{}
		Start(buf, wm, len(moves))

		assertFileResult(t, "result_case_3.txt", buf.String())
	})

	// case 4
	t.Run("all aliens were trapped", func(t *testing.T) {
		wm := parseNormalMap(t)

		moves := map[int]*moveList{
			0: newMoveList(wm, "C3", "C2"),
			1: newMoveList(wm, "C8", "C2"),
			2: newMoveList(wm, "C19", "C13", "C7"),
			3: newMoveList(wm, "C14", "C8", "C7"),
			4: newMoveList(wm, "C13", "C7", "C1"),
			5: newMoveList(wm, "C26", "C27", "C28", "C29"),
			6: newMoveList(wm, "C21", "C22", "C23", "C29"),
			7: newMoveList(wm, "C5", "C6", "C12", "C18", "C24"),
			8: newMoveList(wm, "C20", "C21", "C22", "C23", "C24"),
			9: newMoveList(wm, "C30", "C24", "C30", "C24", "C30"),
		}

		tStarterCityNameForAlienFn = func(alien int) string {
			mvs, ok := moves[alien]
			assert.True(t, ok, "unexpected alien number")
			return mvs.poll().name
		}
		tAlienMoveNextCityFn = func(alien int, curCity *city, possibleCities []*city) *city {
			mvs, ok := moves[alien]
			require.True(t, ok)
			ret := mvs.poll()
			require.NotNil(t, ret, fmt.Sprintf("alien %d cannot have nil moved city", alien))
			assert.Contains(t, possibleCities, ret, fmt.Sprintf("alien %d can't go to %s from %s", alien, ret.name, curCity.name))
			return ret
		}

		buf := &bytes.Buffer{}
		Start(buf, wm, len(moves))

		assertFileResult(t, "result_case_4.txt", buf.String())
	})

	// case 5
	t.Run("only 1 alien free", func(t *testing.T) {
		wm := parseNormalMap(t)

		moves := map[int]*moveList{
			0: newMoveList(wm, "C3", "C2"),
			1: newMoveList(wm, "C8", "C2"),
			2: newMoveList(wm, "C19", "C13", "C7"),
			3: newMoveList(wm, "C14", "C8", "C7"),
			4: newRepeatedMoveList(wm, 100, "C22", "C23"),
		}

		tStarterCityNameForAlienFn = func(alien int) string {
			mvs, ok := moves[alien]
			assert.True(t, ok, "unexpected alien number")
			return mvs.poll().name
		}
		tAlienMoveNextCityFn = func(alien int, curCity *city, possibleCities []*city) *city {
			mvs, ok := moves[alien]
			require.True(t, ok)
			ret := mvs.poll()
			require.NotNil(t, ret, fmt.Sprintf("alien %d cannot have nil moved city", alien))
			assert.Contains(t, possibleCities, ret, fmt.Sprintf("alien %d can't go to %s from %s", alien, ret.name, curCity.name))
			return ret
		}

		buf := &bytes.Buffer{}
		Start(buf, wm, len(moves))

		assertFileResult(t, "result_case_5.txt", buf.String())
	})

	t.Run("all cities were destroyed", func(t *testing.T) {
		wm := parseSmallMap(t)

		moves := map[int]*moveList{
			0:  newMoveList(wm, "C2", "C1"),
			1:  newMoveList(wm, "C4", "C1"),
			2:  newMoveList(wm, "C3", "C2"),
			3:  newMoveList(wm, "C5", "C2"),
			4:  newMoveList(wm, "C2", "C3"),
			5:  newMoveList(wm, "C2", "C3"),
			6:  newMoveList(wm, "C7", "C4"),
			7:  newMoveList(wm, "C5", "C4"),
			8:  newMoveList(wm, "C6", "C5"),
			9:  newMoveList(wm, "C8", "C5"),
			10: newMoveList(wm, "C9", "C6"),
			11: newMoveList(wm, "C9", "C6"),
			12: newMoveList(wm, "C8", "C7"),
			17: newMoveList(wm, "C8", "C7"),
			14: newMoveList(wm, "C9", "C8"),
			15: newMoveList(wm, "C9", "C8"),
			16: newMoveList(wm, "C8", "C9"),
			13: newMoveList(wm, "C8", "C9"),
		}

		tStarterCityNameForAlienFn = func(alien int) string {
			mvs, ok := moves[alien]
			assert.True(t, ok, "unexpected alien number")
			return mvs.poll().name
		}
		tAlienMoveNextCityFn = func(alien int, curCity *city, possibleCities []*city) *city {
			mvs, ok := moves[alien]
			require.True(t, ok)
			ret := mvs.poll()
			require.NotNil(t, ret, fmt.Sprintf("alien %d cannot have nil moved city", alien))
			assert.Contains(t, possibleCities, ret, fmt.Sprintf("alien %d can't go to %s from %s", alien, ret.name, curCity.name))
			return ret
		}

		buf := &bytes.Buffer{}
		Start(buf, wm, len(moves))
		ret := buf.String()
		assert.Contains(t, ret, "C1 has been destroyed by alien 0 and alien 1!")
		assert.Contains(t, ret, "C2 has been destroyed by alien 2 and alien 3!")
		assert.Contains(t, ret, "C3 has been destroyed by alien 4 and alien 5!")
		assert.Contains(t, ret, "C4 has been destroyed by alien 6 and alien 7!")
		assert.Contains(t, ret, "C5 has been destroyed by alien 8 and alien 9!")
		assert.Contains(t, ret, "C6 has been destroyed by alien 10 and alien 11!")
		assert.Contains(t, ret, "C7 has been destroyed by alien 12 and alien 17!")
		assert.Contains(t, ret, "C8 has been destroyed by alien 14 and alien 15!")
		assert.Contains(t, ret, "C9 has been destroyed by alien 13 and alien 16!")
		assert.Contains(t, ret, "All the aliens have been destroyed!")
		assert.Contains(t, ret, "No cities in the map.")
	})

	t.Run("passing nil map or 0 num of aliens", func(t *testing.T) {
		// don't do anything if nil map or 0 numOfAliens are passed
		buf := &bytes.Buffer{}
		Start(buf, nil, 2)
		assert.Empty(t, buf.Bytes())

		buf.Reset()
		wm := parseSmallMap(t)
		Start(buf, wm, 0)
		assert.Empty(t, buf.Bytes())
	})
}

func TestInvasion_remAliensCanReachEachOther(t *testing.T) {
	wm := parseSmallMap(t)
	alienMap := map[int]*alien{
		2: {num: 2, curCity: wm.cities["C4"]},
		4: {num: 4, curCity: wm.cities["C9"]},
		5: {num: 5, trapped: true}, // trapped aliens are ignored
	}

	t.Run("remaining aliens can reach each other", func(t *testing.T) {
		assert.True(t, remAliensCanReachEachOther(alienMap))
	})

	t.Run("remaining aliens cannot reach each other", func(t *testing.T) {
		// if C3, C5 and C8 are destroyed, then it is not possible to go from C4 to C9
		wm.destroyCity("C3")
		wm.destroyCity("C5")
		wm.destroyCity("C8")
		assert.False(t, remAliensCanReachEachOther(alienMap))
	})
}

func TestInvasion_addAlienToCity(t *testing.T) {

	m := make(map[string]*alienSet)

	addAlienToCity(m, "C1", 23)
	addAlienToCity(m, "C1", 12)
	addAlienToCity(m, "C2", 45)
	addAlienToCity(m, "C1", 10)

	s1 := newAlienSet()
	s1.add(23)
	s1.add(12)
	s1.add(10)
	s2 := newAlienSet()
	s2.add(45)
	assert.Equal(t, map[string]*alienSet{
		"C1": s1,
		"C2": s2,
	}, m)

}

func TestInvasion_removeAlienFromCity(t *testing.T) {
	s1 := newAlienSet()
	s1.add(3)
	s1.add(5)
	s2 := newAlienSet()
	s2.add(10)
	m := map[string]*alienSet{
		"C1": s1,
		"C2": s2,
		"C3": nil,
	}
	mSnapshot := make(map[string]*alienSet, len(m))
	for k, v := range m {
		if v != nil {
			mSnapshot[k] = &*v
		} else {
			mSnapshot[k] = nil
		}
	}

	t.Run("remove from city that does not exist", func(t *testing.T) {
		// as "not-exist" city does not exist, this won't do anything
		removeAlienFromCity(m, "not-exist", 3)
		assert.Equal(t, mSnapshot, m)
	})

	t.Run("remove alien that is not in the specified city", func(t *testing.T) {
		// as the alien 10 is not in C1, this won't do anything
		removeAlienFromCity(m, "C1", 10)
		assert.Equal(t, mSnapshot, m)
	})

	t.Run("removing from a city with nil alien set", func(t *testing.T) {
		// this should not panic
		assert.NotPanics(t, func() {
			removeAlienFromCity(m, "C3", 10)
		})
		assert.Equal(t, mSnapshot, m)
	})

	t.Run("successful remove", func(t *testing.T) {
		removeAlienFromCity(m, "C1", 3)
		updSet1 := newAlienSet()
		updSet1.add(5)
		assert.Equal(t, map[string]*alienSet{
			"C1": updSet1,
			"C2": s2,
			"C3": nil,
		}, m)
		removeAlienFromCity(m, "C1", 5)
		assert.Equal(t, map[string]*alienSet{
			"C1": newAlienSet(),
			"C2": s2,
			"C3": nil,
		}, m)
		removeAlienFromCity(m, "C2", 10)
		assert.Equal(t, map[string]*alienSet{
			"C1": newAlienSet(),
			"C2": newAlienSet(),
			"C3": nil,
		}, m)
	})
}

func TestInvasion_alienSet(t *testing.T) {

	aset := newAlienSet()
	assert.NotNil(t, aset.data)

	aset.add(8)
	aset.add(1)
	aset.add(3)

	assert.Equal(t, map[int]struct{}{
		1: {}, 3: {}, 8: {},
	}, aset.data)
	assert.Equal(t, "alien 1, alien 3 and alien 8", aset.String())
	assert.Equal(t, 3, aset.len())

	aset.remove(1)

	assert.Equal(t, map[int]struct{}{
		3: {}, 8: {},
	}, aset.data)
	assert.Equal(t, "alien 3 and alien 8", aset.String())
	assert.Equal(t, 2, aset.len())

	aset.remove(3)

	assert.Equal(t, map[int]struct{}{
		8: {},
	}, aset.data)
	assert.Equal(t, "alien 8", aset.String())
	assert.Equal(t, 1, aset.len())

	aset.remove(8)

	assert.Equal(t, map[int]struct{}{}, aset.data)
	assert.Equal(t, "", aset.String())
	assert.Equal(t, 0, aset.len())

}

// ===============================================================
// test utils
// ===============================================================

func assertFileResult(t *testing.T, fname string, got string) {
	data, err := os.ReadFile(path.Join("testdata", fname))
	require.NoError(t, err)
	assert.Equal(t, string(data), got)
}

type moveList struct {
	cs []*city
}

func newMoveList(wm *WorldMap, cityNames ...string) *moveList {
	cs := make([]*city, 0, len(cityNames))
	for _, n := range cityNames {
		cs = append(cs, wm.cities[n])
	}
	return &moveList{cs: cs}
}

func newRepeatedMoveList(wm *WorldMap, movesLen int, c1, c2 string) *moveList {
	cs := make([]*city, 0, movesLen)
	for i := 0; i < movesLen; i++ {
		if i%2 == 0 {
			cs = append(cs, wm.cities[c1])
		} else {
			cs = append(cs, wm.cities[c2])
		}
	}
	return &moveList{cs: cs}
}

func (l *moveList) poll() *city {
	if len(l.cs) == 0 {
		return nil
	}
	head := l.cs[0]
	l.cs = l.cs[1:]
	return head
}
