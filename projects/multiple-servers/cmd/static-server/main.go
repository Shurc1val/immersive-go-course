package main

import (
	"flag"
	"servers/static"
)

var filePath = flag.String("path", "default", "Path the static files are read from.")
var serverPort = flag.Int("port", 8082, "Port for the server to run on.")

func main() {
	flag.Parse()
	static.Run(*filePath, *serverPort)
}
