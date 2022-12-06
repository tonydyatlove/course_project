package repository

import (
	"log"
	"os"
	"testing"

	"github.com/tonydyatlove/leomessi/internal/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

var (
	PoolM *mongo.Client
)

type ServiceM struct { //service new
	rps Repository
}

func NewServiceM(NewRps Repository) *ServiceM { //create
	return &ServiceM{rps: NewRps}
}
func TestMain(m *testing.M) {
	pool, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatalf("Bad connection: %v", err)
	}
	PoolM = pool
	run := m.Run()
	os.Exit(run)
}
func TestCreateMongo(t *testing.T) {
	testValidData := []model.Game{
		{
			Teams:     "Barcelona - Real Madrid",
			Score:    "5-0",
			xG:      "3.75-0.6",
			MVP: "Lionel Messi",
		},
		{
			Teams:     "Manchester City - Liverpool",
			Score:    "4-1",
			xG:      "2.35-0.93",
			MVP: "Kevin De Bruyne",
		},
	}
	testNoValidData := []model.Person{
		{
			Teams:     "Chelsea - Arsenal",
			Score:    "1-2",
			xG:      "1.01-1.69",
			MVP: "Gabriel Martinelli",
		},
		{
			Teams:     "PSG - Real Madrid",
			Score:    "1-3",
			xG:      "1.48-0.7",
			MVP: "Karim Benzema",
		},
	}
	rps := NewServiceM(&MRepository{Pool: PoolM})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, p := range testValidData {
		_, err := rps.rps.CreateUser(ctx, &p)
		require.NoError(t, err, "create error")
	}
	for _, p := range testNoValidData {
		_, err := rps.rps.CreateUser(ctx, &p)
		require.Error(t, err, "create error")
	}
}
func TestSelectAllMongo(t *testing.T) {
	rps := NewServiceM(&MRepository{Pool: PoolM})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	users, err := rps.rps.GetAllUsers(ctx)
	require.NoError(t, err, "select all: problems with select all games")
	require.Equal(t, 2, len(users), "select all: the values are`t equals")
	collection := PoolM.Database("admin").Collection("admin")
	_, err = collection.InsertOne(ctx, bson.D{
		{Key: "id", Value: "6"},
		{Key: "game", Value: "Zenit - Ural"},
		{Key: "xG", Value: "0.5-0.8"},
		{Key: "score", Value: "0-0"},
		{Key: "mvp", Value: "Malcom"},
	})
	require.NoError(t, err, "select all: insert error")
	users, _ = rps.rps.GetAllUsers(ctx)
	require.Equal(t, 3, len(users), "select all: the values are`t equals")
	require.NotEqual(t, 4, len(users), "select all: the values are equals")
}
func TestSelectByIdMongo(t *testing.T) {
	rps := NewServiceM(&MRepository{Pool: PoolM})
	ctx, cancel := context.WithCancel(context.Background())
	_, err := rps.rps.GetUserByID(ctx, "1223")
	require.NoError(t, err, "select user by id: this id dont exist")
	_, err = rps.rps.GetUserByID(ctx, "20")
	require.Error(t, err, "select user by id: this id already exist")
	cancel()
}
func TestUpdateMongo(t *testing.T) {
	rps := NewServiceM(&MRepository{Pool: PoolM})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	testValidData := []model.Person{
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
	testNoValidData := []model.Person{
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
		err := rps.rps.UpdateUser(ctx, "1223", &p)
		require.NoError(t, err, "update error")
	}
	for _, p := range testNoValidData {
		err := rps.rps.UpdateUser(ctx, "1223", &p)
		require.Error(t, err, "update error")
	}
}