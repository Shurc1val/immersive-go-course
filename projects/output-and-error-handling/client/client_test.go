package client

import (
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetchWeather(t *testing.T) {
	t.Run("Status code 200 response", func(t *testing.T) {
		serverResponseBody := `The weather is bloody awful!`
		test_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("HTTP request had method %v, not \"GET\", as expected.", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(serverResponseBody))
		}))
		defer test_server.Close()
		test_client := Client{ServerAddress: test_server.URL}
		out, err := test_client.FetchWeather()
		if err != nil {
			t.Error("Unexpected error on request: ", err)
		}
		if out != serverResponseBody {
			t.Errorf("Output, %q, did not match expected, %q.", out, serverResponseBody)
		}
	})

	t.Run("Status code 500 response", func(t *testing.T) {
		serverResponseBody := `The weather is bloody awful!`
		test_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("HTTP request had method %v, not \"GET\", as expected.", r.Method)
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(serverResponseBody))
		}))
		defer test_server.Close()
		test_client := Client{ServerAddress: test_server.URL}
		_, err := test_client.FetchWeather()
		if err != ErrorServer {
			t.Error("Unexpected error on request: ", err)
		}
	})

	t.Run("Connection dropped", func(t *testing.T) {
		test_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("HTTP request had method %v, not \"GET\", as expected.", r.Method)
			}
			conn, _, _ := w.(http.Hijacker).Hijack()
			conn.Close()
		}))
		defer test_server.Close()
		test_client := Client{ServerAddress: test_server.URL}
		_, err := test_client.FetchWeather()
		if err != ErrorConnectionDropped {
			t.Error("Unexpected error on request: ", err)
		}
	})
}

func TestExtractRetryDelay(t *testing.T) {
	timeout := 5
	tests := map[string]struct {
		timeString string
		expectedOutput time.Duration
		expectedError error
	}{
		"Integer time string < timeout": {
			timeString: "4",
			expectedOutput: time.Duration(4*time.Second),
			expectedError: nil,
		},
		"Integer time string > timeout": {
			timeString: "7",
			expectedOutput: time.Duration(0),
			expectedError: ErrorRetryTimeTooLong,
		},
		"Datetime time string < timeout": {
			timeString: time.Now().UTC().Add(time.Duration(3) * time.Second).Format(http.TimeFormat),
			expectedOutput: time.Duration(3*time.Second),
			expectedError: nil,
		},
		"Datetime time string > timeout": {
			timeString: time.Now().UTC().Add(time.Duration(12) * time.Second).Format(http.TimeFormat),
			expectedOutput: time.Duration(0),
			expectedError: ErrorRetryTimeTooLong,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := extractRetryDelay(test.timeString, timeout)
			if err != test.expectedError {
				t.Errorf("Error, %q, does not match expected error %q", err.Error(), test.expectedError.Error())
			}
			// Rounded because time passes between setting the variables and running the test
			if int(math.Round(output.Seconds())) != int(test.expectedOutput.Seconds()) {
				t.Errorf("Output, %v, does not match expected output %v", output, test.expectedOutput)
			}
		})
	}
}