package bl

import (
	"go.uber.org/zap"
	"time"
	"vk-inter-test-go/internal/db/repo"
	"vk-inter-test-go/internal/io/models"
	"vk-inter-test-go/internal/utils"
)

func (b *BL) CreateMovie(movie models.MovieIo) (models.MovieIo, error) {
	b.logger.Info("create movie")

	err := b.Db.Movie.CreateMovie(&movie.Movie)
	if err != nil {
		return models.MovieIo{}, err
	}
	actorIDs := make([]int, 0)
	for i, actor := range movie.Actors {
		dbActor, err := b.Db.Actor.GetActorByName(actor.Name)
		if err != nil {
			err := b.Db.Actor.CreateActor(&movie.Actors[i])
			if err != nil {
				b.logger.Info("err :", zap.Error(err))
			}
			actorIDs = append(actorIDs, movie.Actors[i].ID)
			continue
		}
		actorIDs = append(actorIDs, dbActor.ID)
		movie.Actors[i].ID = dbActor.ID
	}

	err = b.Db.MovieActor.CreateMovieActorRelation(movie.Movie.ID, actorIDs)
	if err != nil {
		return models.MovieIo{}, err
	}
	return movie, nil

}

func (b *BL) DeleteMovie(id int) (int64, error) {
	b.logger.Info("delete movie")

	rows, err := b.Db.Movie.DeleteMovieById(id)
	if err != nil {
		return 0, err
	}
	return rows, nil
}

func (b *BL) UpdateMovie(movie repo.Movie) (repo.Movie, error) {
	dbMovie, err := b.Db.Movie.GetMovieById(movie.ID)
	if err != nil {
		return repo.Movie{}, err
	}

	if len(movie.Title) == 0 {
		movie.Title = dbMovie.Title
	}

	if len(movie.Description) == 0 {
		movie.Description = dbMovie.Description
	}

	if len(movie.ReleaseDateJson) == 0 {
		movie.ReleaseDateJson = dbMovie.ReleaseDate.Format("2006-01-02")
		movie.ReleaseDate = dbMovie.ReleaseDate
	} else {
		movie.ReleaseDate, err = time.Parse("2006-01-02", movie.ReleaseDateJson)
		if err != nil {
			return repo.Movie{}, err
		}
	}

	if movie.Rating == 0 {
		movie.Rating = dbMovie.Rating
	}

	_, err = b.Db.Movie.UpdateMovie(movie)
	if err != nil {
		return repo.Movie{}, err
	}

	movie, err = b.Db.Movie.GetMovieById(movie.ID)
	if err != nil {
		return repo.Movie{}, err
	}

	return movie, nil
}

func (b *BL) GetAllMoviesByTitle(title string, orderBy string) ([]models.MovieIo, error) {
	b.logger.Info("get movies by title")

	allMovies, err := b.Db.Movie.GetMoviesLikeTitle(title, orderBy)
	if err != nil {
		return nil, err
	}

	var movies []models.MovieIo
	var movieIDs []int

	for _, movie := range allMovies {
		movies = append(movies, models.MovieIo{Movie: movie})
		movieIDs = append(movieIDs, movie.ID)
	}

	movieIDsWithActorIDs, err := b.Db.MovieActor.GetRelationByMovieIDs(movieIDs)
	if err != nil {
		return nil, err
	}

	actorIDs := utils.UniqueValues(movieIDsWithActorIDs)

	actorMap, err := b.Db.Actor.GetActorMapByIDs(actorIDs)
	if err != nil {
		return nil, err
	}

	for i, _ := range movies {
		for _, actorID := range movieIDsWithActorIDs[movies[i].Movie.ID] {
			movies[i].Actors = append(movies[i].Actors, actorMap[actorID]...)
		}
	}

	return movies, nil
}

func (b *BL) GetAllMoviesByNameActor(name string, orderBy string) ([]models.MovieIo, error) {
	b.logger.Info("get movies by name actor")

	allActor, err := b.Db.Actor.GetAllActorsLikeName(name, "")
	if err != nil {
		return nil, err
	}

	var movies []models.MovieIo
	var actorsIDs []int

	for _, actor := range allActor {
		actorsIDs = append(actorsIDs, actor.ID)
	}

	actorIDsWithMovieIDs, err := b.Db.MovieActor.GetRelationByActorIDs(actorsIDs)
	if err != nil {
		return nil, err
	}

	movieIDs := utils.UniqueValues(actorIDsWithMovieIDs)
	movieIDsWithActorIDs, err := b.Db.MovieActor.GetRelationByMovieIDs(movieIDs)
	if err != nil {
		return nil, err
	}

	movieMap, err := b.Db.Movie.GetMovieMapByIDs(movieIDs, orderBy)
	actorMap, err := b.Db.Actor.GetActorMapByIDs(actorsIDs)
	if err != nil {
		return nil, err
	}
	var i int
	for k, v := range movieMap {
		movies = append(movies, models.MovieIo{Movie: v})
		for _, actorID := range movieIDsWithActorIDs[k] {
			movies[i].Actors = append(movies[i].Actors, actorMap[actorID]...)
		}
		i++
	}

	return movies, nil
}
