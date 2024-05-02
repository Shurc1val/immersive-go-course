package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
)

func Run(databaseURL string, port int) {
	dbConn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error connecting to the database.")
		os.Exit(1)
	}
	defer dbConn.Close(context.Background())

	http.HandleFunc("/images.json", func(w http.ResponseWriter, r *http.Request) {
		log.Println(port)
		MarshalIndentString := r.URL.Query().Get("indent")
		MarshalIndentNum := 0
		if MarshalIndentString != "" {
			MarshalIndentNum, err = strconv.Atoi(MarshalIndentString)
			if err != nil {
				http.Error(w, "Invalid indent query parameter value (must be a positive integer).", http.StatusBadRequest)
			}
		}

		var b []byte

		if r.Method == "GET" {
			images, err := fetchImages(dbConn)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error fetching images from the database.")
				os.Exit(1)
			}

			b, err = json.MarshalIndent(images, "", strings.Repeat(" ", MarshalIndentNum))
			if err != nil {
				http.Error(w, "Unable to marshal images data to json.", http.StatusInternalServerError)
			}

		} else if r.Method == "POST" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Unable to read request JSON data")
				os.Exit(1)
			}
			var image Image
			err = json.Unmarshal(body, &image)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Request JSON data invalid")
				os.Exit(1)
			}

			image, err = insertImage(dbConn, image)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Unable to add imaged to database")
				os.Exit(1)
			}

			b, err = json.MarshalIndent(image, "", strings.Repeat(" ", MarshalIndentNum))
			if err != nil {
				http.Error(w, "Unable to marshal image data to json.", http.StatusInternalServerError)
			}
		}

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Write(b)
	})

	http.ListenAndServe(":"+strconv.Itoa(port), nil)

}
