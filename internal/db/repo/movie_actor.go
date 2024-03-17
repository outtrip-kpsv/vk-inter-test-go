package repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"strings"
)

type MovieActorRepositoryImpl struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewMovieActorRepository(db *pgxpool.Pool, logger *zap.Logger) *MovieActorRepositoryImpl {
	logger.Info("create")
	return &MovieActorRepositoryImpl{db: db, logger: logger}
}

type MovieActor struct {
	ID      int `db:"id" json:"-"`
	MovieID int `db:"movie_id" json:"movieID"`
	ActorID int `db:"actor_id" json:"actorID"`
}

type MovieActorRepository interface {
	CreateMovieActorRelation(movieID int, actorIDs []int) error
	GetRelationByMovieIDs(movieIDs []int) (map[int][]int, error)
	GetRelationByActorIDs(actorIDs []int) (map[int][]int, error)
}

func (m MovieActorRepositoryImpl) CreateMovieActorRelation(movieID int, actorIDs []int) error {
	var values []interface{}
	values = append(values, movieID)

	for _, actorID := range actorIDs {
		values = append(values, actorID)
	}

	sql := "INSERT INTO movies_actors (movie_id, actor_id) VALUES "
	placeholders := make([]string, len(actorIDs))
	for i := range actorIDs {
		placeholders[i] = fmt.Sprintf("($1, $%d)", i+2)
	}
	sql += strings.Join(placeholders, ", ")

	_, err := m.db.Exec(context.Background(), sql, values...)
	return err
}

func (m MovieActorRepositoryImpl) GetRelationByActorIDs(actorIDs []int) (map[int][]int, error) {
	relations := make(map[int][]int)

	query := "SELECT actor_id, movie_id FROM movies_actors WHERE actor_id = ANY($1)"
	rows, err := m.db.Query(context.Background(), query, actorIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var actorID, movieID int
		if err := rows.Scan(&actorID, &movieID); err != nil {
			return nil, err
		}
		relations[actorID] = append(relations[actorID], movieID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return relations, nil
}

func (m MovieActorRepositoryImpl) GetRelationByMovieIDs(movieIDs []int) (map[int][]int, error) {
	relations := make(map[int][]int)

	query := "SELECT movie_id, actor_id FROM movies_actors WHERE movie_id = ANY($1)"
	rows, err := m.db.Query(context.Background(), query, movieIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var movieID, actorID int
		if err := rows.Scan(&movieID, &actorID); err != nil {
			return nil, err
		}
		relations[movieID] = append(relations[movieID], actorID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return relations, nil
}
