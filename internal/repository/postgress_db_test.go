package repository

import (
	"context"
	"log"
	"testing"

	"github.com/tonydyatlove/leomessi/internal/model"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

var (
	Pool *pgxpool.Pool
)

type Service struct { // Service new
	rps Repository
}

func NewService(newRps Repository) *Service { // create
	return &Service{rps: newRps}
}

func TestCreate(t *testing.T) {
	testValidData := []model.Person{
		{
			Teams:     "Roma - Lazio",
			xG:    "1.39-1.58",
			Age:      "0-2",
			MVP:      "Milinkovic-Savic",
		},
		{
			Teams:     "Barcelona - Atletico Madrid",
			xG:    "1.87-0.2",
			Age:      "3-0",
			MVP:      "Andreas Iniesta",
		},
	}
	testNoValidData := []model.Person{
		{
			Teams:     "Bayern - Dortmund",
			xG:    "2.77-2.45",
			Age:      "3-3",
			MVP:      "Leroy Sane",
		},
		{
			Teams:     "Liverpool - Chelsea",
			xG:    "1.28-2.35",
			Age:      "2-1",
			MVP:      "Mohammed Salah",
		},
		{
			Teams:     "Zenit-Arsenal",
			xG:    "0.25-4.12",
			Age:      "0-6",
			MVP:      "Mykhailo Mudryk",
		},
	}
	rps := NewService(&PRepository{Pool: Pool})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, p := range testValidData {
		_, err := rps.rps.CreateGame(ctx, &p)
		require.NoError(t, err, "create error")
	}
	for _, p := range testNoValidData {
		_, err := rps.rps.CreateGame(ctx, &p)
		require.Error(t, err, "create error")
	}
}
func TestSelectAll(t *testing.T) {
	rps := NewService(&PRepository{Pool: Pool})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := model.Game{
		ID:       "8",
		Teams:     "Manchester City - Real Madrid",
		xG:    "4.63-0.4",
		Score:     "6-0",
		MVP: "Phil Foden",
	}

	users, err := rps.rps.GetAllGames(ctx)
	require.NoError(t, err, "select all: problems with select all games")
	require.Equal(t, 2, len(games), "select all: the values are`t equals")

	_, err = Pool.Exec(ctx, "insert into persons(id,teams,xG,score,mvp) values($1,$2,$3,$4,$5)", &p.ID, &p.Teams, &p.xG, &p.Score, &p.MVP)
	require.NoError(t, err, "select all: insert error")
	users, err = rps.rps.GetAllGames(ctx)
	if err != nil {
		defer log.Fatalf("error with select all: %v", err)
	}
	require.NotEqual(t, 5, len(games), "select all: the values are equals")
}

func TestSelectById(t *testing.T) {
	rps := NewService(&PRepository{Pool: Pool})
	ctx, cancel := context.WithCancel(context.Background())
	_, err := rps.rps.GetGameByID(ctx, "12")
	require.NoError(t, err, "select user by id: this id not exist")
	_, err = rps.rps.GetGameByID(ctx, "20")
	require.Error(t, err, "select user by id: this id already exist")
	cancel()
}

func TestUpdate(t *testing.T) {
	rps := NewService(&PRepository{Pool: Pool})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testValidData := []*model.Person{
		{
			Teams:  "Manchester City - Real Madrid",
			Score: "4-3",
			xG:   "2.45-1.36",
		},
		{
			Teams:  "Manchester United - Real Madrid",
			Score: "4-3",
			xG:   "3.12-2.66",
		},
	}
	testNoValidData := []*model.Person{
		{
			Teams:  "Bayern - Wolfsburg",
			Score: "4-0",
			xG:   "3.15-0.62",
		},
		{
			Teams:  "Manchester City - Real Madrid",
			Score: "4-3",
			xG:   "2.45-1.36",
		},
		{
			Teams:  "Juventus - Inter",
			Score: "3-0",
			xG:   "1.76-1.08",
		},
	}
	for _, p := range testValidData {
		err := rps.rps.UpdateGame(ctx, "d57d1026-c79a-443d-9d81-714381a37a80", p)
		require.NoError(t, err, "update error")
	}
	for _, p := range testNoValidData {
		err := rps.rps.UpdateGame(ctx, "bb839db7-4be3-41a8-a53b-403ad26593ca", p)
		require.Error(t, err, "update error")
	}
	err := rps.rps.UpdateGame(ctx, "bb839db7-4be3-41a8-a53b-403ad26593ca", testValidData[0])
	require.Error(t, err, "update error")
}
func TestPRepository_UpdateAuth(t *testing.T) {
	rps := NewService(&PRepository{Pool: Pool})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := rps.rps.UpdateAuth(ctx, "d57d1026-c79a-443d-9d81-714381a37a80",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTg0OTk0NzgsImp0aSI6IjNhZjYyMjY5LTAxZmYtNGM2YS04MmUwLTBhNjIwZTVlY2ZmZCIsInVzZXJuYW1lIjoiRWdvclRpaG9ub3YifQ.d4kAjYeGkObPF-kcm7TaFRducO7rsUjabu_8h-Sy8ZE")
	require.NoError(t, err, "thereis an error")
	err = rps.rps.UpdateAuth(ctx, "3",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTg0OTk0NzgsImp0aSI6IjNhZjYyMjY5LTAxZmYtNGM2YS04MmUwLTBhNjIwZTVlY2ZmZCIsInVzZXJuYW1lIjoiRWdvclRpaG9ub3YifQ.d4kAjYeGkObPF-kcm7TaFRducO7rsUjabu_8h-Sy8ZE")
	require.Error(t, err, "there isnt an error")
}
func TestPRepository_SelectByIdAuth(t *testing.T) {
	rps := NewService(&PRepository{Pool: Pool})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := rps.rps.SelectByIDAuth(ctx, "d57d1026-c79a-443d-9d81-714381a37a80")
	require.NoError(t, err, "there is an error")
	_, err = rps.rps.SelectByIDAuth(ctx, "3")
	require.Error(t, err, "there isn`t an error")
}

func TestPRepository_Delete(t *testing.T) {
	rps := NewService(&PRepository{Pool: Pool})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := rps.rps.SelectByIDAuth(ctx, "d57d1026-c79a-443d-9d81-714381a37a80")
	require.NoError(t, err, "there is an error")
	_, err = rps.rps.SelectByIDAuth(ctx, "3")
	require.Error(t, err, "there isn`t an error")
}