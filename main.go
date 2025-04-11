package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var ErrInvalidReleased = fmt.Errorf("invalid release")

func main() {
	for {
		fmt.Println("REQUEST")
		r, err := GetRandomRelease(generateRandomId())
		if err != nil {
			if errors.Is(err, ErrInvalidReleased) {
				continue
			} else {
				log.Fatal(err)
			}
		}
		fmt.Printf("%+v", r)
		break
	}
}

type Artist struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"resource_url"`
}

type Image struct {
	Type string `json:"type"`
	Uri  string `json:"uri"`
}

type Track struct {
	Duration string `json:"duration"`
	Position string `json:"position"`
	Title    string `json:"title"`
	Type     string `json:"type_"`
}

type Video struct {
	Uri string `json:"uri"`
}

type Release struct {
	Id        int      `json:"id"`
	Title     string   `json:"title"`
	Artists   []Artist `json:"artists"`
	Year      int      `json:"year"`
	Genres    []string `json:"genres"`
	Country   string   `json:"country"`
	Images    []Image  `json:"images"`
	Tracklist []Track  `json:"tracklist"`
	Uri       string   `json:"uri"`
	Videos    []Video  `json:"videos"`
}

const (
	releaseEndpoint = "https://api.discogs.com/releases/%s"
	userAgent       = "HitMeApp/0.1 +development_mode"
)

func GetRandomRelease(id string) (Release, error) {
	var r Release

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		fmt.Sprintf(releaseEndpoint, id),
		nil,
	)
	if err != nil {
		panic(err)
	}

	req.Header.Add("User-Agent", userAgent)

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return Release{}, ErrInvalidReleased
		} else {
			return Release{}, fmt.Errorf("error: %v", res.Status)
		}
	}

	err = json.NewDecoder(res.Body).Decode(&r)

	if err != nil {
		return Release{}, err
	}

	if len(r.Videos) == 0 {
		return Release{}, ErrInvalidReleased
	}

	return r, nil
}

func generateRandomId() string {
	max := 9_999_999
	return strconv.Itoa(rand.Intn(max))
}
