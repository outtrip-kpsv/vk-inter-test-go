package tests_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"sort"
	"testing"
	"time"
	"vk-inter-test-go/internal/bl"
	"vk-inter-test-go/internal/config"
	"vk-inter-test-go/internal/db"
	"vk-inter-test-go/internal/db/repo"
	"vk-inter-test-go/internal/io/models"
	utilsJwt "vk-inter-test-go/internal/utils"
)

func init() {
	s, err := config.InitConfServ()
	if err != nil {
		return
	}
	s.Logger.Info("init")
}

type mockMovieRepo struct {
	mock.Mock
}

type mockActorRepo struct{}
type mockActorMovieRepo struct{}
type mockUserRepo struct{}
type mockRoleRepo struct{}

func (m mockRoleRepo) GetRoleById(id int) (repo.Role, error) {
	return repo.Role{
		ID:   2,
		Name: "admin",
	}, nil
}

func (m mockUserRepo) GetUserByLogin(login string) (repo.User, error) {
	if login != "testuser" {
		return repo.User{}, errors.New("err")
	}
	password, _ := utilsJwt.HashPassword("password")

	return repo.User{
		Login:  "testuser",
		Pass:   password,
		RoleID: 2,
	}, nil
}

func (m mockUserRepo) CreateUser(user repo.User) error {
	return nil
}

var (
	mok = &db.DBRepo{
		User:       &mockUserRepo{},
		Role:       &mockRoleRepo{},
		Actor:      &mockActorRepo{},
		Movie:      &mockMovieRepo{},
		MovieActor: &mockActorMovieRepo{},
	}

	exempl = bl.NewBL(mok, zap.NewExample())
)

func (m mockMovieRepo) CreateMovie(movie *repo.Movie) error {
	return nil
}

func (m mockMovieRepo) GetMovieMapByIDs(movieIDs []int, orderBy string) (map[int]repo.Movie, error) {
	res := make(map[int]repo.Movie)
	res[1] = repo.Movie{
		ID:              1,
		Title:           "Oppenheimer",
		Description:     "The story of J. Robert Oppenheimer's role in the development of the atomic bomb during World War II.",
		ReleaseDateJson: "2023-07-21",
		Rating:          8,
	}
	res[2] = repo.Movie{
		ID:              2,
		Title:           "Retreat",
		Description:     "Kate and Martin escape from personal tragedy to an Island Retreat. Cut off from the outside world, their attempts to recover are shattered when a man is washed ashore, with news of airborne killer disease that is sweeping through Europe.",
		ReleaseDateJson: "2011-10-14",
		Rating:          5,
	}
	return res, nil
}

func (m mockMovieRepo) DeleteMovieById(id int) (int64, error) {
	return 1, nil
}

func (m mockMovieRepo) UpdateMovie(movie repo.Movie) (int64, error) {
	return 1, nil
}

func (m mockMovieRepo) GetMovieById(id int) (repo.Movie, error) {
	if id > 200 {
		return repo.Movie{}, errors.New("err")
	}
	return repo.Movie{
		ID:          1,
		Title:       "Old Title",
		Description: "Old Description",
		ReleaseDate: time.Now(),
		Rating:      5,
	}, nil
}

func (m mockMovieRepo) GetMoviesLikeTitle(title string, orderBy string) ([]repo.Movie, error) {
	res := []repo.Movie{
		{ID: 1,
			Title:           "Oppenheimer",
			Description:     "The story of J. Robert Oppenheimer's role in the development of the atomic bomb during World War II.",
			ReleaseDateJson: "2023-07-21",
			Rating:          8},
	}
	return res, nil
}

func (m mockActorMovieRepo) CreateMovieActorRelation(movieID int, actorIDs []int) error {
	return nil
}

func (m mockActorMovieRepo) GetRelationByMovieIDs(movieIDs []int) (map[int][]int, error) {
	res := make(map[int][]int)
	res[1] = []int{1}
	res[2] = []int{1}
	return res, nil
}

func (m mockActorMovieRepo) GetRelationByActorIDs(actorIDs []int) (map[int][]int, error) {
	if len(actorIDs) == 0 {
		return nil, errors.New("err")
	}
	res := make(map[int][]int)
	res[1] = []int{1, 2}
	return res, nil
}

func (m *mockActorRepo) DeleteActorByName(name string) (int64, error) {
	if name == "test" {
		return 1, nil

	}
	return 0, errors.New("not actor")
}

