package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/boltdb/bolt"
)

var Movies []Movie

const API_URL = "https://gist.githubusercontent.com/Vaasu-Dhand/f2aa21c91c903413e26bb7ffb1389f6d/raw/a3e072bcf81c659485c69ad884f5b0342a385137/movies.json"
const MOVIE_BUCKET = "Movie"

// Create Scruct to represent Movie entity
type Movie struct {
	Id      string   `json:"id"`
	Title   string   `json:"title"`
	Year    int16    `json:"year"`
	Runtime int16    `json:"runtime"`
	Genre   []string `json:"genre"`
}

// Seed a local movies slice
func Seed(db *bolt.DB) {

	resp, err := http.Get(API_URL)
	if err != nil {
		fmt.Println("Error fetching seed data:", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	var movies []Movie

	json.Unmarshal(body, &movies)

	fmt.Println(movies)

	for _, movie := range movies {
		save(db, movie)
	}

	fmt.Println("🌱 Seeded")
}

// Store a movie in database
func save(db *bolt.DB, movie Movie) error {
	// Store the movie struct in the Movies bucket using the Id as the key.
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(MOVIE_BUCKET))
		if err != nil {
			return err
		}

		encoded, err := json.Marshal(movie)
		if err != nil {
			return err
		}
		return b.Put([]byte(movie.Id), encoded)
	})
	return err
}

// CRUD functions
func GetAll(db *bolt.DB) ([]Movie, error) {
	var movies []Movie
	err := db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte(MOVIE_BUCKET))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var movie Movie
			json.Unmarshal(v, &movie)
			movies = append(movies, movie)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error reading bucket:", err)
		return nil, fmt.Errorf("could not get movies")
	}
	return movies, nil

}

func GetById(db *bolt.DB, id string) (Movie, error) {

	var movie Movie
	err := db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte(MOVIE_BUCKET))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			if string(k) == id {
				json.Unmarshal(v, &movie)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Could not find movie with id: %s:", id)
		return Movie{}, fmt.Errorf("could not find movie")
	}
	return movie, nil

}

func CreateMovie(db *bolt.DB, movie Movie) (Movie, error) {

	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(MOVIE_BUCKET))

		numOfMovies := bucket.Stats().KeyN
		movie.Id = strconv.Itoa(numOfMovies + 1)

		encoded, err := json.Marshal(movie)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(movie.Id), encoded)
	})

	if err != nil {
		fmt.Println("Error reading bucket:", err)
		return Movie{}, fmt.Errorf("could not get movies")
	}

	return movie, nil
}

func UpdateMovie(db *bolt.DB, id string, body Movie) (Movie, error) {
	var updatedMovie Movie

	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(MOVIE_BUCKET))
		if bucket == nil {
			return fmt.Errorf("bucket %s does not exist", MOVIE_BUCKET)
		}

		val := bucket.Get([]byte(id))
		if val == nil {
			return fmt.Errorf("movie with id %s not found", id)
		}

		// Unmarshal existing movie
		if err := json.Unmarshal(val, &updatedMovie); err != nil {
			return fmt.Errorf("failed to decode movie: %v", err)
		}

		// Update fields if provided
		if body.Title != "" {
			updatedMovie.Title = body.Title
		}
		if len(body.Genre) != 0 {
			updatedMovie.Genre = body.Genre
		}
		if body.Runtime != 0 {
			updatedMovie.Runtime = body.Runtime
		}
		if body.Year != 0 {
			updatedMovie.Year = body.Year
		}

		// Marshal and save back to DB
		encoded, err := json.Marshal(updatedMovie)
		if err != nil {
			return fmt.Errorf("failed to encode movie: %v", err)
		}

		return bucket.Put([]byte(id), encoded)
	})

	if err != nil {
		return Movie{}, err
	}

	return updatedMovie, nil
}

func DeleteMovie(db *bolt.DB, id string) (bool, error) {

	err := db.Update(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte(MOVIE_BUCKET))

		keyToDelete := []byte(id)
		if err := bucket.Delete(keyToDelete); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Could not find delete with id: %s:", id)
		return false, fmt.Errorf("could not find movie")
	}
	return true, nil
}
