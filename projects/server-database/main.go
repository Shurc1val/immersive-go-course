package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
)

type Image struct {
	Title   string `json:"title"`
	AltText string `json:"alt_text"`
	URL     string `json:"url"`
}

func fetchImages(dbConn *pgx.Conn) ([]Image, error) {
	var images []Image

	rows, err := dbConn.Query(context.Background(), "SELECT title, url, alt_text FROM public.images")
	if err != nil {
		return images, err
	}

	for rows.Next() {
		var title, url, altText string
		err = rows.Scan(&title, &url, &altText)
		if err != nil {
			return images, err
		}
		images = append(images, Image{Title: title, URL: url, AltText: altText})
	}

	return images, nil
}

func insertImage(dbConn *pgx.Conn, image Image) (Image, error) {
	rows, err := dbConn.Query(context.Background(), "INSERT INTO public.images(title, url, alt_text) VALUES ($1, $2, $3)", image.Title, image.URL, image.AltText)
	if err != nil {
		return Image{}, err
	}
	for rows.Next() {
		err = rows.Scan(&image.Title, &image.URL, &image.AltText)
	}
	if err != nil {
		return Image{}, err
	}
	return image, nil
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fmt.Fprintln(os.Stderr, "Required environment variable DATABASE_URL not set")
		os.Exit(1)
	}

	dbConn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer dbConn.Close(context.Background())

	http.HandleFunc("/images.json", func(w http.ResponseWriter, r *http.Request) {
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
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			b, err = json.MarshalIndent(images, "", strings.Repeat(" ", MarshalIndentNum))
			if err != nil {
				http.Error(w, "Unable to marshal images data to json.", http.StatusInternalServerError)
			}

		} else if r.Method == "POST" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			var image Image
			err = json.Unmarshal(body, &image)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			image, err = insertImage(dbConn, image)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			b, err = json.MarshalIndent(image, "", strings.Repeat(" ", MarshalIndentNum))
			if err != nil {
				http.Error(w, "Unable to marshal image data to json.", http.StatusInternalServerError)
			}
		}

		w.Header().Add("Content-Type", "application/json")

		w.Write(b)
	})

	http.ListenAndServe(":8080", nil)

}
