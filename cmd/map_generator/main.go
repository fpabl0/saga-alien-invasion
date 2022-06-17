package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fpabl0/saga-alien-invasion/mapgen"
)

func main() {
	var (
		outputFile string
		width      int
		height     int
	)

	flag.StringVar(&outputFile, "out", "", "Output file where the generated map will be written. Ignoring this, the generated map will be printed in STDOUT.")
	flag.IntVar(&width, "width", 20, "The width of the map.")
	flag.IntVar(&height, "height", 20, "The height of the map.")
	flag.Parse()

	if outputFile != "" {
		// check if the file already exists
		_, err := os.Stat(outputFile)
		if err == nil {
			log.Fatalf("The file %q already exists.", outputFile)
		}
	}

	g := mapgen.NewGenerator(width, height)
	data := g.Generate()

	if outputFile == "" {
		fmt.Println(string(data))
	} else {
		if err := os.WriteFile(outputFile, data, os.ModePerm); err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("Created %s file\n", outputFile)
	}

}
