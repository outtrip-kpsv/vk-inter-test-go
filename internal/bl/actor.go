package bl

import (
	"vk-inter-test-go/internal/db/repo"
	"vk-inter-test-go/internal/io/models"
	"vk-inter-test-go/internal/utils"
)

func (b *BL) CreateActor(actor repo.Actor) (repo.Actor, error) {
	b.logger.Info("create actor")

	err := b.Db.Actor.CreateActor(&actor)
	if err != nil {
		return repo.Actor{}, err
	}
	return actor, nil
}

func (b *BL) DeleteActor(name string) (int64, error) {
	b.logger.Info("delete actor")

	res, err := b.Db.Actor.DeleteActorByName(name)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (b *BL) UpdateActor(actor repo.Actor) (repo.Actor, error) {
	b.logger.Info("update actor")
	dbActor, err := b.Db.Actor.GetActorById(actor.ID)
	if err != nil {
		return repo.Actor{}, err
	}
	if len(actor.Name) == 0 {
		actor.Name = dbActor.Name
	}
	if len(actor.BirthDateJson) == 0 {
		actor.BirthDateJson = dbActor.BirthDate.String()
		actor.BirthDate = dbActor.BirthDate
	}

	_, err = b.Db.Actor.UpdateActor(actor)
	if err != nil {
		return repo.Actor{}, err
	}
	actor, err = b.Db.Actor.GetActorById(actor.ID)
	if err != nil {
		return repo.Actor{}, err
	}
	return actor, nil
}

func (b *BL) GetAllActorsLikeName(name string, orderBy string) ([]models.ActorIo, error) {
	b.logger.Info("get actors like name")

	var actors []models.ActorIo
	allActors, err := b.Db.Actor.GetAllActorsLikeName(name, orderBy)
	if err != nil {
		return nil, err
	}
	var actorsIDs []int
	for _, actor := range allActors {
		actorsIDs = append(actorsIDs, actor.ID)
		actors = append(actors, models.ActorIo{Actor: actor})
	}

	actorIDsWithMovieIDs, err := b.Db.MovieActor.GetRelationByActorIDs(actorsIDs)
	if err != nil {
		return nil, err
	}
	movieIds := utils.UniqueValues(actorIDsWithMovieIDs)
	movieMap, err := b.Db.Movie.GetMovieMapByIDs(movieIds, "")
	if err != nil {
		return nil, err
	}

	for i, _ := range actors {
		for _, val := range actorIDsWithMovieIDs[actors[i].Actor.ID] {
			actors[i].Movies = append(actors[i].Movies, movieMap[val])
		}
	}
	return actors, nil
}
