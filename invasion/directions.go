package invasion

import "fmt"

// direction defines direction type (N, S, E, O)
type direction int

// direction options
const (
	dirNorth direction = iota
	dirSouth
	dirEast
	dirWest
)

// directionFromString converts a valid string direction
// into a direction type. If the string is not a valid direction
// this will return an error.
func directionFromString(s string) (direction, error) {
	switch s {
	case "north":
		return dirNorth, nil
	case "south":
		return dirSouth, nil
	case "east":
		return dirEast, nil
	case "west":
		return dirWest, nil
	}
	return -1, fmt.Errorf("%s is not a valid direction", s)
}

// opposite returns the opposite direction of this direction.
// For example calling opposite() with dirNorth should return
// dirSouth.
func (d direction) opposite() direction {
	switch d {
	case dirNorth:
		return dirSouth
	case dirSouth:
		return dirNorth
	case dirEast:
		return dirWest
	case dirWest:
		return dirEast
	}
	return -1
}

// String implements fmt.Stringer. Will return the equivalent direction
// string in lowercase letters.
func (d direction) String() string {
	switch d {
	case dirNorth:
		return "north"
	case dirSouth:
		return "south"
	case dirEast:
		return "east"
	case dirWest:
		return "west"
	}
	return "invalid direction"
}
