# Alien Invasion Simulator

This project simulates a world alien invasion.

# Table of contents

1. [Invasion rules](#1-invasion-rules)
2. [Map file format](#2-map-file-format)
3. [Project structure](#3-project-structure)
4. [Executables](#4-executables)
5. [Tests](#5-tests)

## 1. Invasion rules

1. N ​aliens​ ​start​ ​out​ ​at​ ​random​ ​places (cities) ​on​ ​the​ ​map,​ ​and​ ​wander​ ​around​ ​randomly following​ ​links.​ ​Each​ ​iteration,​ ​the​ ​aliens​ ​can​ ​travel​ ​in​ ​any​ ​of​ ​the​ ​directions leading​ ​out​ ​of​ ​a​ ​city.​
2. If two or more aliens end up in​ ​the​ ​same​ ​place,​ ​they​ ​fight,​ ​and​ ​in​ ​the​ ​process​ ​kill each​ ​other​ ​and​ ​destroy​ ​the​ ​city.​
3. Once​ ​a​ ​city​ ​is​ ​destroyed,​ ​aliens​ ​can​ ​no​ ​longer​ ​travel​ ​to​ ​or​ ​through​ ​it. ​This may​ ​lead​ ​to​ ​aliens​ ​getting​ ​"trapped".
4. If at the starting point, there is more than 1 alien in a same city, they will not fight until the first move/iteration has been done.
5. The​ ​program​ ​will ​run​ ​until​ ​all​ ​the​ ​aliens or the cities​ ​have​ ​been destroyed,​ ​or​ ​each​ ​non-trapped alien​ ​has​ ​moved​ ​at​ ​least​ ​10,000​ ​times.

The simulation will finish in the following identified cases:

- **Case 1:** All the aliens have been destroyed.
- **Case 2:** Each non-trapped alien has moved 10,000 times.
- **Case 3:** Remaining aliens cannot reach each other.
- **Case 4:** All the remaining aliens were trapped.
- **Case 5:** There is only 1 free alien, then no city can be destroyed.

## 2. Map file format

Map file format should be '.txt' file containing a city with its surroundings in each line. For example:

```
C1 south=C4 east=C2
C2 south=C5 east=C3 west=C1
C3 south=C6 west=C2
C4 north=C1 south=C7 east=C5
C5 north=C2 south=C8 east=C6 west=C4
C6 north=C3 south=C9 west=C5
C7 north=C4 east=C8
C8 north=C5 east=C9 west=C7
C9 north=C6 west=C8
```

**Rules:**

1. City names must have only letters [A-Z a-z] or numbers [0-9] or both. Spaces and other symbols are not allowed.
2. Directions must be any of the lowered-cases words: `north | south | east | west`. The directions order does not matter.

## 3. Project structure

This project is organized in 3 main folders:

1. **cmd/**: Contains all the executable main packages of this project: `map_generator` and `simulator`.
2. **invasion/**: Contains all the business logic about the invasion simulator.
3. **mapgen/**: Contains all the business logic to generate world maps with a given width and height.

## 4. Executables

### 4.1. Map generator (cmd/map_generator)

Map generator easily creates a map with a given width and height specified by the user in the passed output file.

```
$ go run cmd/map_generator/main.go -h
Usage of map_generator:
  -height int
        The height of the map. (default 20)
  -out string
        Output file where the generated map will be written. Ignoring this, the generated map will be printed in STDOUT.
  -width int
        The width of the map. (default 20)
```

For example, if we want to generate a map file in the current directory called `my_map.txt` with a width of 20 and a height of 30, then you should write:

```
$ go run cmd/map_generator/main.go -out my_map.txt -width 20 -height 30
```

### 4.2. Simulator (cmd/simulator)

Simulator command line starts invasion simulator with a given number of aliens, map file and output file specified by the user.

```
$ go run cmd/simulator/main.go -h
Usage of simulator:
  -m string
        Specify the world map file used for the invasion. (default "invasion/testdata/small_map.txt")
  -n int
        Specify the number of aliens for the invasion. (default 10)
  -o string
        Output file where the simulator result will be written. Ignoring this, the result will be redirected to STDOUT.
```

For example, if we want to simulate the invasion of 1000 aliens using a map file called `map1.txt` and save the result in `result.txt`, you should write:

```
$ go run cmd/invasion/main.go -n 1000 -m map1.txt -o result.txt
```

## 5. Tests

The available test files (`*_test.go`) can be found inside `invasion/` and `mapgen/` folders.

To run all the project tests:

`go test ./...`

To run only `invasion` related tests:

`go test ./invasion/...`

To run only `mapgen` related tests:

`go test ./mapgen/...`

**\*Note:** You can add `-count 1` flag to ensure the tests run again by skipping cached results.
