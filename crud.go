package main

import (
	"fmt"
	"strconv"
)

var Movies []Movie

// Create Scruct to represent Movie entity
type Movie struct {
	Id      string   `json:"id"`
	Title   string   `json:"title"`
	Year    int16    `json:"year"`
	Runtime int16    `json:"runtime"`
	Genre   []string `json:"genre"`
}

// Seed a local movies slice
func Seed() {
	zodiacMovie := Movie{Id: "1", Title: "Zodiac", Year: 2007, Runtime: 137, Genre: []string{"suspense", "thriller"}}
	prionersMovie := Movie{Id: "2", Title: "Prisoners", Year: 2013, Runtime: 153, Genre: []string{"crime", "drama"}}
	theKillerMovie := Movie{Id: "3", Title: "The Killer", Year: 2023, Runtime: 118, Genre: []string{"crime", "action"}}
	Movies = append(Movies, zodiacMovie, prionersMovie, theKillerMovie)
}

// CRUD functions
func GetAll() []Movie {
	return Movies
}

func GetById(id string) (Movie, error) {
	for _, movie := range Movies {
		if movie.Id == id {
			return movie, nil
		}
	}
	return Movie{}, fmt.Errorf("movie with id %s not found", id)
}

func CreateMovie(movie Movie) (Movie, error) {
	movie.Id = strconv.Itoa(len(Movies) + 1)

	// Append to Movies slice
	Movies = append(Movies, movie)

	return movie, nil
}

func UpdateMovie(id string, body Movie) (Movie, error) {
	for i, movie := range Movies {
		if movie.Id == id {
			if body.Title != "" {
				Movies[i].Title = body.Title
			}
			if len(body.Genre) != 0 {
				Movies[i].Genre = body.Genre
			}
			if body.Runtime != 0 {
				Movies[i].Runtime = body.Runtime
			}
			if body.Year != 0 {
				Movies[i].Year = body.Year
			}
			return Movies[i], nil
		}
	}
	return Movie{}, fmt.Errorf("movie with id %s not found", id)
}

func DeleteMovie(id string) (bool, error) {
	// Loop over array
	for i, movie := range Movies {
		if movie.Id == id {
			Movies = append(Movies[:i], Movies[i+1:]...)
			return true, nil
		}
	}
	return false, fmt.Errorf("Could not delete movie with id %s", id)

}
