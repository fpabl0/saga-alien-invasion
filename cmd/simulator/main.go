package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fpabl0/saga-alien-invasion/invasion"
)

func main() {

	var (
		numOfAliens int
		mapFile     string
		outputFile  string
	)

	flag.IntVar(&numOfAliens, "n", 10, "Specify the number of aliens for the invasion.")
	flag.StringVar(&mapFile, "m", "invasion/testdata/small_map.txt", "Specify the world map file used for the invasion.")
	flag.StringVar(&outputFile, "o", "", "Output file where the simulator result will be written. Ignoring this, the result will be redirected to STDOUT.")
	flag.Parse()

	if numOfAliens <= 0 {
		log.Fatalln("The number of aliens must be greater than 0")
	}

	var out io.Writer
	if outputFile == "" {
		out = os.Stdout
	} else {
		buf := &bytes.Buffer{}
		buf.Grow(1000)
		out = buf
	}

	worldMap, err := parseWorldMapFile(mapFile)
	if err != nil {
		log.Fatalln(err)
	}

	invasion.Start(out, worldMap, numOfAliens)

	if buf, ok := out.(*bytes.Buffer); ok {
		if err := os.WriteFile(outputFile, buf.Bytes(), os.ModePerm); err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("File %q was created successfully.\n", outputFile)
	}
}

func parseWorldMapFile(fname string) (*invasion.WorldMap, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	worldMap, err := invasion.ParseWorldMap(s)
	if err != nil {
		return nil, err
	}
	return worldMap, nil
}
