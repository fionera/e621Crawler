package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

var arguments = struct {
	StartId     int
	Concurrency int
	Output      string
	Verbose     bool
}{}

func parseArgs(args []string) {
	// Create new parser object
	parser := argparse.NewParser("e621Crawler", "")

	// Create flags
	startId := parser.Int("s", "start", &argparse.Options{
		Required: false,
		Default:  0,
		Help:     "ID from where to start crawling"})

	concurrency := parser.Int("c", "concurrency", &argparse.Options{
		Required: false,
		Help:     "Concurrency",
		Default:  4})

	output := parser.String("o", "output", &argparse.Options{
		Required: false,
		Default:  "Download",
		Help:     "Output folder"})

	verbose := parser.Flag("v", "verbose", &argparse.Options{
		Required: false,
		Default:  false,
	})

	// Parse input
	err := parser.Parse(args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		os.Exit(0)
	}

	// Fill arguments structure
	arguments.StartId = *startId
	arguments.Concurrency = *concurrency
	arguments.Output = *output
	arguments.Verbose = *verbose
}
