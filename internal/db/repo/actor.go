package repo

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

type ActorRepositoryImpl struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewActorRepository(db *pgxpool.Pool, logger *zap.Logger) *ActorRepositoryImpl {
	logger.Info("create")
	return &ActorRepositoryImpl{db: db, logger: logger}
}

type Actor struct {
	ID            int       `db:"id" json:"ID"`
	Name          string    `db:"name" json:"name,omitempty"`
	Gender        string    `db:"gender" json:"gender,omitempty"`
	BirthDateJson string    `db:"-" json:"birthDate,omitempty"`
	BirthDate     time.Time `db:"birth-date" json:"-"`
}

type ActorRepository interface {
	CreateActor(actor *Actor) error
	DeleteActorByName(name string) (int64, error)
	UpdateActor(actor Actor) (int64, error)
	GetActorById(id int) (Actor, error)
	GetActorByName(name string) (Actor, error)
	GetAllActorsLikeName(name string, orderBy string) ([]Actor, error)
	GetActorMapByIDs(actorIDs []int) (map[int][]Actor, error)
}

func (a ActorRepositoryImpl) CreateActor(actor *Actor) error {
	sql := "INSERT INTO actors (name, gender, birth_date) VALUES ($1, $2, $3) RETURNING id"
	err := a.db.QueryRow(context.Background(), sql, actor.Name, actor.Gender, actor.BirthDate).Scan(&actor.ID)
	if err != nil {
		return err
	}
	return nil
}

func (a ActorRepositoryImpl) DeleteActorByName(name string) (int64, error) {
	sql := "DELETE FROM actors WHERE name = $1"
	res, err := a.db.Exec(context.Background(), sql, name)

	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func (a ActorRepositoryImpl) UpdateActor(actor Actor) (int64, error) {
	sql := "UPDATE actors SET name = $2, gender = $3, birth_date = $4 WHERE id = $1"
	res, err := a.db.Exec(context.Background(), sql, actor.ID, actor.Name, actor.Gender, actor.BirthDate)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func (a ActorRepositoryImpl) GetActorById(id int) (Actor, error) {
	var actor Actor

	sql := "SELECT id, name, gender, birth_date FROM actors WHERE id = $1"
	err := a.db.QueryRow(context.Background(), sql, id).Scan(&actor.ID, &actor.Name, &actor.Gender, &actor.BirthDate)
	if err != nil {
		return Actor{}, err
	}
	actor.BirthDateJson = actor.BirthDate.Format("2006-01-02")
	return actor, nil
}

func (a ActorRepositoryImpl) GetActorByName(name string) (Actor, error) {
	var actor Actor

	sql := "SELECT id, name, gender, birth_date FROM actors WHERE name = $1"
	err := a.db.QueryRow(context.Background(), sql, name).Scan(&actor.ID, &actor.Name, &actor.Gender, &actor.BirthDate)
	if err != nil {
		return Actor{}, err
	}
	actor.BirthDateJson = actor.BirthDate.Format("2006-01-02")
	return actor, nil
}

func (a ActorRepositoryImpl) GetAllActorsLikeName(name string, orderBy string) ([]Actor, error) {
	var actors []Actor
	sql := "SELECT id, name, gender, birth_date FROM actors WHERE name LIKE '%' || $1 || '%' ORDER BY "
	switch orderBy {
	case "name":
		sql += "name"
	case "date":
		sql += "birth_date DESC"
	default:
		sql += "id"
	}
	rows, err := a.db.Query(context.Background(), sql, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var actor Actor
		if err := rows.Scan(&actor.ID, &actor.Name, &actor.Gender, &actor.BirthDate); err != nil {
			return nil, err
		}
		actor.BirthDateJson = actor.BirthDate.Format("2006-01-02")
		actors = append(actors, actor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return actors, nil
}

func (a ActorRepositoryImpl) GetActorMapByIDs(actorIDs []int) (map[int][]Actor, error) {
	actorMap := make(map[int][]Actor)

	sql := "SELECT id, name, gender, birth_date FROM actors WHERE id = ANY($1)"

	rows, err := a.db.Query(context.Background(), sql, actorIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var actor Actor
		err := rows.Scan(&actor.ID, &actor.Name, &actor.Gender, &actor.BirthDate)
		if err != nil {
			return nil, err
		}
		actor.BirthDateJson = actor.BirthDate.Format("2006-01-02")
		actorMap[actor.ID] = append(actorMap[actor.ID], actor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return actorMap, nil
}
