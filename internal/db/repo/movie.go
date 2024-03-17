package repo

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

type MovieRepositoryImpl struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewMovieRepository(db *pgxpool.Pool, logger *zap.Logger) *MovieRepositoryImpl {
	logger.Info("create")
	return &MovieRepositoryImpl{db: db, logger: logger}
}

type Movie struct {
	ID              int       `db:"id" json:"ID"`
	Title           string    `db:"title" json:"title,omitempty"`
	Description     string    `db:"description" json:"description,omitempty"`
	ReleaseDateJson string    `db:"-" json:"releaseDate,omitempty"`
	ReleaseDate     time.Time `db:"release_date" json:"-"`
	Rating          int       `db:"rating" json:"rating,omitempty"`
}

type MovieRepository interface {
	CreateMovie(movie *Movie) error
	GetMovieMapByIDs(movieIDs []int, orderBy string) (map[int]Movie, error)
	DeleteMovieById(id int) (int64, error)
	UpdateMovie(movie Movie) (int64, error)
	GetMovieById(id int) (Movie, error)
	GetMoviesLikeTitle(title string, orderBy string) ([]Movie, error)
}

func (m MovieRepositoryImpl) CreateMovie(movie *Movie) error {
	sql := "INSERT INTO movies (title, description, release_date, rating) VALUES ($1, $2, $3, $4) RETURNING id"
	err := m.db.QueryRow(context.Background(), sql, movie.Title, movie.Description, movie.ReleaseDate, movie.Rating).Scan(&movie.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m MovieRepositoryImpl) GetMovieMapByIDs(movieIDs []int, orderBy string) (map[int]Movie, error) {
	movieMap := make(map[int]Movie)

	sql := "SELECT id, title, description, release_date, rating FROM movies WHERE id = ANY($1) ORDER BY "
	switch orderBy {
	case "rating":
		sql += "rating DESC"
	case "date":
		sql += "release_date DESC"
	case "title":
		sql += "title"
	default:
		sql += "rating DESC"
	}

	rows, err := m.db.Query(context.Background(), sql, movieIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var movie Movie
		err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating)
		movie.ReleaseDateJson = movie.ReleaseDate.Format("2006-01-02")
		if err != nil {
			return nil, err
		}
		movieMap[movie.ID] = movie
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movieMap, nil
}

func (m MovieRepositoryImpl) DeleteMovieById(id int) (int64, error) {
	sql := "DELETE FROM movies WHERE id = $1"
	res, err := m.db.Exec(context.Background(), sql, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func (m MovieRepositoryImpl) UpdateMovie(movie Movie) (int64, error) {
	sql := "UPDATE movies SET title = $2, description = $3, release_date = $4, rating = $5 WHERE id = $1"
	res, err := m.db.Exec(context.Background(), sql, movie.ID, movie.Title, movie.Description, movie.ReleaseDate, movie.Rating)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func (m MovieRepositoryImpl) GetMovieById(id int) (Movie, error) {
	sql := "SELECT id, title, description, release_date, rating FROM movies WHERE id = $1 ORDER BY rating DESC "
	row := m.db.QueryRow(context.Background(), sql, id)

	var movie Movie

	err := row.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating)
	if err != nil {
		return Movie{}, err
	}
	movie.ReleaseDateJson = movie.ReleaseDate.Format("2006-01-02")
	return movie, nil
}

func (m MovieRepositoryImpl) GetMoviesLikeTitle(title string, orderBy string) ([]Movie, error) {
	sql := "SELECT * FROM movies WHERE title LIKE '%' || $1 || '%' ORDER BY "
	switch orderBy {
	case "rating":
		sql += "rating DESC"
	case "date":
		sql += "release_date DESC"
	case "title":
		sql += "title"
	default:
		sql += "rating DESC"
	}
	rows, err := m.db.Query(context.Background(), sql, title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie

	for rows.Next() {
		var movie Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating); err != nil {
			return nil, err
		}
		movie.ReleaseDateJson = movie.ReleaseDate.Format("2006-01-02")
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return movies, nil
}
