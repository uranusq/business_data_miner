package main

// Config ... Holds structure of TOML configuration file
type Config struct {
	General generalConfig
	Common  commonConfig
	Google  googleConfig
	Colly   collyConfig
}

type generalConfig struct {
	Database string
}

type commonConfig struct {
	Use   bool
	Path  string
	Debug bool
}

type googleConfig struct {
	Use   bool
	Path  string
	Debug bool
}

type collyConfig struct {
	Use   bool
	Path  string
	Debug bool
}
