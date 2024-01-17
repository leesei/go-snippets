package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var opts struct {
	// Slice of bool will append 'true' each time the option
	// is encountered (can be set multiple times, like -vvv)
	Verbosity []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	Verbose   int    // internal use only
	LogFile   string `short:"l" long:"log" description:"Log file (NDJSON format)"`

	Args struct {
		Input string `positional-arg-name:"INPUT" description:"Input file"`
	} `positional-args:"yes" required:"yes"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		return
	}
	opts.Verbose = len(opts.Verbosity)

	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.SetGlobalLevel(zerolog.InfoLevel - zerolog.Level(opts.Verbose))
	var f *os.File
	if opts.LogFile != "" {
		f, _ = os.Create(opts.LogFile)
		log.Logger = log.Output(f)
	} else {
		stdout, _ := IsRedirected()
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: stdout, TimeFormat: "15:04:05"})
	}

	log.Info().Interface("opts", opts).Send()

	stat, err := os.Stat(opts.Args.Input)
	if err != nil || stat.IsDir() {
		log.Error().Msg(fmt.Sprintf("Input %s is not a file", opts.Args.Input))
		return
	}

	log.Trace().Msg("trace message")
	log.Debug().Msg("debug message")
	log.Info().Msg("info message")
	log.Warn().Msg("warn message")
	log.Error().Msg("error message")
}
