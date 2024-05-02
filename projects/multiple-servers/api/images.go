package api

import (
	"context"

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
