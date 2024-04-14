package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)


var ErrorRetryTimeTooLong error = errors.New("retry time too long")
var ErrorConnectionDropped error = errors.New("connection dropped")
var ErrorServer error = errors.New("server error")
var ErrorReadingResponseBody error = errors.New("error reading response body")

type Client struct {
	ServerAddress string
	RetryTimeout int
}


func extractRetryDelay(timeString string, timeLimit int) (time.Duration, error) {
	delay := time.Duration(5e9)
	if retryTime, err := strconv.Atoi(timeString); err == nil {
		delay = time.Duration(retryTime)*time.Second
	} else if retryTime, err := time.Parse(http.TimeFormat, timeString); err == nil {
		delay = time.Until(retryTime)
	}
	if int(delay.Seconds()) > timeLimit {
		return time.Duration(0), ErrorRetryTimeTooLong
	}
	return delay, nil
}


func (c Client) FetchWeather() (string, error) {
	resp, err := http.Get(c.ServerAddress)
	if err != nil {
		return "", ErrorConnectionDropped
	}

	switch resp.StatusCode {
	case 429:
		// Server overloaded; wait to retry
		retryDelay, err := extractRetryDelay(resp.Header.Get("Retry-After"), c.RetryTimeout)
		if err != nil {
			return "", err
		}
		if retryDelay.Seconds() > 1 {
			fmt.Fprintf(os.Stderr, "Server overload; retrying request in %v seconds.\n", int(retryDelay)/1e9)
		}
		time.Sleep(retryDelay)
		return c.FetchWeather()
	case 500:
		// Server dropped connection
		return "", ErrorServer
	default:
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", ErrorReadingResponseBody
		}
		body := string(bodyBytes)
		return body, nil
	}
}