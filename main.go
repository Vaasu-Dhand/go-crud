package main

/*
 * BONUS
 * 1. Connect to a database
 */

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Seed Movies
	Seed()

	// Serve them with Gorilla Mux
	router := mux.NewRouter()
	router.HandleFunc("/movies", MoviesHandler)
	router.HandleFunc("/movies/{id:[0-9]+}", MovieHandler)
	fmt.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func MoviesHandler(res http.ResponseWriter, req *http.Request) {
	// vars := mux.Vars(req)
	res.Header().Set("Content-Type", "application/json")

	fmt.Println(req.Method + " " + req.URL.Path)

	switch req.Method {
	// getAll()
	case "GET":
		response, _ := json.Marshal(GetAll())
		res.WriteHeader(http.StatusOK)
		res.Write(response)
	// createMovie()
	case "POST":
		body, _ := io.ReadAll(req.Body)

		var movie Movie
		json.Unmarshal(body, &movie)
		fmt.Println(movie)
		newMovie, err := CreateMovie(movie)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("Could not add movie at this time. Please try again!"))
			return
		}
		res.WriteHeader(http.StatusCreated)
		response, _ := json.Marshal(newMovie)
		res.Write(response)
	}
}

func MovieHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	fmt.Println(req.Method + " " + req.URL.Path)

	vars := mux.Vars(req)
	id := vars["id"]

	switch req.Method {
	// getById()
	case "GET":
		movie, err := GetById(id)
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			res.Write([]byte(fmt.Sprintf("Movie with Id: %v not found", id)))
			return
		}
		res.WriteHeader(http.StatusOK)
		response, _ := json.Marshal(movie)
		res.Write(response)
	// updateMovie
	case "PUT":
		// Get body
		body, _ := io.ReadAll(req.Body)
		var movie Movie
		json.Unmarshal(body, &movie)
		// Pass it on to the Crud function for update

		newMovie, err := UpdateMovie(id, movie)
		if err != nil {
			fmt.Print(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("Could not update movie. Please try again!"))
			return
		}
		res.WriteHeader(http.StatusOK)
		response, _ := json.Marshal(newMovie)
		res.Write(response)
	// deleteMovie
	case "DELETE":
		// Extract the id, call the crud func, return deleted status
		deleted, err := DeleteMovie(id)
		if err != nil || !deleted {
			fmt.Print((err))
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("Could not delete movie. Please try again!"))
			return
		}
		res.WriteHeader(200)
		// Ananoymous Struct
		response, _ := json.Marshal(struct {
			Deleted bool
		}{
			Deleted: true,
		})
		res.Write(response)
	}

}
