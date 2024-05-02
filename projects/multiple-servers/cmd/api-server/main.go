package main

import (
	"flag"
	"fmt"
	"os"
	"servers/api"
)

var serverPort = flag.Int("port", 8082, "Port for the server to run on.")

func main() {
	flag.Parse()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fmt.Fprintln(os.Stderr, "Required environment variable DATABASE_URL not set")
		os.Exit(1)
	}
	api.Run(dbURL, *serverPort)
}
