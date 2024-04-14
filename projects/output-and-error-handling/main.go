package main

import (
	"fmt"
	"os"

	"github.com/CodeYourFuture/immersive-go-course/projects/output-and-error-handling/client"
)

const serverURL string = "http://localhost:8080"

func main() {
	weatherClient := client.Client{ServerAddress: serverURL}
	weather, err := weatherClient.FetchWeather()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to fetch weather:", err.Error())
		os.Exit(1)
	}
	fmt.Println(weather)
	os.Exit(0)
}