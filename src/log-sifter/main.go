package main

import (
	"bufio"
	"fmt"
	"os"

	levenshtein "github.com/agnivade/levenshtein"
	flags "github.com/jessevdk/go-flags"
	util "github.com/woanware/goutil"
)

// ##### Structs ##############################################################

type Options struct {
	Input    string  `short:"i" long:"input" description:"Input file" required:"true"`
	Distance float64 `short:"d" long:"distance" description:"Levenshtein distance" default:"0.30" required:"false"`
}

// ##### Variables ############################################################

var (
	options Options
)

// ##### Constants ############################################################

const APP_NAME string = "log-sifter"
const APP_VERSION string = "0.0.1"

// ##### Methods ##############################################################

// main is the application entry point!
func main() {

	//fmt.Println(fmt.Sprintf("\n%s v%s - woanware\n", APP_NAME, APP_VERSION))

	parseCommandLine()

	if util.DoesFileExist(options.Input) == false {
		fmt.Println("Input file (-i) does not exist")
		os.Exit(-1)
	}

	file, err := os.Open(options.Input)
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		os.Exit(-1)
	}
	defer file.Close()

	data := make(map[string]struct{})
	line := ""
	distance := 0
	d := ""
	display := false
	normalisedLev := 0.0
	lineLen := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()

		display = true
		lineLen = len(line)
		for d = range data {
			distance = levenshtein.ComputeDistance(d, line)
			normalisedLev = float64(distance) / (float64(lineLen + len(d)))

			if normalisedLev <= options.Distance {
				display = false
				break
			}
		}

		if display == true {
			fmt.Println(line)
			data[line] = struct{}{}
		}
	}
}

// parseCommandLine parses the command line options
func parseCommandLine() {

	var parser = flags.NewParser(&options, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}