func (m *mockActorRepo) UpdateActor(actor repo.Actor) (int64, error) {
	if actor.Name == "err" {
		return 0, errors.New("err")
	}
	return 1, nil
}

func (m *mockActorRepo) GetAllActorsLikeName(name string, orderBy string) ([]repo.Actor, error) {
	if name == "!" {
		return nil, errors.New("err")
	}
	if name == "err" {
		return nil, nil
	}
	res := []repo.Actor{
		{
			ID:            1,
			Name:          "Cillian Murphy",
			Gender:        "male",
			BirthDateJson: "1976-05-25",
		},
	}
	return res, nil
}

func (m *mockActorRepo) GetActorById(id int) (repo.Actor, error) {
	if id == -1 {
		return repo.Actor{}, errors.New("err")
	}
	return repo.Actor{
		ID:            1,
		Name:          "test",
		Gender:        "male",
		BirthDateJson: "1984-02-24",
		BirthDate:     time.Now(),
	}, nil
}

func (m *mockActorRepo) GetActorByName(name string) (repo.Actor, error) {
	if name == "Actor1" {
		return repo.Actor{ID: 1, Name: "Actor1"}, nil
	}
	if name == "Actor2" {
		return repo.Actor{ID: 2, Name: "Actor2"}, nil
	}
	return repo.Actor{}, errors.New("err")
}

func (m *mockActorRepo) GetActorMapByIDs(actorIDs []int) (map[int][]repo.Actor, error) {
	res := make(map[int][]repo.Actor)
	res[1] = []repo.Actor{
		{
			ID:            1,
			Name:          "Cillian Murphy",
			Gender:        "male",
			BirthDateJson: "1976-05-25",
		},
	}
	return res, nil
}

func (m *mockActorRepo) CreateActor(actor *repo.Actor) error {
	if actor.Name == "err" {
		return errors.New("err")
	}
	actor.ID = 1
	return nil
}

//---------------
//TESTS
//---------------

func TestCreateActor1(t *testing.T) {
	actor := repo.Actor{
		ID:            0,
		Name:          "test",
		Gender:        "male",
		BirthDateJson: "1984-02-24",
		BirthDate:     time.Time{},
	}

	actorNew, err := exempl.CreateActor(actor)

	assert.NoError(t, err, "Unexpected error during actor creation")

	assert.Equal(t, 1, actorNew.ID, "Expected actor ID to be 1")
}
func TestCreateActor2(t *testing.T) {
	actor := repo.Actor{
		ID:            0,
		Name:          "err",
		Gender:        "male",
		BirthDateJson: "1984-02-24",
		BirthDate:     time.Time{},
	}

	actorNew, err := exempl.CreateActor(actor)

	assert.Error(t, err, "Unexpected error during actor creation")

	assert.Equal(t, 0, actorNew.ID, "Expected actor ID to be 1")
}
func TestGetAllMoviesByNameActor(t *testing.T) {
	actorName := "Cillian"
	orderBy := "rating"

	expectedMovies := []models.MovieIo{
		{
			Movie: repo.Movie{
				ID:              1,
				Title:           "Oppenheimer",
				Description:     "The story of J. Robert Oppenheimer's role in the development of the atomic bomb during World War II.",
				ReleaseDateJson: "2023-07-21",
				Rating:          8,
			},
			Actors: []repo.Actor{
				{
					ID:            1,
					Name:          "Cillian Murphy",
					Gender:        "male",
					BirthDateJson: "1976-05-25",
				},
			},
		},
		{
			Movie: repo.Movie{
				ID:              2,
				Title:           "Retreat",
				Description:     "Kate and Martin escape from personal tragedy to an Island Retreat. Cut off from the outside world, their attempts to recover are shattered when a man is washed ashore, with news of airborne killer disease that is sweeping through Europe.",
				ReleaseDateJson: "2011-10-14",
				Rating:          5,
			},
			Actors: []repo.Actor{
				{
					ID:            1,
					Name:          "Cillian Murphy",
					Gender:        "male",
					BirthDateJson: "1976-05-25",
				},
			},
		},
	}

	actualMovies, err := exempl.GetAllMoviesByNameActor(actorName, orderBy)

	assert.NoError(t, err, "Unexpected error during GetAllMoviesByNameActor")

	assert.Equal(t, len(expectedMovies), len(actualMovies), "Number of movies is not as expected")

	sort.Slice(expectedMovies, func(i, j int) bool {
		return expectedMovies[i].Movie.ID < expectedMovies[j].Movie.ID
	})
	sort.Slice(actualMovies, func(i, j int) bool {
		return actualMovies[i].Movie.ID < actualMovies[j].Movie.ID
	})
	assert.Equal(t, expectedMovies, actualMovies)
}

