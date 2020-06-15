package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/BurntSushi/toml"
)

// Requirements:
//
//  - Load a configuration file

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

var configPath = flag.String("config", "./mercury.toml", "Path to configuration")

type feed struct {
	Name string
	Feed string
}

type Config struct {
	Name    string
	URL     string `toml:url`
	Owner   string
	Email   string
	Cache   string
	Timeout duration
	Theme   string
	Output  string
	Feed    []feed
}

func main() {
	var config Config

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", path.Base(os.Args[0]))

		flag.PrintDefaults()
	}

	flag.Parse()
	if _, err := toml.DecodeFile(*configPath, &config); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", config)
}
