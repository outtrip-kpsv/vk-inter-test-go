package ioutils

import (
	"time"
	"vk-inter-test-go/internal/db/repo"
	"vk-inter-test-go/internal/io/models"
)

func UserJsonValidate(user repo.User) bool {
	if len(user.Login) > 0 && len(user.Pass) > 0 {
		return true
	}
	return false
}

func ActorJsonValidate(actor *repo.Actor) bool {
	if len(actor.Name) == 0 {
		return false
	}
	if !(actor.Gender == "male" || actor.Gender == "female") {
		return false
	}
	var err error

	actor.BirthDate, err = time.Parse("2006-01-02", actor.BirthDateJson)
	if err != nil {
		return false
	}
	return true
}

func MovieJsonValidate(movie *models.MovieIo) bool {
	if len(movie.Movie.Title) < 1 || len(movie.Movie.Title) > 150 {
		return false
	}
	if len(movie.Movie.Description) > 1000 {
		return false
	}
	if movie.Movie.Rating > 10 || movie.Movie.Rating < 0 {
		return false
	}
	var err error
	movie.Movie.ReleaseDate, err = time.Parse("2006-01-02", movie.Movie.ReleaseDateJson)
	if err != nil {
		return false
	}
	for i, _ := range movie.Actors {
		if !ActorJsonValidate(&movie.Actors[i]) {
			return false
		}
	}
	return true
}
