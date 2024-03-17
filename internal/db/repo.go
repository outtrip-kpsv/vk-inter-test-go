package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"vk-inter-test-go/internal/config"
	"vk-inter-test-go/internal/db/repo"
)

type DBRepo struct {
	db *pgxpool.Pool

	User       repo.UserRepository
	Role       repo.RoleRepository
	Actor      repo.ActorRepository
	Movie      repo.MovieRepository
	MovieActor repo.MovieActorRepository
}

func NewDBRepo(conf *config.ConfSrv) *DBRepo {
	db, _ := NewDb(conf.Options.DbString("postgres"))
	return &DBRepo{
		db: db,

		User:       repo.NewUserRepository(db, conf.Logger.Named("RepoUser")),
		Actor:      repo.NewActorRepository(db, conf.Logger.Named("RepoActor")),
		Role:       repo.NewRoleRepository(db, conf.Logger.Named("RepoRole")),
		Movie:      repo.NewMovieRepository(db, conf.Logger.Named("RepoMovie")),
		MovieActor: repo.NewMovieActorRepository(db, conf.Logger.Named("RepoMovieActor")),
	}
}

func (d *DBRepo) Close() {
	d.db.Close()
}
