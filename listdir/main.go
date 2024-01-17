package main

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var opts struct {
	Args struct {
		Input string `positional-arg-name:"INPUT" description:"Input folder to scan"`
	} `positional-args:"yes" required:"yes"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		return
	}

	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	stat, err := os.Stat(opts.Args.Input)
	if err != nil || !stat.IsDir() {
		log.Error().Msg(fmt.Sprintf("Input %s is not a directory", opts.Args.Input))
		return
	}

	items, err := ListDirectory(opts.Args.Input, func(fileInfo fs.FileInfo) bool {
		log.Info().Str("path", opts.Args.Input).Str("name", fileInfo.Name()).Bool("dir", fileInfo.IsDir()).Msg("ListDirectory() func")
		return true
	})
	if err != nil {
		log.Error().Err(err).Msg("ListDirectory() error")
	} else {
		log.Info().Interface("items", items).Msg("ListDirectory() success")
	}

	items = []string{}
	err = fs.WalkDir(os.DirFS(opts.Args.Input), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		log.Info().Str("path", path).Str("name", d.Name()).Bool("dir", d.IsDir()).Msg("WalkDir() func")
		if path == "." {
			return nil
		}
		if d.IsDir() {
			path = path + "/"
		}
		items = append(items, path)
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("WalkDir() error")
	} else {
		log.Info().Interface("items", items).Msg("WalkDir() success")
	}

}
