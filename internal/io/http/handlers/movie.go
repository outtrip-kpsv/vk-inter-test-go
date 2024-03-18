package handlers

import (
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
	"vk-inter-test-go/internal/db/repo"
	"vk-inter-test-go/internal/io/ioutils"
	"vk-inter-test-go/internal/io/models"
)

// CreateMovie создает новый фильм.
//
// @Summary Создает новый фильм
// @Description Создает новый фильм с данными, предоставленными в теле запроса.
// @Tags Movies
// @Accept  json
// @Produce  json
// @Param body body models.MovieIo true "Данные фильма"
// @Param Authorization header string true "Bearer"
// @Security bearerAuth
// @Success 200 {object} models.MovieIo "Успешно созданный фильм"
// @Failure 400 {object} models.ErrorResponse "Неверный формат данных"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/movie [post]
func (c *Controller) CreateMovie(w http.ResponseWriter, req *http.Request) {

	var movie models.MovieIo
	err := ioutils.DecodeRequestBody(req, &movie)
	if err != nil || !ioutils.MovieJsonValidate(&movie) {
		ioutils.HandleInvalidJson(w)
		return
	}

	var answer interface{}

	movieIo, err := c.Bl.CreateMovie(movie)
	if err != nil {
		c.logger.Info("err", zap.Error(err))
		answer = models.ErrorResponse{
			Error: "err : '" + err.Error() + "'",
		}
	} else {
		answer = movieIo
	}
	c.logger.Info("resp : ", zap.Reflect("answer :", answer))

	ioutils.RespJson(w, answer)
}

// DeleteMovie удаляет фильм по его ID.
//
// @Summary Удаляет фильм по его ID
// @Description Удаляет фильм с указанным ID.
// @Tags Movies
// @Param id query integer true "ID фильма для удаления"
// @Param Authorization header string true "Bearer"
// @Security bearerAuth
// @Success 200 {object} models.OkResponse "Успешное удаление фильма"
// @Failure 400 {object} models.ErrorResponse "Неверное значение ID или ошибка удаления"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/movie [delete]
func (c *Controller) DeleteMovie(w http.ResponseWriter, req *http.Request) {
	idStr := req.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ioutils.RespErrorText("Не верное значение ID", w)
		return
	}
	rows, err := c.Bl.DeleteMovie(id)
	if err != nil || rows == 0 {
		ioutils.RespErrorText("Ошибка удалеия", w)
		return
	}
	answer := models.OkResponse{Ok: "Запись удалена"}
	c.logger.Info("resp : ", zap.Reflect("answer :", answer))

	ioutils.RespJson(w, answer)
}

// UpdateMovie обновляет данные фильма.
//
// @Summary Обновляет данные фильма
// @Description Обновляет данные фильма с данными, предоставленными в теле запроса.
// @Tags Movies
// @Accept  json
// @Produce  json
// @Param body body repo.Movie true "Данные фильма для обновления"
// @Param Authorization header string true "Bearer"
// @Security bearerAuth
// @Success 200 {object} repo.Movie "Обновленные данные фильма"
// @Failure 400 {object} models.ErrorResponse "Неверный формат данных"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/movie [put]
func (c *Controller) UpdateMovie(w http.ResponseWriter, req *http.Request) {
	var movie repo.Movie
	err := ioutils.DecodeRequestBody(req, &movie)

	if len(movie.ReleaseDateJson) > 0 {
		movie.ReleaseDate, err = time.Parse("2006-01-02", movie.ReleaseDateJson)
		if err != nil {
			c.logger.Info("err :", zap.Error(err))
			ioutils.HandleInvalidJson(w)
			return
		}
	}

	if movie.Rating < 0 || movie.Rating > 10 {
		ioutils.HandleInvalidJson(w)
		return
	}

	var answer interface{}
	movie, err = c.Bl.UpdateMovie(movie)
	if err != nil {
		c.logger.Info("err", zap.Error(err))
		answer = models.ErrorResponse{
			Error: "err : '" + err.Error() + "'",
		}
	} else {
		movie.ReleaseDateJson = movie.ReleaseDate.Format("2006-01-02")
		answer = movie
	}
	c.logger.Info("resp : ", zap.Reflect("answer :", answer))

	ioutils.RespJson(w, answer)
}

// GetMovies получает все фильмы или фильмы с определенным заголовком или именем актера.
//
// @Summary Получает все фильмы или фильмы с определенным заголовком или именем актера
// @Description Получает все фильмы, если ни один из параметров не указан, или фильмы с определенным заголовком или именем актера.
// @Tags Movies
// @Param title query string false "Заголовок фильма для фильтрации"
// @Param name query string false "Имя актера для фильтрации"
// @Param sort query string false "Поле для сортировки (например, 'title')" Доступные значения: 'rating', 'title', 'date'
// @Param Authorization header string true "Bearer"
// @Security bearerAuth
// @Accept  json
// @Produce  json
// @Success 200 {array} models.MovieIo "Список фильмов"
// @Failure 400 {object} models.ErrorResponse "Неверный формат данных"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/movie [get]
func (c *Controller) GetMovies(w http.ResponseWriter, req *http.Request) {
	title := req.URL.Query().Get("title")
	name := req.URL.Query().Get("name")
	orderBy := req.URL.Query().Get("sort")

	if len(name) != 0 {
		movieIo, err := c.Bl.GetAllMoviesByNameActor(name, orderBy)
		if err != nil {
			ioutils.RespErrorText(err.Error(), w)
			return
		}
		ioutils.RespJson(w, movieIo)
		return
	}
	movieIo, err := c.Bl.GetAllMoviesByTitle(title, orderBy)
	if err != nil {
		ioutils.RespErrorText(err.Error(), w)
		return
	}
	c.logger.Info("resp : ", zap.Reflect("answer :", movieIo))

	ioutils.RespJson(w, movieIo)
	return

}
