package invasion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirections_fromString(t *testing.T) {
	t.Run("invalid direction", func(t *testing.T) {
		_, err := directionFromString("invalid")
		assert.Error(t, err)
		assert.Equal(t, "invalid is not a valid direction", err.Error())

		_, err = directionFromString("North")
		assert.Error(t, err)
		assert.Equal(t, "North is not a valid direction", err.Error())
	})

	t.Run("valid direction", func(t *testing.T) {
		d, err := directionFromString("north")
		assert.NoError(t, err)
		assert.Equal(t, dirNorth, d)

		d, err = directionFromString("south")
		assert.NoError(t, err)
		assert.Equal(t, dirSouth, d)

		d, err = directionFromString("east")
		assert.NoError(t, err)
		assert.Equal(t, dirEast, d)

		d, err = directionFromString("west")
		assert.NoError(t, err)
		assert.Equal(t, dirWest, d)
	})

}

func TestDirections_opposite(t *testing.T) {
	t.Run("invalid direction", func(t *testing.T) {
		assert.Equal(t, direction(-1), direction(5).opposite())
	})
	t.Run("valid direction", func(t *testing.T) {
		assert.Equal(t, dirSouth, dirNorth.opposite())
		assert.Equal(t, dirNorth, dirSouth.opposite())
		assert.Equal(t, dirEast, dirWest.opposite())
		assert.Equal(t, dirWest, dirEast.opposite())
	})
}

func TestDirections_String(t *testing.T) {
	assert.Equal(t, "invalid direction", direction(8).String())
	assert.Equal(t, "north", dirNorth.String())
	assert.Equal(t, "south", dirSouth.String())
	assert.Equal(t, "east", dirEast.String())
	assert.Equal(t, "west", dirWest.String())
}
