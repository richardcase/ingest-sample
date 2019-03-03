package main

// Options represents the runtime options for the service
type Options struct {
	ConfigFile    string
	ListenAddress string
	LogLevel      string
	DbURL         string
	DbName        string
	CollName      string
}