func TestDeleteActor1(t *testing.T) {
	deletedCount, err := exempl.DeleteActor("test")
	assert.NoError(t, err, "Unexpected error")
	assert.Equal(t, int64(1), deletedCount, "Expected one actor to be deleted")
}

func TestDeleteActor2(t *testing.T) {
	deletedCount, _ := exempl.DeleteActor("tet")
	assert.Equal(t, int64(0), deletedCount, "Expected one actor to be deleted")
}

func TestUpdateActor1(t *testing.T) {

	actor := repo.Actor{
		ID:            1,
		Name:          "test",
		Gender:        "male",
		BirthDateJson: "1984-02-24",
		BirthDate:     time.Now(),
	}

	updatedActor, err := exempl.UpdateActor(actor)

	assert.NoError(t, err, "Unexpected error")
	assert.Equal(t, actor.ID, updatedActor.ID, "Expected ID to match")
	assert.Equal(t, actor.Name, updatedActor.Name, "Expected Name to match")
	assert.Equal(t, actor.Gender, updatedActor.Gender, "Expected Gender to match")
}
func TestUpdateActor2(t *testing.T) {

	actor := repo.Actor{
		ID: 1,
	}

	updatedActor, err := exempl.UpdateActor(actor)

	assert.NoError(t, err, "Unexpected error")
	assert.Equal(t, actor.ID, updatedActor.ID, "Expected ID to match")
	assert.Equal(t, "test", updatedActor.Name, "Expected Name to match")
	assert.Equal(t, "male", updatedActor.Gender, "Expected Gender to match")
}
func TestUpdateActor3(t *testing.T) {

	actor := repo.Actor{
		ID: -1,
	}

	updatedActor, err := exempl.UpdateActor(actor)

	assert.Error(t, err, "Unexpected error")
	assert.Equal(t, 0, updatedActor.ID, "Expected ID to match")
}
func TestUpdateActor4(t *testing.T) {

	actor := repo.Actor{
		ID:   1,
		Name: "err",
	}

	updatedActor, err := exempl.UpdateActor(actor)

	assert.Error(t, err, "Unexpected error")
	assert.Equal(t, "", updatedActor.Name, "Expected Name to match")
}

func TestGetAllActorsLikeName1(t *testing.T) {
	testName := "Cillian"
	testOrderBy := "rating"

	expectedActors := []models.ActorIo{
		{Actor: repo.Actor{
			ID:            1,
			Name:          "Cillian Murphy",
			Gender:        "male",
			BirthDateJson: "1976-05-25",
		}, Movies: []repo.Movie{
			{
				ID:              1,
				Title:           "Oppenheimer",
				Description:     "The story of J. Robert Oppenheimer's role in the development of the atomic bomb during World War II.",
				ReleaseDateJson: "2023-07-21",
				Rating:          8,
			},
			{
				ID:              2,
				Title:           "Retreat",
				Description:     "Kate and Martin escape from personal tragedy to an Island Retreat. Cut off from the outside world, their attempts to recover are shattered when a man is washed ashore, with news of airborne killer disease that is sweeping through Europe.",
				ReleaseDateJson: "2011-10-14",
				Rating:          5,
			},
		}},
	}

	actualActors, err := exempl.GetAllActorsLikeName(testName, testOrderBy)

	assert.NoError(t, err, "Unexpected error")
	assert.ElementsMatch(t, expectedActors, actualActors, "Actors do not match")
}

func TestGetAllActorsLikeName2(t *testing.T) {
	testName := "!"
	testOrderBy := "rating"
	_, err := exempl.GetAllActorsLikeName(testName, testOrderBy)

	assert.Error(t, err, "Unexpected error")
}

func TestGetAllActorsLikeName3(t *testing.T) {
	testName := "err"
	testOrderBy := "rating"
	_, err := exempl.GetAllActorsLikeName(testName, testOrderBy)

	assert.Error(t, err, "Unexpected error")
}

