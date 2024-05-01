package main

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

func isAuthorised(username string, password string) bool {
	if (username != os.Getenv("AUTH_USERNAME")) || (password != os.Getenv("AUTH_PASSWORD")) {
		return false
	}
	return true
}

func main() {
	godotenv.Load()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Add("Content-Type", "text/html")
			params, err := url.ParseQuery(r.URL.RawQuery)
			if err != nil {
				w.Write([]byte("500 - Internal Server Error"))
			} else {
				w.Write([]byte(`<!DOCTYPE html>
<html>
<p>Query parameters:
<ul>`))
				for param, value := range params {
					user_input := html.EscapeString(fmt.Sprint(param+":", value))
					output := fmt.Sprintln("<li>", user_input, "</li>")
					w.Write([]byte(output))
				}
				w.Write([]byte(`</ul>
</html>`))
			}
		case "POST":
			rb := r.Body
			data, err := io.ReadAll(rb)
			if err != nil {
				w.Write([]byte("500 - Internal Server Error"))
			} else {
				w.Header().Add("Content-Type", "text/html")
				w.Write([]byte(`<!DOCTYPE html>
<html>
`))
				w.Write([]byte(html.EscapeString(string(data[:]))))
			}
		}
	})

	http.HandleFunc("/200", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("200 - OK"))
	})

	http.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
	})

	http.Handle("/404", http.NotFoundHandler())

	http.HandleFunc("/authenticated", func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Add("WWW-Authenticate", `Basic realm="localhost", charset="UTF-8"`)
			w.WriteHeader(http.StatusUnauthorized)
		}

		if !isAuthorised(username, password) {
			w.Header().Add("WWW-Authenticate", `Basic realm="localhost", charset="UTF-8"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid credentials."))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("'sup"))
	})

	http.HandleFunc("/limited", func(w http.ResponseWriter, r *http.Request) {
		limiter := rate.NewLimiter(100, 30)

		if limiter.Allow() {
			w.Write([]byte("Hello"))
		} else {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
		}
	})

	http.ListenAndServe(":8080", nil)
}
