package handlers

import (
	"go.uber.org/zap"
	"net/http"
	"time"
	"vk-inter-test-go/internal/db/repo"
	"vk-inter-test-go/internal/io/ioutils"
	"vk-inter-test-go/internal/io/models"
)

// CreateActor создает нового актера.
//
// @Summary Создает нового актера
// @Description Создает нового актера с данными, предоставленными в теле запроса.
// @Tags Actors
// @Accept  json
// @Produce  json
// @Param body body repo.Actor true "Данные актера"
// @Param Authorization header string true "Bearer"
// @Security bearerAuth
// @Success 200 {object} repo.Actor "Успешно созданный актер"
// @Failure 400 {object} models.ErrorResponse "Неверный формат данных"
// @Failure 401 {object} models.ErrorResponse "Отказано в доступе: ошибка извлечения токена"
// @Failure 403 {object} models.ErrorResponse "Доступ запрещен: отсутствие необходимой роли"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/actor [post]
func (c *Controller) CreateActor(w http.ResponseWriter, req *http.Request) {

	var actor repo.Actor
	err := ioutils.DecodeRequestBody(req, &actor)
	if err != nil || !ioutils.ActorJsonValidate(&actor) {
		ioutils.HandleInvalidJson(w)
		return
	}

	var answer interface{}

	createActor, err := c.Bl.CreateActor(actor)
	if err != nil {
		c.logger.Info("err", zap.Error(err))
		answer = models.ErrorResponse{
			Error: "err : '" + err.Error() + "'",
		}
	} else {
		answer = createActor
	}
	c.logger.Info("resp : ", zap.Reflect("answer :", createActor))
	ioutils.RespJson(w, answer)
}

// DeleteActor удаляет актера по имени.
//
// @Summary Удаляет актера по имени
// @Description Удаляет актера с указанным именем.
// @Tags Actors
// @Param name query string true "Имя актера для удаления"
// @Param Authorization header string true "Bearer"
// @Security bearerAuth
// @Success 200 {object} models.OkResponse "Успешное удаление актера"
// @Failure 400 {object} models.ErrorResponse "Неверный формат данных"
// @Failure 403 {object} models.ErrorResponse "Доступ запрещен: отсутствие необходимой роли"
// @Failure 404 {object} models.ErrorResponse "Актер не найден"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/actor [delete]
func (c *Controller) DeleteActor(w http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")

	var answer interface{}

	res, err := c.Bl.DeleteActor(name)
	if res == 0 {
		answer = models.ErrorResponse{
			Error: "актера с таким именем нет в базе",
		}
		ioutils.RespJson(w, answer)
		return
	}
	if err != nil {
		c.logger.Info("err", zap.Error(err))
		answer = models.ErrorResponse{
			Error: "err : '" + err.Error() + "'",
		}
	} else {
		answer = models.OkResponse{Ok: "Актер " + name + " удален"}
	}
	c.logger.Info("resp : ", zap.Reflect("answer :", answer))

	ioutils.RespJson(w, answer)
}

// UpdateActor обновляет данные актера.
//
// @Summary Обновляет данные актера
// @Description Обновляет данные актера с данными, предоставленными в теле запроса.
// @Tags Actors
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer"
// @Security bearerAuth
// @Param body body repo.Actor true "Данные актера для обновления"
// @Success 200 {object} repo.Actor "Обновленные данные актера"
// @Failure 400 {object} models.ErrorResponse "Неверный формат данных"
// @Failure 403 {object} models.ErrorResponse "Доступ запрещен: отсутствие необходимой роли"
// @Failure 404 {object} models.ErrorResponse "Актер не найден"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/actor [put]
func (c *Controller) UpdateActor(w http.ResponseWriter, req *http.Request) {
	var actor repo.Actor
	err := ioutils.DecodeRequestBody(req, &actor)

	if len(actor.BirthDateJson) > 0 {
		actor.BirthDate, err = time.Parse("2006-01-02", actor.BirthDateJson)
		if err != nil {
			c.logger.Info("err :", zap.Error(err))
			ioutils.HandleInvalidJson(w)
			return
		}
	}
	var answer interface{}
	if len(actor.Gender) > 0 && !(actor.Gender == "male" || actor.Gender == "female") {
		ioutils.HandleInvalidJson(w)
		return
	}
	actor, err = c.Bl.UpdateActor(actor)
	if err != nil {
		c.logger.Info("err", zap.Error(err))
		answer = models.ErrorResponse{
			Error: "err : '" + err.Error() + "'",
		}
	} else {
		actor.BirthDateJson = actor.BirthDate.Format("2006-01-02")
		answer = actor
	}
	c.logger.Info("resp : ", zap.Reflect("answer :", answer))

	ioutils.RespJson(w, answer)
}

// GetAllActors получает всех актеров.
//
// @Summary Получает всех актеров или актеров с определенным именем
// @Description Получает всех актеров, если имя не указано, или актеров с определенным именем, если имя указано в запросе.
// @Tags Actors
// @Param name query string false "Имя актера для фильтрации"
// @Param sort query string false "Поле для сортировки (например, 'name')"
// @Param Authorization header string true "Bearer"
// @Security bearerAuth
// @Accept  json
// @Produce  json
// @Success 200 {array} models.ActorIo "Список актеров"
// @Failure 400 {object} models.ErrorResponse "Неверный формат данных"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/actor [get]
func (c *Controller) GetAllActors(w http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")
	orderBy := req.URL.Query().Get("sort")

	actors, err := c.Bl.GetAllActorsLikeName(name, orderBy)
	if err != nil {
		return
	}
	c.logger.Info("resp : ", zap.Reflect("answer :", actors))

	ioutils.RespJson(w, actors)
}