func TestGetAllMoviesByTitle(t *testing.T) {

	testTitle := "eimer"
	testOrderBy := "rating"

	expectedMovies := []models.MovieIo{
		{Movie: repo.Movie{
			ID:              1,
			Title:           "Oppenheimer",
			Description:     "The story of J. Robert Oppenheimer's role in the development of the atomic bomb during World War II.",
			ReleaseDateJson: "2023-07-21",
			Rating:          8,
		}, Actors: []repo.Actor{{
			ID:            1,
			Name:          "Cillian Murphy",
			Gender:        "male",
			BirthDateJson: "1976-05-25",
		},
		}},
	}

	actualMovies, err := exempl.GetAllMoviesByTitle(testTitle, testOrderBy)

	assert.NoError(t, err, "Unexpected error")
	assert.ElementsMatch(t, expectedMovies, actualMovies, "Movies do not match")
}

func TestCreateUser(t *testing.T) {
	testUser := repo.User{
		Login: "testuser",
		Pass:  "password",
	}

	actualToken, err := exempl.CreateUser(testUser)

	assert.NoError(t, err, "Unexpected error")
	assert.Equal(t, true, exempl.CheckJwt(actualToken), "Generated token does not match")
	name, _ := utilsJwt.ExtractUsernameFromToken(actualToken)

	assert.Equal(t, testUser.Login, name, "names does not match")

}

func TestAuthUser(t *testing.T) {
	testUser := repo.User{
		Login: "testuser",
		Pass:  "password",
	}

	actualToken, err := exempl.AuthUser(testUser)

	assert.NoError(t, err, "Unexpected error")
	assert.Equal(t, true, exempl.CheckJwt(actualToken), "Generated token does not match")
}

func TestSanitize(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
		changed  bool
	}{
		{"Hello, World!", "Hello World", false},      // Заменяется запятая
		{"123-456-7@89", "123-456-789", false},       // Не изменяется
		{"!@#$%^&*()_+", "", false},                  // Удалены все спецсимволы
		{"This is a test.", "This is a test", false}, // Удалена точка
		{"", "", true},                       // Пустая строка
		{"   ", "   ", true},                 // Только пробелы
		{"hello_world", "helloworld", false}, // Удален символ подчеркивания
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result, changed := utilsJwt.Sanitize(tc.input)

			if result != tc.expected {
				t.Errorf("Expected sanitized string: %s, got: %s", tc.expected, result)
			}

			if changed != tc.changed {
				t.Errorf("Expected changed value: %v, got: %v", tc.changed, changed)
			}
		})
	}
}

func TestCheckRole(t *testing.T) {
	assert.True(t, exempl.CheckRole("testuser", "admin"))
	assert.False(t, exempl.CheckRole("test_error", "admin"))
	assert.False(t, exempl.CheckRole("test", "admin"))
	assert.False(t, exempl.CheckRole("", "admin"))
}

func TestUpdateMovie(t *testing.T) {

	mockMovie := repo.Movie{
		ID:          1,
		Title:       "Old Title",
		Description: "Old Description",
		ReleaseDate: time.Now(),
		Rating:      5,
	}

	newMovie, err := exempl.UpdateMovie(repo.Movie{})
	assert.NoError(t, err)
	assert.Equal(t, mockMovie.Title, newMovie.Title)
	assert.Equal(t, mockMovie.Description, newMovie.Description)
	assert.Equal(t, mockMovie.Rating, newMovie.Rating)

	invalidReleaseDate := "invalid date"
	_, err = exempl.UpdateMovie(repo.Movie{ID: mockMovie.ID, ReleaseDateJson: invalidReleaseDate})
	assert.Error(t, err)

	_, err = exempl.UpdateMovie(repo.Movie{ID: 999})
	assert.Error(t, err)
}

func TestCreateMovie(t *testing.T) {
	mockMovie := models.MovieIo{
		Movie: repo.Movie{
			ID:          1,
			Title:       "Test Movie",
			Description: "Test Description",
			ReleaseDate: time.Now(),
			Rating:      8,
		},
		Actors: []repo.Actor{
			{ID: 5, Name: "Actor1"},
			{ID: 2, Name: "Actor2"},
			{ID: 1, Name: "Actor3"},
		},
	}

	createdMovie, err := exempl.CreateMovie(mockMovie)

	assert.NoError(t, err)
	assert.Equal(t, mockMovie, createdMovie)

}

func TestDeleteMovie(t *testing.T) {
	rowsAffected, err := exempl.DeleteMovie(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

}
